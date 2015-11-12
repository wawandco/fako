package fako

import (
	"reflect"

	"github.com/wawandco/fako/Godeps/_workspace/src/github.com/icrowley/fake"
)

var typeMapping = map[string]func() string{
	"name":     fake.FullName,
	"email":    fake.EmailAddress,
	"phone":    fake.Phone,
	"username": fake.UserName,
	"password": fake.SimplePassword,
	"address":  fake.StreetAddress,
}

//Fill fill all the fields that have a fako: tag
func Fill(strukt interface{}) {
	fillWithDetails(strukt, []string{}, []string{})
}

//FillOnly fill fields that have a fako: tag and its name is on the second argument array
func FillOnly(strukt interface{}, fields ...string) {
	fillWithDetails(strukt, fields, []string{})
}

//FillExcept fill fields that have a fako: tag and its name is not on the second argument array
func FillExcept(strukt interface{}, fields ...string) {
	fillWithDetails(strukt, []string{}, fields)
}

func fillWithDetails(strukt interface{}, only []string, except []string) {
	elem := reflect.ValueOf(strukt).Elem()
	elemT := reflect.TypeOf(strukt).Elem()

	for i := 0; i < elem.NumField(); i++ {

		field := elem.Field(i)
		fieldt := elemT.Field(i)
		fakeType := fieldt.Tag.Get("fako")

		inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
		notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

		if field.CanSet() && fakeType != "" && inOnly && notInExcept {
			for kind, function := range typeMapping {
				if fakeType == kind {
					field.SetString(function())
					break
				}
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
