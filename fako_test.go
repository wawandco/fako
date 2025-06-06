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

func TestUsesTagValueNotFieldName(t *testing.T) {
	// Register a custom generator for the tag value "special_tag"
	Register("special_tag", func() string {
		return "FROM_TAG"
	})

	// Register a custom generator for the field name "FieldWithMismatch"
	Register("field_with_mismatch", func() string {
		return "FROM_FIELD_NAME"
	})

	// Create a struct where field name != tag value
	type TestStruct struct {
		FieldWithMismatch string `fako:"special_tag"`
	}

	var test TestStruct
	Fill(&test)

	// If the bug existed, it would use "FieldWithMismatch" (field name)
	// and return "FROM_FIELD_NAME". With the fix, it should use "special_tag" 
	// (tag value) and return "FROM_TAG"
	if test.FieldWithMismatch != "FROM_TAG" {
		t.Errorf("Expected field to be filled using tag value 'special_tag' returning 'FROM_TAG', but got: %s", test.FieldWithMismatch)
	}
}

func TestConcurrentRegisterAndFill(t *testing.T) {
	// Regression test to ensure no race conditions when registering custom generators
	// concurrently while filling structs. This test would fail before the mutex
	// protection was added to the customGenerators map.
	const numGoroutines = 100
	const numIterations = 10
	
	type TestStruct struct {
		Field1 string `fako:"test_gen_1"`
		Field2 string `fako:"test_gen_2"`
	}
	
	done := make(chan bool, numGoroutines*2)
	
	// Goroutines registering generators
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numIterations; j++ {
				Register("test_gen_1", func() string {
					return "value1"
				})
				Register("test_gen_2", func() string {
					return "value2"
				})
			}
			done <- true
		}(i)
	}
	
	// Goroutines filling structs
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				var test TestStruct
				Fill(&test)
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines*2; i++ {
		<-done
	}
}

func TestConcurrentFuzz(t *testing.T) {
	// Regression test to ensure no race conditions in the Fuzz function.
	// This test verifies that concurrent calls to Fuzz() are safe and that
	// random generation doesn't cause data races.
	const numGoroutines = 100
	const numIterations = 10
	
	type TestStruct struct {
		Value string
	}
	
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				var test TestStruct
				Fuzz(&test)
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
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

func TestFieldNameMatching(t *testing.T) {
	// Test struct with field names that match faker functions
	type AutoMatchStruct struct {
		FullName      string // Should match "FullName" faker
		EmailAddress  string // Should match "EmailAddress" faker
		Phone         string // Should match "Phone" faker
		UserName      string // Should match "UserName" faker
		Company       string // Should match "Company" faker
		NoMatchField  string // Should not be filled (no matching faker)
		TaggedField   string `fako:"city"` // Should use tag, not field name
	}

	var test AutoMatchStruct
	Fill(&test)

	// Check that fields with matching faker functions are filled
	if test.FullName == "" {
		t.Errorf("Expected FullName to be filled by field name matching, but got empty string")
	}
	if test.EmailAddress == "" {
		t.Errorf("Expected EmailAddress to be filled by field name matching, but got empty string")
	}
	if !strings.Contains(test.EmailAddress, "@") {
		t.Errorf("Expected EmailAddress to contain '@', but got: %s", test.EmailAddress)
	}
	if test.Phone == "" {
		t.Errorf("Expected Phone to be filled by field name matching, but got empty string")
	}
	if test.UserName == "" {
		t.Errorf("Expected UserName to be filled by field name matching, but got empty string")
	}
	if test.Company == "" {
		t.Errorf("Expected Company to be filled by field name matching, but got empty string")
	}

	// Check that field without matching faker is not filled
	if test.NoMatchField != "" {
		t.Errorf("Expected NoMatchField to remain empty (no matching faker), but got: %s", test.NoMatchField)
	}

	// Check that tagged field uses tag value, not field name
	if test.TaggedField == "" {
		t.Errorf("Expected TaggedField to be filled using tag value 'city', but got empty string")
	}
}

func TestFieldNameMatchingWithVariousFormats(t *testing.T) {
	// Test field names in different formats that should still match
	type VariousFormatsStruct struct {
		FirstName      string // Should match "FirstName" faker
		First_Name     string // Should match "FirstName" faker (camelized)
		LastName       string // Should match "LastName" faker
		Street_Address string // Should match "StreetAddress" faker (camelized)
	}

	var test VariousFormatsStruct
	Fill(&test)

	if test.FirstName == "" {
		t.Errorf("Expected FirstName to be filled, but got empty string")
	}
	if test.First_Name == "" {
		t.Errorf("Expected First_Name to be filled by camelization matching, but got empty string")
	}
	if test.LastName == "" {
		t.Errorf("Expected LastName to be filled, but got empty string")
	}
	if test.Street_Address == "" {
		t.Errorf("Expected Street_Address to be filled by camelization matching, but got empty string")
	}
}

func TestFieldNameMatchingWithFillOnly(t *testing.T) {
	// Test that field name matching works with FillOnly
	type TestStruct struct {
		FullName     string // Should match "FullName" faker
		EmailAddress string // Should match "EmailAddress" faker
		Phone        string // Should match "Phone" faker
		NoMatchField string // Should not be filled (no matching faker)
	}

	var test TestStruct
	FillOnly(&test, "FullName", "Phone")

	// Only FullName and Phone should be filled
	if test.FullName == "" {
		t.Errorf("Expected FullName to be filled by FillOnly with field name matching, but got empty string")
	}
	if test.Phone == "" {
		t.Errorf("Expected Phone to be filled by FillOnly with field name matching, but got empty string")
	}
	
	// EmailAddress should not be filled (not in FillOnly list)
	if test.EmailAddress != "" {
		t.Errorf("Expected EmailAddress to be empty (not in FillOnly list), but got: %s", test.EmailAddress)
	}
	
	// NoMatchField should not be filled (no matching faker)
	if test.NoMatchField != "" {
		t.Errorf("Expected NoMatchField to be empty (no matching faker), but got: %s", test.NoMatchField)
	}
}

func TestFieldNameMatchingWithFillExcept(t *testing.T) {
	// Test that field name matching works with FillExcept
	type TestStruct struct {
		FullName     string // Should match "FullName" faker
		EmailAddress string // Should match "EmailAddress" faker
		Phone        string // Should match "Phone" faker
		NoMatchField string // Should not be filled (no matching faker)
	}

	var test TestStruct
	FillExcept(&test, "EmailAddress")

	// FullName and Phone should be filled (not in except list)
	if test.FullName == "" {
		t.Errorf("Expected FullName to be filled by FillExcept with field name matching, but got empty string")
	}
	if test.Phone == "" {
		t.Errorf("Expected Phone to be filled by FillExcept with field name matching, but got empty string")
	}
	
	// EmailAddress should not be filled (in except list)
	if test.EmailAddress != "" {
		t.Errorf("Expected EmailAddress to be empty (in FillExcept list), but got: %s", test.EmailAddress)
	}
	
	// NoMatchField should not be filled (no matching faker)
	if test.NoMatchField != "" {
		t.Errorf("Expected NoMatchField to be empty (no matching faker), but got: %s", test.NoMatchField)
	}
}