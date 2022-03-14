package avrogen

import (
	"testing"

	"github.com/andrewinci/rap/configuration"
)

func TestHappyPathAvroGen(t *testing.T) {
	testSchema := `
	{
		"type" : "record",
		"name" : "Example",
		"fields" : [
			{ "name": "booleanField", "type": "boolean" },
			{ "name": "intField", "type": "int" },
			{ "name": "longField", "type": "long" },
			{ "name": "floatField", "type": "float" },
			{ "name": "doubleField", "type": "double" },
			{ "name": "stringField", "type": "string" },
			{ "name": "nullField", "type": "null" }
		]
	 }
	`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res1, key1, err := sut.Generate()
	if err != nil {
		t.FailNow()
	}
	res2, key2, err := sut.Generate()
	if err != nil {
		t.FailNow()
	}
	if string(res1) == string(res2) {
		t.Error("generated 2 identical records")
	}
	if string(key1) == string(key2) {
		t.Error("generated 2 identical keys")
	}
}

func TestHappyPath2AvroGen(t *testing.T) {
	testSchema := `
	{
		"type" : "record",
		"name" : "Example",
		"fields" : [
			{ "name": "booleanField",   "type": "boolean" },
			{ "name": "intField", 		"type": "int" },
			{ "name": "longField", 		"type": "long" },
			{ "name": "floatField", 	"type": "float" },
			{ "name": "doubleField", 	"type": "double" },
			{ "name": "stringField", 	"type": "string" }
		]
	 }
	`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		},
		GenerationRules: map[string]string{
			"boolean": "booleanGen",
			"int":     "intGen",
			"long":    "longGen",
			"float":   "floatGen",
			"double":  "doubleGen",
			"string":  "stringGen",
			"key":     "keyGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"booleanGen": "{boolean}[false]{1}",
			"intGen":     "{int}[1321]{1}",
			"longGen":    "{long}[9876]{1}",
			"floatGen":   "{float}[12.12]{1}",
			"doubleGen":  "{double}[123.321]{1}",
			"stringGen":  "{string}[stringValue]{1}",
			"keyGen":     "{string}[fixed-key]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res1, key1, _ := sut.Generate()
	res2, key2, _ := sut.Generate()
	if string(res1) != string(res2) {
		t.Error("the records should be identical")
	}
	if key1 != key2 {
		t.Error("the record keys should be identical")
	}
}

func TestHappyAvroGenUnion(t *testing.T) {
	testSchema := `
	{
		"type" : "record",
		"name" : "Example",
		"fields" : [
			{ "name": "testField", "type": ["boolean", "null"] }
		]
	 }
	`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	rawRes, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	res := rawRes.(map[string]interface{})["testField"]
	if res != true {
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
	if res.(map[string]interface{})["testField"].(map[string]interface{})["testNestedField"] != 45678923 {
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
		}}, 1)
	if err != nil {
		t.FailNow()
	}
	_, err = sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
}

func TestHappyAvroGenEnum(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": {
					"type": "enum",
					"name": "Suit",
					"symbols" : ["SPADES", "HEARTS", "DIAMONDS", "CLUBS"]
				}
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
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	if res.(map[string]interface{})["testField"] != "HEARTS" {
		t.FailNow()
	}
}

func TestHappyAvroGenEnum2(t *testing.T) {
	testSchema := `
	{
		"type": "record",
		"name": "Example",
		"fields": [
			{
				"name": "testField",
				"type": {
					"type": "enum",
					"name": "Suit",
					"symbols" : ["SPADES", "HEARTS", "DIAMONDS", "CLUBS"]
				}
			}
		]
	}`
	sut, err := NewAvroGen(configuration.AvroGenConfiguration{
		Schema: configuration.SchemaConfiguration{
			Raw: testSchema,
			Id:  1,
		},
		GenerationRules: map[string]string{
			".testField": "stringGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"stringGen": "{string}[DIAMONDS|CLUBS]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res, err := sut.generate(sut.getSchema(), "")
	if err != nil {
		t.FailNow()
	}
	if res.(map[string]interface{})["testField"] != "DIAMONDS" {
		t.FailNow()
	}
	_, _, err = sut.Generate()
	if err != nil {
		t.FailNow()
	}
}
