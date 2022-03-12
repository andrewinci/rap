package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type KafkaConfig struct {
	ClusterEndpoint string `yaml:"clusterEndpoint,omitempty"`
}

type AvroConfig struct {
	// avro schema
	Schema string
	// list of generators available
	// in the rules
	Generators map[string]string
	// set of rules to customize the
	// avro generation
	GenerationRules map[string]string `yaml:"generationRules,omitempty"`
}

type ProducerConfig struct {
	Name             string
	NumberOfMessages int `yaml:"numberOfMessages,omitempty"`
	Avro             AvroConfig
}

type Configuration struct {
	Kafka     KafkaConfig
	Producers []ProducerConfig
}

func LoadConfiguration(fileName string) (*Configuration, error) {
	var configuration Configuration
	rawConfig, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(rawConfig, &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}
