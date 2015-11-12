### Fako

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
    Email     string `fako:"email"` //Notice the fako:"email" tag
    FullName  string `fako:"name"`
}

func main(){
  var user User
  fako.Fill(&user)

  fmt.Println(&user.Email)
  // This Prints AnthonyMeyer@Twimbo.biz
  // or another valid email
}
```

**Fako** currently supports the following data types for string fields:

- name
- username
- email
- phone
- password
- address

But our goal is to map all of the [fake](https://github.com/icrowley/fake) functions into **Fako**.


#### Credits
As you may have notices this is based on [fake](https://github.com/icrowley/fake) library, which does all the work to generate data.

#### Copyright
Copyright (c) 2015-2019 Wawandco SAS. See LICENSE for further details.
