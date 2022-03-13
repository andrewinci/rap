package kafka

import (
	"log"
	"sync"

	"github.com/Shopify/sarama"
	c "github.com/andrewinci/rap/configuration"
)

type asyncProducer struct {
	producer sarama.AsyncProducer
	success  int
	errors   int
	wg       *sync.WaitGroup
}

type Producer interface {
	ProduceAsync(key string, value []byte, topicName string)
	Close()
	GetErrorsCount() int
	GetSuccessCount() int
}

func NewProducer(config c.KafkaConfiguration) (Producer, error) {
	var wg sync.WaitGroup
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	switch config.Security {
	case c.None:
	case c.Sasl:
		configureSasl(config.Sasl, conf)
	case c.MTLS:
	}
	producer, err := sarama.NewAsyncProducer([]string{config.ClusterEndpoint}, conf)
	if err != nil {
		return nil, err
	}
	res := asyncProducer{
		producer: producer,
		wg:       &wg,
		success:  0,
		errors:   0,
	}
	// drain success
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			res.success += 1
		}
	}()

	// drain errors
	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range producer.Errors() {
			log.Println("Error", m.Msg.Topic, m.Msg)
			res.errors += 1
		}
	}()

	return &res, nil
}

func configureSasl(config c.SaslConfiguration, saramaConfig *sarama.Config) {
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	saramaConfig.Net.SASL.User = config.Username
	saramaConfig.Net.SASL.Password = config.Password
	saramaConfig.Net.SASL.Handshake = true
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.TLS.Enable = true
}

func (p asyncProducer) ProduceAsync(key string, value []byte, topicName string) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
}

func (p asyncProducer) Close() {
	p.producer.AsyncClose()
	p.wg.Wait()
}

func (p asyncProducer) GetErrorsCount() int {
	return p.errors
}
func (p asyncProducer) GetSuccessCount() int {
	return p.success
}
