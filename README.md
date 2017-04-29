
### Adds support for ints and floats. However no support for decimals
```go
package main

import (
	"encoding/json"
	"io/ioutil"
	//"github.com/wawandco/fako"
	"github.com/MintyOwl/fako"

	"github.com/hokaccha/go-prettyjson"
)

type Person struct {
	Phone    int64  `fako:"10";json:"phone"` // Just add number of digits you want. In this example its 10
	Name     string `fako:"full_name";json:"name"`
	Username string `fako:"user_name";json:"user_name"`
	Email    string `fako:"email_address";json:"email"`
	Password string `fako:"simple_password";json:"password"`
	Address  string `fako:"street_address";json:"address"`
}

var persons []*Person
var per []Person

func init() {
	//GetPersons()
}

func formatter() *prettyjson.Formatter {
	return &prettyjson.Formatter{
		StringMaxLength: 0,
		DisabledColor:   true,
		Indent:          2,
	}
}

func CreatePersons() {
	for i := 0; i < 5; i++ {
		var person Person
		fako.FillElem(&person)
		persons = append(persons, &person)
	}
	b, _ := json.Marshal(persons)
	b, _ = formatter().Format(b)
	//fmt.Println("JSON ERR >", err, string(b))
	ioutil.WriteFile("people.json", b, 0666)
}

func GetPersons() []*Person {
	if persons != nil {
		return persons
	}
	jsn, _ := ioutil.ReadFile("people.json")
	var personns = make([]*Person, 0)
	json.Unmarshal(jsn, &personns)
	persons = personns
	return personns
}

func main() {	
	CreatePersons()
}

```

### Fako

[![Circle CI](https://circleci.com/gh/wawandco/fako.svg?style=svg)](https://circleci.com/gh/wawandco/fako) [![Godoc](https://img.shields.io/badge/godoc-docs-blue.svg)](https://godoc.org/github.com/wawandco/fako)
[![Go Report Card](https://goreportcard.com/badge/github.com/wawandco/fako)](https://goreportcard.com/report/github.com/wawandco/fako)

**Fako** is a library intended to fake Golang structs with fake but coherent data, **Fako** maps struct field tags and generates fake data accordingly.

We find it useful when writing specs to generate fake database data, hope you too.

#### Example

This is an example of how **Fako** works.

```go
import(
  "fmt"
  "github.com/wawandco/fako"
)

type User struct {
    Name     string `fako:"full_name"`
  	Username string `fako:"username"`
  	Email    string `fako:"email_address"`//Notice the fako:"email_address" tag
  	Phone    string `fako:"phone"`
  	Password string `fako:"simple_password"`
  	Address  string `fako:"street_address"`
}

func main(){
  var user User
  fako.Fill(&user)

  fmt.Println(&user.Email)
  // This prints something like AnthonyMeyer@Twimbo.biz
  // or another valid email

  var userWithOnlyEmail User
  fako.FillOnly(&userWithOnlyEmail, "Email")
  //This will fill all only the email

  var userWithoutEmail User
  fako.FillExcept(&userWithoutEmail, "Email")
  //This will fill all the fields except the email

}
```

Fako provides 3 built in functions `Fill`, `FillOnly`, and `FillExcept`, please go to [godoc](https://godoc.org/github.com/wawandco/fako) for details.

**Fako** support most of the fields on the [fake](https://github.com/icrowley/fake)  library, below you can see a list of the field types you can use.

- brand
- character
- characters
- city
- color
- company
- continent
- country
- credit_card_type
- currency
- currency_code
- digits
- domain_name
- domain_zone
- email_address
- email_body
- email_subject
- female_first_name
- female_full_name
- female_full_name_with_prefix
- female_full_name_with_suffix
- female_last_name
- female_patronymic
- first_name
- full_name
- full_name_with_prefix
- full_name_with_suffix
- gender
- gender_abbrev
- hex_color
- hex_color_short
- ip_v4
- industry
- job_title
- language
- last_name
- latitude_direction
- longitude_direction
- male_first_name
- male_full_name
- male_full_name_with_prefix
- male_full_name_with_suffix
- male_last_name
- male_patronymic
- model
- month
- month_short
- paragraph
- paragraphs
- patronymic
- phone
- product
- product_name
- sentence
- sentences
- simple_password
- state
- state_abbrev
- street
- street_address
- title
- top_level_domain
- user_name
- week_day
- week_day_short
- word
- words
- zip

#### Custom Generators

Fako provides a function called `Register` to add custom data generators in case you need something that our provided generators cannot cover.

To add a custom generator simply call the `Register` function as in the following example:

```go
import(
  "fmt"
  "github.com/wawandco/fako"
)

type User struct {
    Name     string `fako:"full_name"`
    Username string `fako:"username"`
    Email    string `fako:"email_address"`//Notice the fako:"email_address" tag
    Phone    string `fako:"phone"`
    Password string `fako:"simple_password"`
    Address  string `fako:"street_address"`

    AValue   string `fako:"a_gen"`
}

func main(){
  fako.Register("a_gen", func() string {
    return "My Value"
  })

  var user User
  fako.Fill(&user)
  fmt.Println(user.AValue) //should print My Value
}
```

When using custom generators please keep the following in mind:
  1. Call Register function before calling `Fill` and its brothers.
  2. Custom generators override base generators, if you pick the same name as one of the existing generators, we will override the existing generator with yours.


#### Fuzzing

Sometimes you just want to generate random data inside a struct, for those cases you wont want to fill fako types (yes, we understand that part). Fako provides you a `Fuzz` function you can use to fuzz your structs with random data that simply matches the struct's field types.

You can use it as in the following example:

```go
import "fako"

type Instance struct {
   Name string
   Number int
}

func main(){
  instance := Instance{}
  fako.Fuzz(&instance) // This fills your instance variable
}
```

Note, Fuzz function works for the following types `string, bool, int, int32, int64, float32, float64`. other types like `Array` or `Chan` are out of our scope.

#### Credits
As you may have notices this is based on [fake](https://github.com/icrowley/fake) library, which does all the work to generate data.

#### Copyright
Fako is Copyright © 2008-2015 Wawandco SAS. It is free software, and may be redistributed under the terms specified in the LICENSE file.
