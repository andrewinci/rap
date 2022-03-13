package avrogen

import (
	"fmt"

	"github.com/hamba/avro"
)

type AvroGenConfiguration struct {
	Schema          string
	Generators      map[string]string
	GenerationRules map[string]string
}

type avroGen struct {
	schema          avro.Schema
	generatorsCache map[string]func() (interface{}, error)
}

type AvroGen interface {
	Generate() ([]byte, error)
	generate(schema avro.Schema, fieldPath string) (interface{}, error)
	generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error)
}

func NewAvroGen(config AvroGenConfiguration, seed int64) (AvroGen, error) {
	// parse avro schema
	schema, err := avro.Parse(config.Schema)
	if err != nil {
		return nil, err
	}

	generatorsCache := map[string]func() (interface{}, error){
		string(avro.Boolean): defaultBooleanFieldGen(seed),
		string(avro.Int):     defaultIntFieldGen(seed),
		string(avro.Long):    defaultLongFieldGen(seed),
		string(avro.Float):   defaultFloatFieldGen(seed),
		string(avro.Double):  defaultDoubleFieldGen(seed),
		string(avro.String):  defaultStringFieldGen(seed),
	}

	fieldGenerators := map[string]func() (interface{}, error){}

	for k, v := range config.Generators {
		fieldGen := newFieldGen(v, seed)
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

	return avroGen{
		schema:          schema,
		generatorsCache: generatorsCache,
	}, nil
}

func (g avroGen) Generate() ([]byte, error) {
	generated, err := g.generate(g.schema, "")
	if err != nil {
		return nil, err
	}
	return avro.Marshal(g.schema, generated)
}

func (g avroGen) generate(schema avro.Schema, fieldPath string) (interface{}, error) {
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

func (g avroGen) generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error) {
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
