package main

import (
	"log"
	"time"

	generator "github.com/andrewinci/rap/avrogen"
	config "github.com/andrewinci/rap/config"
)

func main() {
	fileName := "./examples/example1.yaml"
	run(fileName, time.Now().UTC().UnixNano())
}

type producer struct {
	config        config.ProducerConfig
	avroGenerator generator.Generator
	//todo: kafka producer
}

func run(fileName string, seed int64) {
	c, err := config.LoadConfiguration(fileName)
	if err != nil {
		log.Fatalf("Unable to load the config file %s: %s", fileName, err.Error())
	}
	producers := []producer{}
	for _, p := range c.Producers {
		g, err := generator.NewGenerator(p.Avro, seed)
		if err != nil {
			log.Fatalf("Unable to build the producer %s: %s", p.Name, err.Error())
		}
		producers = append(producers, producer{config: p, avroGenerator: g})
	}
	for _, p := range producers {
		for i := 0; i < p.config.NumberOfMessages; i++ {
			p.avroGenerator.Generate()
			//todo: produce
		}
	}
}
