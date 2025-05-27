package fako

import (
	"strings"
	"testing"
)

type OU struct {
	Name     string `fako:"full_name"`
	Username string `fako:"user_name"`
	Email    string `fako:"email_address"`
	Phone    string `fako:"phone"`
	Password string `fako:"simple_password"`
	Address  string `fako:"street_address"`
	IgnoreMe string ``
	AValue   string `fako:"a_gen"`
}

func TestFillsFields(t *testing.T) {
	var myCase OU
	Fill(&myCase)
	var first = myCase.Name

	if myCase.Name == "" {
		t.Errorf("Expected Name to be filled, but got empty string")
	}
	Fill(&myCase)
	if myCase.Name == first {
		t.Errorf("Expected Name to change on second fill, but got same value: %s", first)
	}
}

func TestFillsEmail(t *testing.T) {
	var myCase OU
	Fill(&myCase)

	if !strings.Contains(myCase.Email, "@") {
		t.Errorf("Expected Email to contain '@', but got: %s", myCase.Email)
	}
}

func TestFillOnly(t *testing.T) {
	var myCase OU
	FillOnly(&myCase, "Name")

	if myCase.Email != "" {
		t.Errorf("Expected Email to be empty when filling only Name, but got: %s", myCase.Email)
	}
}

func TestFillExcept(t *testing.T) {
	var myCase OU
	FillExcept(&myCase, "Name")

	if myCase.Name != "" {
		t.Errorf("Expected Name to be empty when excluded from fill, but got: %s", myCase.Name)
	}
	if myCase.Email == "" {
		t.Errorf("Expected Email to be filled when not excluded, but got empty string")
	}
}

func TestCustomGenerator(t *testing.T) {
	Register("a_gen", func() string {
		return "A"
	})

	var myCase OU
	Fill(&myCase)

	if myCase.AValue != "A" {
		t.Errorf("Expected AValue to be 'A' from custom generator, but got: %s", myCase.AValue)
	}
}

type Simple struct {
	Attribute string
	Value     int
	Value32   int32
	Value64   int64
	ValueF32  float32
	ValueF64  float64
	Active    bool
}

func TestFuzz(t *testing.T) {
	simple := Simple{}
	Fuzz(&simple)

	if simple.Value == 0 {
		t.Errorf("Expected Value to be non-zero, but got: %d", simple.Value)
	}
	if simple.Value32 == 0 {
		t.Errorf("Expected Value32 to be non-zero, but got: %d", simple.Value32)
	}
	if simple.Value64 == 0 {
		t.Errorf("Expected Value64 to be non-zero, but got: %d", simple.Value64)
	}
	if simple.ValueF32 == 0.0 {
		t.Errorf("Expected ValueF32 to be non-zero, but got: %f", simple.ValueF32)
	}
	if simple.ValueF64 == 0.0 {
		t.Errorf("Expected ValueF64 to be non-zero, but got: %f", simple.ValueF64)
	}

	count := 0
	for i := 0; i < 10; i++ {
		Fuzz(&simple)
		if simple.Active {
			count++
		}
	}

	if count == 0 {
		t.Errorf("Expected at least one true value for Active in 10 iterations, but got none")
	}
}

func TestFillAll(t *testing.T) {
	customName := struct {
		Name string `fako:"full_name"`
	}{}

	customU := struct {
		Username string `fako:"user_name"`
	}{}

	Fill(&customName, &customU)
	if customName.Name == "" {
		t.Errorf("Expected customName.Name to be filled, but got empty string")
	}
	if customU.Username == "" {
		t.Errorf("Expected customU.Username to be filled, but got empty string")
	}
}

func TestFillAliasType(t *testing.T) {
	type DateTime string

	type Custom struct {
		Name DateTime `fako:"full_name"`
	}

	el := &Custom{}
	Fill(el)

	if el.Name == "" {
		t.Errorf("Expected Name to be filled, but got empty string")
	}
}

func TestCamelize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"hello_world_test", "HelloWorldTest"},
		{"hello-world-test", "HelloWorldTest"},
		{"hello world test", "HelloWorldTest"},
		{"", ""},
		{"a", "A"},
		{"user_name", "UserName"},
		{"full_name", "FullName"},
		{"email_address", "EmailAddress"},
		{"simple_password", "SimplePassword"},
		{"street_address", "StreetAddress"},
		{"API_KEY", "ApiKey"},
		{"test_API_key", "TestApiKey"},
		{"multiple___underscores", "MultipleUnderscores"},
		{"multiple---dashes", "MultipleDashes"},
		{"multiple   spaces", "MultipleSpaces"},
		{"mixed_-_delimiters test", "MixedDelimitersTest"},
	}

	for _, test := range tests {
		result := camelize(test.input)
		if result != test.expected {
			t.Errorf("camelize(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}