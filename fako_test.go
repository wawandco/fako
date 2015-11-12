package fako

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wawandco/fako/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

type OU struct {
	Name     string `fako:"name"`
	Username string `fako:"username"`
	Email    string `fako:"email"`
	Phone    string `fako:"phone"`
	Password string `fako:"password"`
	Address  string `fako:"address"`
}

func TestFillsFields(t *testing.T) {
	var myCase OU
	Fill(&myCase)
	var first = myCase.Name
	fmt.Println(myCase.Email)
	assert.NotEqual(t, myCase.Name, "")
	Fill(&myCase)
	assert.NotEqual(t, myCase.Name, first)
}

func TestFillsEmail(t *testing.T) {
	var myCase OU
	Fill(&myCase)

	assert.True(t, strings.Contains(myCase.Email, "@"))
}

func TestFillOnly(t *testing.T) {
	var myCase OU
	FillOnly(&myCase, "Name")

	assert.Equal(t, myCase.Email, "")
}

func TestFillExcept(t *testing.T) {
	var myCase OU
	FillExcept(&myCase, "Name")

	assert.Equal(t, myCase.Name, "")
	assert.NotEqual(t, myCase.Email, "")
}
