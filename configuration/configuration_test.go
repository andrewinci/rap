package configuration

import (
	"os"
	"testing"
)

func TestLoadConfigHappyPath(t *testing.T) {
	testEndpoint, testUsername := "exampleEndpoint", "testUsername"
	os.Setenv("KAFKA_ENDPOINT", testEndpoint)
	os.Setenv("KAFKA_USERNAME", testUsername)
	os.Setenv("SR_ENDPOINT", testEndpoint)
	res, err := LoadConfiguration("../example/sasl_configuration.yaml")
	if err != nil {
		t.Fail()
	}
	if res.Kafka.ClusterEndpoint != testEndpoint {
		t.Fail()
	}
	if res.Kafka.Sasl.Username != testUsername {
		t.Fail()
	}
	if res.Producers[0].NumberOfMessages != 2 {
		t.Fail()
	}
}

func TestLoadConfigHappyPath2(t *testing.T) {
	res, err := LoadConfiguration("../example/local_cluster.yaml")
	if err != nil {
		t.Fail()
	}
	if res.Kafka.ClusterEndpoint != "localhost:9092" {
		t.Fail()
	}
	if res.Producers[0].NumberOfMessages != 2 {
		t.Fail()
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	_, err := LoadConfiguration("/asdf/asdf/asdf")
	if err == nil {
		t.Fail()
	}
}

func TestLoadAnInvalidConfig(t *testing.T) {
	_, err := LoadConfiguration("../readme.md")
	if err == nil {
		t.Fail()
	}
}

func TestValidateConfiguration_Schema(t *testing.T) {
	c := Configuration{
		Kafka: KafkaConfiguration{
			ClusterEndpoint: "endpoint",
		},
		Producers: []ProducerConfiguration{
			{Avro: AvroGenConfiguration{SchemaName: "test-schema"}},
		}}
	res := validateConfiguration(&c)
	if res == nil {
		// validation should fail because the schema registry config is missing
		t.Fail()
	}
}

func TestValidateConfiguration_KafkaEndpoint(t *testing.T) {
	c := Configuration{
		Producers: []ProducerConfiguration{
			{Avro: AvroGenConfiguration{Schema: SchemaConfiguration{Id: 1, Raw: ""}}}}}
	res := validateConfiguration(&c)
	if res == nil {
		// validation should fail because the kafka endpoint is not defined
		t.Fail()
	}
}

func TestValidateConfiguration_EmptyProducers(t *testing.T) {
	c := Configuration{
		Kafka: KafkaConfiguration{
			ClusterEndpoint: "endpoint",
			SchemaRegistry:  SchemaRegistryConfiguration{Endpoint: "endpoint"},
		},
	}
	res := validateConfiguration(&c)
	if res == nil {
		// validation should fail because there are no producers configured
		t.Fail()
	}
}
