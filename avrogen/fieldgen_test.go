package avrogen

import (
	"fmt"
	"math/rand"
	"testing"
	"unicode"

	"github.com/google/uuid"
)

func TestHappyPathGenerateString(t *testing.T) {
	random := rand.New(rand.NewSource(0))
	expected := "res1"
	res, _ := newFieldGen("{string}[res1]{1}", random)()
	if res != expected {
		t.Errorf("expected %s, received %d", expected, res)
	}

	res, _ = newFieldGen("{string}[a-Z]{30}", random)()
	if res != "yizABeAsCmGcYwewHIgmAhUHCEecEP" {
		t.Fail()
	}
	res, _ = newFieldGen("{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}", random)()
	if res != "mpotqplyjo@fkueeaobht.org" {
		t.Fail()
	}
	res, _ = newFieldGen("{string}[a-z|A-Z|0-9|test]{30}", random)()
	if res != "test2mptesttestttest8testtestx9test2test09LNf7test0test39g61" {
		t.Fail()
	}
	res, _ = newFieldGen("{boolean}[false|true]{1}", random)()
	if res != false {
		t.Fail()
	}
}

func TestHappyPathGenerateFloat(t *testing.T) {
	random := rand.New(rand.NewSource(0))
	var expected float32 = 0.678
	res, _ := newFieldGen("{float}[0]{1}[.]{1}[ 0-9 ]{3}", random)()
	if res.(float32)-expected > 0.000001 {
		t.Errorf("expected %f, received %f", expected, res)
	}
}

func TestHappyPathGenerateUUID(t *testing.T) {
	random := rand.New(rand.NewSource(0))
	res, _ := newFieldGen("{string}[uuid()]{1}", random)()
	_, err := uuid.Parse(res.(string))
	if err != nil {
		t.Fail()
	}
}

func TestHappyPathGenerateInt(t *testing.T) {
	random := rand.New(rand.NewSource(0))
	expected := 132
	res, _ := newFieldGen("{int}[132]{1}", random)()
	if res != expected {
		t.Errorf("expected %d, received %d", expected, res)
	}
	expected = 132231
	res, _ = newFieldGen("{int}[132]{1}[231]{1}", random)()
	if res != expected {
		t.Errorf("expected %d, received %d", expected, res)
	}
	expected = 911
	res, _ = newFieldGen("{int}[1|9]{3}", random)()
	if res != expected {
		t.Errorf("expected %d, received %d", expected, res)
	}
	var expectedL int64 = 303232103112013030
	res, _ = newFieldGen("{long}[0|1|2|3]{18}", random)()
	if res != expectedL {
		t.Errorf("expected %d, received %d", expectedL, res)
	}
}

func TestProbability(t *testing.T) {
	// test that the probability of the `or` clauses is the same
	// i.e. a-z|0-9 both the alphabetic values and the numeric values
	// have the same prob to be picked
	random := rand.New(rand.NewSource(0))
	res, _ := newFieldGen("{string}[a-z|0-9|A-Z]{100000}", random)()
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
	random := rand.New(rand.NewSource(0))
	if newFieldGen("{string}[]{1}", random) != nil {
		t.Fail()
	}
}
