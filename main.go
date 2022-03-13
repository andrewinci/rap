package main

import (
	"log"
	"os"
	"sync"
	"time"

	ag "github.com/andrewinci/rap/avrogen"
	c "github.com/andrewinci/rap/configuration"
	k "github.com/andrewinci/rap/kafka"
	"github.com/hamba/avro/registry"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expected 1 argument with the configuration file path")
	}
	// load the configurations
	configFilePath := os.Args[1]
	config, err := c.LoadConfiguration(configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize the schema registry client if required
	var schemaRegistry = buildSchemaRegistry(config.Kafka)

	// replace the schema name with the actual schema in the producers
	for i := range config.Producers {
		avroConfig := &config.Producers[i].Avro
		if avroConfig.SchemaName != "" {
			schemaInfo, err := schemaRegistry.GetLatestSchemaInfo(avroConfig.SchemaName)
			if err != nil {
				log.Fatalf("unable to retrieve the schema %s: %s", avroConfig.SchemaName, err.Error())
			}
			avroConfig.Schema.Raw = schemaInfo.Schema.String()
			avroConfig.Schema.Id = schemaInfo.ID
		}
	}

	// initialize the kafka producer

	kafkaProducer, err := k.NewProducer(config.Kafka)
	if err != nil {
		log.Fatalf("unable to initialize the kafka producer: %s", err.Error())
	}

	// setup all the producers and then execute
	producers := setupProducers(*config, *schemaRegistry, kafkaProducer)

	start := time.Now()
	var wg sync.WaitGroup
	for _, p := range producers {
		wg.Add(1)
		go p(&wg)
	}
	wg.Wait()
	log.Println("All records have been generated successfully")
	kafkaProducer.Close()
	log.Println("All records have been produced successfully")
	elapsed := time.Since(start)

	log.Printf("Produced %d records, %d errors, in %s\n",
		kafkaProducer.GetSuccessCount(),
		kafkaProducer.GetErrorsCount(),
		elapsed)

}

func setupProducers(config c.Configuration, schemaRegistry registry.Client, producer k.Producer) []func(wg *sync.WaitGroup) {
	var producers []func(wg *sync.WaitGroup)
	// setup random avro generators
	seed := time.Now().UnixMilli() //todo: add to the arguments
	log.Printf("Initializing the avro-generators with seed: %d", seed)

	for _, p := range config.Producers {
		gen, err := ag.NewAvroGen(p.Avro, seed)
		if err != nil {
			log.Fatalf("unable to initialize the generator for the producer %s: %s", p.Name, err.Error())
		}
		count := p.NumberOfMessages
		topicName := p.Topic
		producerName := p.Name
		producers = append(producers, func(wg *sync.WaitGroup) {
			defer wg.Done()
			log.Printf("Producer %s started", producerName)
			for i := 0; i < count; i++ {
				msg, key, err := gen.Generate()
				if err != nil {
					log.Fatalf("unable to generate record")
				}
				producer.ProduceAsync(key, msg, topicName)
			}
			log.Printf("Producer %s completed. %d records has been produced", producerName, count)
		})
	}
	return producers
}

func buildSchemaRegistry(config c.KafkaConfiguration) *registry.Client {
	if config.SchemaRegistry.Endpoint != "" {
		schemaRegistry, err := registry.NewClient(config.SchemaRegistry.Endpoint)
		if err != nil {
			log.Fatalf("unable to initialize the schema registry")
		}
		registry.WithBasicAuth(
			config.SchemaRegistry.Username,
			config.SchemaRegistry.Password)(schemaRegistry)
		return schemaRegistry
	}
	return nil
}
