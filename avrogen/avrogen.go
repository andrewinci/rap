package avrogen

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	c "github.com/andrewinci/rap/configuration"
	"github.com/hamba/avro"
)

type avroGen struct {
	schema         avro.Schema
	schemaId       int
	generatorsRepo map[string]fieldGen
	randomSource   *rand.Rand
}

type AvroGen interface {
	// return the avro record value and the key
	Generate() ([]byte, string, error)
	generate(schema avro.Schema, fieldPath string) (interface{}, error)
	generateRecord(schema *avro.RecordSchema, parentSchema string) (map[string]interface{}, error)
	getSchema() avro.Schema
}

func NewAvroGen(config c.AvroGenConfiguration, seed int64) (AvroGen, error) {
	// parse avro schema
	schema, err := avro.Parse(config.Schema.Raw)
	if err != nil {
		return nil, err
	}
	randomSource := rand.New(rand.NewSource(seed))
	generatorsRepo := map[string]fieldGen{
		"key":                defaultKeyGen(randomSource),
		string(avro.Boolean): defaultBooleanFieldGen(randomSource),
		string(avro.Int):     defaultIntFieldGen(randomSource),
		string(avro.Long):    defaultLongFieldGen(randomSource),
		string(avro.Float):   defaultFloatFieldGen(randomSource),
		string(avro.Double):  defaultDoubleFieldGen(randomSource),
		string(avro.String):  defaultStringFieldGen(randomSource),
		string(avro.Null):    defaultNullFieldGen(),
	}

	fieldGenerators := map[string]fieldGen{}

	for k, v := range config.Generators {
		fieldGen := newFieldGen(v, randomSource)
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
		generatorsRepo[k] = g
	}

	return avroGen{
		schema:         schema,
		schemaId:       config.Schema.Id,
		generatorsRepo: generatorsRepo,
		randomSource:   randomSource,
	}, nil
}

func (g avroGen) getSchema() avro.Schema {
	return g.schema
}

func (g avroGen) Generate() ([]byte, string, error) {
	generated, err := g.generate(g.schema, "")
	if err != nil {
		return nil, "", err
	}
	key, err := g.generatorsRepo["key"]()
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
	fieldGen, ok := g.generatorsRepo[fieldPath]
	if ok {
		return fieldGen()
	}
	typeGen, ok := g.generatorsRepo[string(schema.Type())]
	if ok {
		return typeGen()
	}
	// no customization found for the field
	if schema.Type() == avro.Union {
		return g.generateUnionField(schema.(*avro.UnionSchema), fieldPath)
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

func (g avroGen) generateUnionField(schema *avro.UnionSchema, fieldPath string) (interface{}, error) {
	// pick a random type among the union options
	tIndex := g.randomSource.Intn(len(schema.Types()))
	return g.generate(schema.Types()[tIndex], fieldPath)
}
