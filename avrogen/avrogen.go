package avrogen

import (
	"encoding/binary"
	"fmt"

	c "github.com/andrewinci/rap/configuration"
	"github.com/hamba/avro"
)

type avroGen struct {
	schema          avro.Schema
	schemaId        int
	generatorsCache map[string]func() (interface{}, error)
}

type AvroGen interface {
	// return the avro record value and the key
	Generate() ([]byte, string, error)
	generate(schema avro.Schema, fieldPath string) (interface{}, error)
	generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error)
}

func NewAvroGen(config c.AvroGenConfiguration, seed int64) (AvroGen, error) {
	// parse avro schema
	schema, err := avro.Parse(config.Schema.Raw)
	if err != nil {
		return nil, err
	}

	generatorsCache := map[string]func() (interface{}, error){
		"key":                defaultKeyGen(seed),
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
		schemaId:        config.Schema.Id,
		generatorsCache: generatorsCache,
	}, nil
}

func (g avroGen) Generate() ([]byte, string, error) {
	generated, err := g.generate(g.schema, "")
	if err != nil {
		return nil, "", err
	}
	key, err := g.generatorsCache["key"]()
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate the key, %s", err.Error())
	}
	raw, err := avro.Marshal(g.schema, generated)
	if err != nil {
		return nil, "", fmt.Errorf("unable to marshal the record, %s", err.Error())
	}
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(g.schemaId))
	msg := append([]byte{0x00}, bs...)
	msg = append(msg, raw...)
	return msg, key.(string), err
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
