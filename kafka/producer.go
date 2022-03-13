package kafka

import (
	"github.com/Shopify/sarama"
	c "github.com/andrewinci/rap/configuration"
)

type syncProducer struct {
	producer sarama.SyncProducer
}

type Producer interface {
	Produce(key string, value []byte, topicName string) error
}

func NewProducer(config c.KafkaConfiguration) (Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	switch config.Security {
	case c.None:
	case c.Sasl:
		configureSasl(config.Sasl, conf)
	case c.MTLS:
	}
	producer, err := sarama.NewSyncProducer([]string{config.ClusterEndpoint}, conf)
	if err != nil {
		return nil, err
	}
	return syncProducer{producer: producer}, nil
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

func (p syncProducer) Produce(key string, value []byte, topicName string) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	return err
}
