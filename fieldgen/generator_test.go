package fieldgen

import (
	"fmt"
	"testing"
	"unicode"
)

func TestHappypathGenerate(t *testing.T) {
	const seed = 0
	res, _ := NewGenerator("{string}[res1]{1}")(seed)
	if res != "res1" {
		t.Fail()
	}
	res, _ = NewGenerator("{int}[132]{1}")(seed)
	if res != 132 {
		t.Fail()
	}
	res, _ = NewGenerator("{int}[132]{1}[231]{1}")(seed)
	if res != 132231 {
		t.Fail()
	}
	res, _ = NewGenerator("{string}[a-Z]{30}")(seed)
	if res != "UyizABeAsCmGcYwewHIgmAhUHCEecE" {
		t.Fail()
	}
	res, _ = NewGenerator("{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}")(seed)
	if res != "uyizabeasc@gcywewhigm.org" {
		t.Fail()
	}
	res, _ = NewGenerator("{string}[0|1|2|3]{10}")(seed)
	if res != "2133003032" {
		t.Fail()
	}
	res, _ = NewGenerator("{string}[a-z|A-Z|0-9|test]{30}")(seed)
	if res != "4Ytesttestabtestatest4test0CytestEW3iGtestatestu782testCe" {
		t.Fail()
	}
}

func TestProbability(t *testing.T) {
	// test that the probability of the or cases is the same
	// i.e. a-z|0-9 both the alphabetic values and the numeric values
	// have the same prob to be picked

	const seed = 0
	res, _ := NewGenerator("{string}[a-z|0-9|A-Z]{100000}")(seed)
	lower, digits, upper := 0, 0, 0
	for _, l := range res.(string) {
		if unicode.IsNumber(l) {
			digits += 1
		} else if unicode.IsLower(l) {
			lower += 1
		} else if unicode.IsUpper(l) {
			upper += 1
		}
	}
	if lower < 30000 || digits < 30000 || upper < 30000 {
		fmt.Println(lower, digits, upper)
		t.Fail()
	}
}

func TestNilIfInvalidPattern(t *testing.T) {
	if NewGenerator("{string}[]{1}") != nil {
		t.Fail()
	}
}
