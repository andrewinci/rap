package avrogen

import (
	"testing"
	"time"
)

func TestParseInvalidPattern(t *testing.T) {
	// parse pattern with invalid type
	if parsePattern("{asdf}[a]{1}") != nil {
		t.Fail()
	}
	// parse pattern with invalid content
	if parsePattern("{asdf}[]{1}") != nil {
		t.Fail()
	}
	// parse pattern with invalid count
	if parsePattern("{asdf}[a]{1a}") != nil {
		t.Fail()
	}
}

func TestTrimOrClauses(t *testing.T) {
	// the first option should be all the letters
	option1 := parsePattern("{string}[ a-Z | 0 ]{1}").content[0].options[0]
	if len(option1()) == 1 || option1()[0][0] == ' ' {
		t.Fail()
	}
}

func TestParseUUIDFunction(t *testing.T) {
	uuidGen := parsePattern("{string}[ uuid() ]{1}").content[0].options[0]
	uuid1 := uuidGen()
	uuid2 := uuidGen()
	if len(uuid1) != 1 || len(uuid2) != 1 {
		// only 1 uuid should be available
		t.Fail()
	}
	if uuid1[0] == uuid2[0] {
		// any call should generate a new uuid
		t.Fail()
	}
}

func TestParseTimestampFunction(t *testing.T) {
	timestampGen := parsePattern("{string}[ timestamp_ms() ]{1}").content[0].options[0]
	time1 := timestampGen()
	time.Sleep(1 * time.Millisecond)
	time2 := timestampGen()
	if len(time1) != 1 || len(time2) != 1 {
		// only 1 uuid should be available
		t.Fail()
	}
	if time1[0] == time2[0] {
		// any call should generate a new uuid
		t.Fail()
	}
}
