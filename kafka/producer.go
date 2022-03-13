package kafka

import "github.com/Shopify/sarama"

type syncProducer struct {
	producer sarama.SyncProducer
}

type Producer interface {
	Produce(key string, value []byte, topicName string) error
}

type Security int

const (
	None Security = iota
	Sasl
	mTLS //todo
)

type SaslConfiguration struct {
	User     string
	Password string
}

type KafkaConfiguration struct {
	Endpoint string
	Security Security
	Sasl     SaslConfiguration
}

func NewProducer(config KafkaConfiguration) (Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	switch config.Security {
	case None:
	case Sasl:
		configureSasl(config.Sasl, conf)
	case mTLS:
	}
	producer, err := sarama.NewSyncProducer([]string{config.Endpoint}, conf)
	if err != nil {
		return nil, err
	}
	return syncProducer{producer: producer}, nil
}

func configureSasl(config SaslConfiguration, saramaConfig *sarama.Config) {
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	saramaConfig.Net.SASL.User = config.User
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
