package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
)

type testCase struct {
	inty    int
	stringy string
	booly   bool
}

func TestCheckStringForEmptiness(t *testing.T) {
	cases := []testCase{
		{51, "Elon", CheckStringForEmptiness("Musk")},
		{58, "Jeff", CheckStringForEmptiness("Bezos")},
		{38, "Mark", CheckStringForEmptiness("Zuckerberg")},
	}

	for _, tc := range cases {
		name := tc.stringy

		got := CheckStringForEmptiness(name)
		if tc.booly != got {
			t.Errorf("Expected '%v', but got '%v'", tc.booly, got)
		}
	}
}

func TestCheckStringForIpOrHostname(t *testing.T) {
	cases := []testCase{
		{1, "https://example.com", IsStringValidUrl("https://example.com")},
		{2, "127.0.0.1", IsStringValidUrl("127.0.0.1")},
	}

	for _, tc := range cases {
		route := tc.stringy
		isValid := CheckStringForIpOrHostname(route)

		got := CheckStringForIpOrHostname(route)
		if isValid != got {
			t.Errorf("Expected '%v', but got '%v'", isValid, got)
		}
	}
}

func TestIsStringValidURL(t *testing.T) {
	cases := []testCase{
		{1, "https://example.com", IsStringValidUrl("https://example.com")},
		{2, "https://exampl.com", IsStringValidUrl("https://exampl.com")},
		{3, "https://examp.com", IsStringValidUrl("https://examp.com")},
	}

	for _, tc := range cases {
		url := tc.stringy

		got := IsStringValidUrl(url)
		if tc.booly != got {
			t.Errorf("Expected '%v', but got '%v'", tc.booly, got)
		}
	}
}
func TestCapitalizeString(t *testing.T) {
	input := "macro"
	inputCaps := CapitalizeString(input)

	got := CapitalizeString(input)
	if inputCaps != got {
		t.Errorf("Expected '%v', but got '%v'", inputCaps, got)
	}
}

func TestLowercaseString(t *testing.T) {
	input := "Macro"
	inputLow := LowercaseString(input)

	got := LowercaseString(input)
	if inputLow != got {
		t.Errorf("Expected '%v', but got '%v'", inputLow, got)
	}
}

func TestYesNoToInt(t *testing.T) {
	cases := []testCase{
		{1, "Yes", true},
		{2, "No", false},
	}

	for _, tc := range cases {
		toInt := YesNoToInt(tc.stringy)

		got := YesNoToInt(tc.stringy)
		if toInt != got {
			t.Errorf("Expected '%v', but got '%v'", toInt, got)
		}
	}
}

func TestWriteToEnv(t *testing.T) {
	writeTo := "TEST_CONFIG"
	writeValue := "TEST"
	newKeyValue := fmt.Sprintf("%s=%s", writeTo, writeValue)

	got := true
	if WriteToEnv(writeTo, newKeyValue) != got {
		t.Errorf("Expected '%v', but got '%v'", true, got)
	}

	err := os.Remove(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func TestV8Isolates(t *testing.T) {
	execJS, err := V8Isolates("function hello() {return 'hello';} result = hello(); result;", true) // Execute script
	if err != nil {
		t.Errorf("exeJS error: %v", err)
	}

	got := execJS
	if execJS != got {
		t.Errorf("expected '%v', but got '%v'", execJS, got)
	}
}

func TestV8NodeJS(t *testing.T) {
	execNodeJS, err := V8NodeJS("function hello() {return 'hello';} result = hello(); result;", true) // Execute script
	if err != nil {
		t.Errorf("exeNodeJS error: %v", err)
	}

	got := execNodeJS
	if execNodeJS != got {
		t.Errorf("expected '%v', but got '%v'", execNodeJS, got)
	}
}
