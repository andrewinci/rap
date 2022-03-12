package avrogen

import (
	"fmt"

	config "github.com/andrewinci/rap/config"
	fieldGen "github.com/andrewinci/rap/fieldgen"
	"github.com/hamba/avro"
)

type generator struct {
	schema          avro.Schema
	generatorsCache map[string]func() (interface{}, error)
}

type Generator interface {
	Generate() ([]byte, error)
	generate(schema avro.Schema, fieldPath string) (interface{}, error)
	generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error)
}

func NewGenerator(config config.AvroConfig, seed int64) (Generator, error) {
	// parse avro schema
	schema, err := avro.Parse(config.Schema)
	if err != nil {
		return nil, err
	}

	generatorsCache := map[string]func() (interface{}, error){
		string(avro.String): func() (interface{}, error) { return "hello", nil },
		string(avro.Int):    func() (interface{}, error) { return 123, nil },
	}

	fieldGenerators := map[string]func() (interface{}, error){}

	for k, v := range config.Generators {
		fieldGen := fieldGen.NewFieldGenerator(v, seed)
		fieldGenerators[k] = func() (interface{}, error) {
			res, err := fieldGen()
			return res, err
		}
	}

	for k, v := range config.GenerationRules {
		g, ok := fieldGenerators[v]
		if !ok {
			return nil, fmt.Errorf("missing generator %s for the rule %s", v, k)
		}
		generatorsCache[k] = g
	}

	return generator{
		schema:          schema,
		generatorsCache: generatorsCache,
	}, nil
}

func (g generator) Generate() ([]byte, error) {
	generated, err := g.generate(g.schema, "")
	if err != nil {
		return nil, err
	}
	return avro.Marshal(g.schema, generated)
}

func (g generator) generate(schema avro.Schema, fieldPath string) (interface{}, error) {
	if schema.Type() == avro.Record {
		recordSchema := schema.(*avro.RecordSchema)
		return g.generateRecord(recordSchema, fieldPath)
	}
	fieldGen, ok := g.generatorsCache[fieldPath]
	if ok {
		return fieldGen()
	}
	typeGen, ok := g.generatorsCache[string(schema.Type())]
	if ok {
		return typeGen()
	}
	return nil, fmt.Errorf("no generator found for type %s, path %s", string(schema.Type()), fieldPath)
}

func (g generator) generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, f := range schema.Fields() {
		val, err := g.generate(f.Type(), parentSchema+"."+f.Name())
		if err != nil {
			return nil, err
		}
		res[f.Name()] = val
	}
	return res, nil
}
