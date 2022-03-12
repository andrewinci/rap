package fieldgen

import "testing"

func TestHappypathGenerate(t *testing.T) {
	const seed = 0
	if NewGenerator("{string}[res1]{1}")(seed) != "res1" {
		t.Fail()
	}
	if NewGenerator("{int}[132]{1}")(seed) != 132 {
		t.Fail()
	}
	if NewGenerator("{int}[132]{1}[231]{1}")(seed) != 132231 {
		t.Fail()
	}
	if NewGenerator("{string}[a-Z]{30}")(seed) != "CuByHIzZkAkbLEEanSOCPMigNCkYrw" {
		t.Fail()
	}
	if NewGenerator("{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}")(seed) != "cubyhizzka@bleeansocp.org" {
		t.Fail()
	}
	if NewGenerator("{string}[0|1|2|3]{10}")(seed) != "2212303100" {
		t.Fail()
	}
}
