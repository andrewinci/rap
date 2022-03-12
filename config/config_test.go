package config

import (
	"testing"
)

func TestLoadConfigHappyPath(t *testing.T) {
	res, err := LoadConfiguration("../examples/example1.yaml")
	if err != nil {
		t.Fail()
	}
	if res.Kafka.ClusterEndpoint != "exampleEndpoint" {
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
