package fako_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawandco/fako"
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
	fako.Fill(&myCase)
	var first = myCase.Name

	assert.NotEqual(t, myCase.Name, "")
	fako.Fill(&myCase)
	assert.NotEqual(t, myCase.Name, first)
}

func TestFillsEmail(t *testing.T) {
	var myCase OU
	fako.Fill(&myCase)

	assert.True(t, strings.Contains(myCase.Email, "@"))
}

func TestFillOnly(t *testing.T) {
	var myCase OU
	fako.FillOnly(&myCase, "Name")

	assert.Equal(t, myCase.Email, "")
}

func TestFillExcept(t *testing.T) {
	var myCase OU
	fako.FillExcept(&myCase, "Name")

	assert.Equal(t, myCase.Name, "")
	assert.NotEqual(t, myCase.Email, "")
}

func TestCustomGenerator(t *testing.T) {
	fako.Register("a_gen", func() string {
		return "A"
	})

	var myCase OU
	fako.Fill(&myCase)

	assert.Equal(t, myCase.AValue, "A")
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
	fako.Fuzz(&simple)

	assert.NotEqual(t, 0, simple.Value)
	assert.NotEqual(t, 0, simple.Value32)
	assert.NotEqual(t, 0, simple.Value64)
	assert.NotEqual(t, 0.0, simple.ValueF32)
	assert.NotEqual(t, 0.0, simple.ValueF64)

	count := 0
	for i := 0; i < 10; i++ {
		fako.Fuzz(&simple)
		if simple.Active {
			count++
		}
	}

	assert.True(t, count > 0)
}

func TestFillAll(t *testing.T) {
	customName := struct {
		Name string `fako:"full_name"`
	}{}

	customU := struct {
		Username string `fako:"user_name"`
	}{}

	fako.Fill(&customName, &customU)
	assert.NotZero(t, customName.Name)
	assert.NotZero(t, customU.Username)
}

func TestFillAliasType(t *testing.T) {
	type DateTime string

	type Custom struct {
		Name DateTime `fako:"full_name"`
	}

	el := &Custom{}
	fako.Fill(el)

	assert.NotZero(t, el.Name)
}
