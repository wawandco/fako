### Fako

[![Circle CI](https://circleci.com/gh/wawandco/fako.svg?style=svg)](https://circleci.com/gh/wawandco/fako) [![Godoc](https://img.shields.io/badge/godoc-docs-blue.svg)](https://godoc.org/github.com/wawandco/fako)
[![Go Report Card](https://goreportcard.com/badge/github.com/wawandco/fako)](https://goreportcard.com/report/github.com/wawandco/fako)

**Fako** is a library intended to fake Golang structs with fake but coherent data, **Fako** maps struct field tags and generates fake data accordingly. When no tag is specified, **Fako** can also automatically match field names with available faker functions.

We find it useful when writing specs to generate fake database data, hope you too.

#### Examples

##### Basic Usage with Tags

This is an example of how **Fako** works with explicit tags.

```go
import(
	"fmt"
	"github.com/wawandco/fako"
)

type User struct {
	Name     string `fako:"full_name"`
	Username string `fako:"user_name"`
	Email    string `fako:"email_address"`//Notice the fako:"email_address" tag
	Phone    string `fako:"phone"`
	Password string `fako:"simple_password"`
	Address  string `fako:"street_address"`
}

func main(){
	var user User
	fako.Fill(&user)

	fmt.Println(user.Email)
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

##### Field Name Matching (No Tags Required)

**Fako** can automatically match field names with faker functions when no `fako` tag is specified. This makes it even easier to generate fake data without needing to add tags to every field.

```go
import(
	"fmt"
	"github.com/wawandco/fako"
)

type User struct {
	FullName     string // Automatically matches "FullName" faker
	EmailAddress string // Automatically matches "EmailAddress" faker
	Phone        string // Automatically matches "Phone" faker
	UserName     string // Automatically matches "UserName" faker
	Company      string // Automatically matches "Company" faker
	
	// Field names with underscores are camelized for matching
	Street_Address string // Matches "StreetAddress" faker
	
	// Fields that don't match any faker remain empty
	CustomField string
	
	// Explicit tags take precedence over field name matching
	Name string `fako:"city"` // Uses "city" faker, not field name
}

func main(){
	var user User
	fako.Fill(&user)

	fmt.Printf("FullName: %s\n", user.FullName)
	fmt.Printf("EmailAddress: %s\n", user.EmailAddress)
	fmt.Printf("Phone: %s\n", user.Phone)
	fmt.Printf("Street_Address: %s\n", user.Street_Address)
	fmt.Printf("CustomField: '%s'\n", user.CustomField) // Empty
	fmt.Printf("Name: %s\n", user.Name) // Uses city faker due to tag
	
	// Field name matching also works with FillOnly and FillExcept
	var user2 User
	fako.FillOnly(&user2, "FullName", "EmailAddress")
	// Only FullName and EmailAddress will be filled
	
	var user3 User
	fako.FillExcept(&user3, "Phone")
	// All matching fields except Phone will be filled
}
```

**How Field Name Matching Works:**
- Field names are converted to CamelCase (e.g., `street_address` → `StreetAddress`)
- The camelized name is matched against available faker functions
- If a match is found, the corresponding faker is used
- If no match is found, the field remains empty
- Explicit `fako` tags always take precedence over field name matching

**Common Field Name Mappings:**
Here are some examples of field names that will automatically match faker functions:

```go
type Example struct {
    // Personal Information
    FirstName      string  // → fake.FirstName()
    LastName       string  // → fake.LastName()
    FullName       string  // → fake.FullName()
    UserName       string  // → fake.UserName()
    
    // Contact Information
    EmailAddress   string  // → fake.EmailAddress()
    Phone          string  // → fake.Phone()
    
    // Address Information
    StreetAddress  string  // → fake.StreetAddress()
    Street_Address string  // → fake.StreetAddress() (camelized)
    City           string  // → fake.City()
    State          string  // → fake.State()
    Country        string  // → fake.Country()
    
    // Business Information
    Company        string  // → fake.Company()
    JobTitle       string  // → fake.JobTitle()
    Industry       string  // → fake.Industry()
    
    // Other
    Color          string  // → fake.Color()
    Brand          string  // → fake.Brand()
    ProductName    string  // → fake.ProductName()
    Title          string  // → fake.Title()
    
    // These won't match (no corresponding faker)
    CustomField    string  // remains empty
    MySpecialData  string  // remains empty
}
```

Fako provides 3 built in functions `Fill`, `FillOnly`, and `FillExcept`, please go to [godoc](https://godoc.org/github.com/wawandco/fako) for details.

**Fako** supports most of the fields from the [fake](https://github.com/icrowley/fake) library. These can be used either as explicit tag values (e.g., `fako:"full_name"`) or as field names for automatic matching (e.g., `FullName string`).

**Available Faker Types:**
```
brand, character, characters, city, color, company, continent, country, credit_card_type, currency, currency_code, digits, domain_name, domain_zone, email_address, email_body, email_subject, female_first_name, female_full_name, female_full_name_with_prefix, female_full_name_with_suffix, female_last_name, female_patronymic, first_name, full_name, full_name_with_prefix, full_name_with_suffix, gender, gender_abbrev, hex_color, hex_color_short, ip_v4, industry, job_title, language, last_name, latitude_direction, longitude_direction, male_first_name, male_full_name, male_full_name_with_prefix, male_full_name_with_suffix, male_last_name, male_patronymic, model, month, month_short, paragraph, paragraphs, patronymic, phone, product, product_name, sentence, sentences, simple_password, state, state_abbrev, street, street_address, title, top_level_domain, user_name, week_day, week_day_short, word, words, zip
```

When using field name matching, convert the snake_case names above to CamelCase for your field names (e.g., `full_name` → `FullName`, `email_address` → `EmailAddress`).

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
	Username string `fako:"user_name"`
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
As you may have noticed this is based on [fake](https://github.com/icrowley/fake) library, which does all the work to generate data.

#### Copyright
Fako is Copyright © 2008-2015 Wawandco SAS. It is free software, and may be redistributed under the terms specified in the LICENSE file.
