package avrogen

import (
	"testing"

	"github.com/andrewinci/rap/configuration"
)

func TestHappyPathAvroGenArray3(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": {
					"type": "array",
					"items" : "string",
					"default": []
				}
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  3373,
		},
		GenerationRules: map[string]string{
			".testField":       "strGen",
			".testField.len()": "lenGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"strGen": "{string}[test1]{1}",
			"lenGen": "{int}[10]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	generatedArray := res.(map[string]interface{})["testField"].([]interface{})

	if generatedArray[0].(string) != "test1" {
		t.FailNow()
	}
	if len(generatedArray) != 10 {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}

func TestHappyPathAvroGenArray2(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": {
					"type": "array",
					"items" : {
						"type" : "record",
						"name" : "ArrayObj",
						"fields" : [
							{ "name": "stringField", "type": "string" }
						]
					}
				}
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  3373,
		},
		GenerationRules: map[string]string{
			".testField.len()": "lenGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"lenGen": "{int}[0]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	generatedArray := res.(map[string]interface{})["testField"].([]interface{})
	if len(generatedArray) != 0 {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}

func TestHappyPathAvroGenArray(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": {
					"type": "array",
					"items" : "string",
					"default": []
				}
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  3373,
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	generatedArray := res.(map[string]interface{})["testField"].([]interface{})

	if generatedArray[0].(string) != "Bhz8K9EN0P" {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}
