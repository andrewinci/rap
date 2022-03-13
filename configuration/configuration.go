package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Kafka     KafkaConfiguration
	Producers []ProducerConfiguration
}

type KafkaConfiguration struct {
	ClusterEndpoint string                      `yaml:"clusterEndpoint"`
	SchemaRegistry  SchemaRegistryConfiguration `yaml:"schemaRegistry"`
	Security        Security
	Sasl            SaslConfiguration `yaml:"sasl"`
}

type Security string

const (
	None Security = "none"
	Sasl Security = "sasl"
	MTLS Security = "mtls"
)

type ProducerConfiguration struct {
	Name             string
	NumberOfMessages int `yaml:"numberOfMessages"`
	Avro             AvroGenConfiguration
	Topic            string `yaml:"topic"`
}

type SchemaRegistryConfiguration struct {
	Endpoint string
	Username string
	Password string
}

type SaslConfiguration struct {
	Username string
	Password string
}

type SchemaConfiguration struct {
	Id  int
	Raw string
}

type AvroGenConfiguration struct {
	// raw avro schema
	Schema SchemaConfiguration
	// schema name only works if the schema registry is configured
	SchemaName string `yaml:"schemaName"`
	// list of generators available
	// in the rules
	Generators map[string]string
	// set of rules to customize the
	// avro generation
	GenerationRules map[string]string `yaml:"generationRules"`
}

// Load the configuration from the provided yaml file path
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
	replaceEnvVariables(&configuration)
	err = validateConfiguration(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

// Replace the fields that references a env variable
// with the actual env variable value
func replaceEnvVariables(config *Configuration) {
	replaceVar := func(v *string) {
		if v == nil || *v == "" {
			return
		}
		if (*v)[0] == '$' {
			*v = os.Getenv((*v)[1:])
		}
	}
	replaceVar(&config.Kafka.ClusterEndpoint)
	replaceVar(&config.Kafka.SchemaRegistry.Endpoint)
	replaceVar(&config.Kafka.SchemaRegistry.Username)
	replaceVar(&config.Kafka.SchemaRegistry.Password)
	replaceVar(&config.Kafka.Sasl.Username)
	replaceVar(&config.Kafka.Sasl.Password)
}

// Validate if the config file is correct
func validateConfiguration(config *Configuration) error {
	if config.Kafka.ClusterEndpoint == "" {
		return fmt.Errorf("validation error: an endpoint for kafka need to be configured in order to produce records")
	}
	// validate non empty producers
	if len(config.Producers) == 0 {
		return fmt.Errorf("validation error: at least one producer must be specified")
	}
	// validate security
	if config.Kafka.Security != Sasl && config.Kafka.Security != None {
		return fmt.Errorf("validation error: security setting `%s` not supported.", config.Kafka.Security)
	}

	schemaRegistryConfigured := config.Kafka.SchemaRegistry.Endpoint != ""
	if schemaRegistryConfigured {
		return nil
	}
	for _, p := range config.Producers {
		// we cannot retrieve the schema from
		// the registry cause if the latter is not configured
		if p.Avro.SchemaName != "" && !schemaRegistryConfigured {
			return fmt.Errorf("validation error: cannot use `SchemaName` when the schema registry is not configured")
		}
	}
	return nil
}
