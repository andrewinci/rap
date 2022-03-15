package avrogen

import (
	"testing"

	"github.com/andrewinci/rap/configuration"
)

func TestHappyAvroGenUnion(t *testing.T) {
	testSchema := `
	{
		"type" : "record",
		"name" : "Example",
		"fields" : [
			{ "name": "booleanField", "type": ["null", "boolean"] },
			{ "name": "intField", "type": ["null", "int"] },
			{ "name": "longField", "type": ["null", "long"] },
			{ "name": "floatField", "type": ["null", "float"] },
			{ "name": "doubleField", "type": ["null", "double"] },
			{ "name": "stringField", "type": ["null", "string"] },
			{ "name": "enumField", "type": ["null", { "type": "enum", "name": "Suit", "symbols" : ["SPADES", "HEARTS", "DIAMONDS", "CLUBS"]}]}
		]
	 }
	`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		}}, 321)
	if err != nil {
		t.FailNow()
	}
	rawRes, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	res := rawRes.(map[string]interface{})["booleanField"]
	if res != false {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}

func TestHappyAvroGenUnion2(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": [
					"string",
					{
						"type": "record",
						"name": "Nested",
						"fields": [
							{
								"name": "testNestedField",
								"type": ["int", "float"]
							}
						]
					}
				]
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		},
		GenerationRules: map[string]string{
			".testField.Nested.testNestedField": "intGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"intGen": "{int}[45678923]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}

	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	if res.(map[string]interface{})["testField"].(map[string]interface{})["Nested"].(map[string]interface{})["testNestedField"] != 45678923 {
		t.Fail()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}

func TestHappyAvroGenUnion3(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": [
					"boolean",
					{
						"type": "record",
						"name": "Nested",
						"fields": [
							{ "name": "testNestedField", "type": ["string", "null"] }
						]
					}
				]
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		}}, 1)
	if err != nil {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}

func TestHappyAvroGenUnion4(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": ["boolean", "null"]
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		},
		GenerationRules: map[string]string{
			".testField": "nullGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"nullGen": "{null}[-]{1}",
		}}, 123)
	if err != nil {
		t.FailNow()
	}
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	if res.(map[string]interface{})["testField"] != nil {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}
