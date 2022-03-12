package fieldgen

import "testing"

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
