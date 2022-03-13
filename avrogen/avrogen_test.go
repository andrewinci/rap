package avrogen

import (
	"testing"
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
			{ "name": "stringField", "type": "string" }
		]
	 }
	`
	sut, err := NewAvroGen(AvroGenConfiguration{Schema: testSchema}, 0)
	if err != nil {
		t.FailNow()
	}
	res1, err := sut.Generate()
	if err != nil {
		t.FailNow()
	}
	res2, err := sut.Generate()
	if err != nil {
		t.FailNow()
	}
	if string(res1) == string(res2) {
		t.Error("generated 2 identical records")
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
	sut, err := NewAvroGen(AvroGenConfiguration{
		Schema: testSchema,
		GenerationRules: map[string]string{
			"boolean": "booleanGen",
			"int":     "intGen",
			"long":    "longGen",
			"float":   "floatGen",
			"double":  "doubleGen",
			"string":  "stringGen",
		},
		// all generators are constants
		Generators: map[string]string{
			"booleanGen": "{boolean}[false]{1}",
			"intGen":     "{int}[1321]{1}",
			"longGen":    "{long}[9876]{1}",
			"floatGen":   "{float}[12.12]{1}",
			"doubleGen":  "{double}[123.321]{1}",
			"stringGen":  "{string}[stringValue]{1}",
		}}, 0)
	if err != nil {
		t.FailNow()
	}
	res1, _ := sut.Generate()
	res2, _ := sut.Generate()
	if string(res1) != string(res2) {
		t.Error("the records should be identical")
	}
}
