package fako

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/serenize/snaker"
)

//Fill fills all the fields that have a fako: tag
func Fill(elems ...interface{}) {
	for _, elem := range elems {
		FillElem(elem)
	}
}

//FillElem provides a way to fill a simple interface
func FillElem(strukt interface{}) {
	fillWithDetails(strukt, []string{}, []string{})
}

//FillOnly fills fields that have a fako: tag and its name is on the second argument array
func FillOnly(strukt interface{}, fields ...string) {
	fillWithDetails(strukt, fields, []string{})
}

//FillExcept fills fields that have a fako: tag and its name is not on the second argument array
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
		fakeTypes := fmt.Sprintf("%v", fieldt.Type)
		//fmt.Println("FROM FAKO", len(fakeTypes))
		if fakeTypes == "string" {
			if fakeType != "" {
				fakeType = snaker.SnakeToCamel(fakeType)
				function := findFakeFunctionFor(fakeType)

				inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
				notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

				if field.CanSet() && fakeType != "" && inOnly && notInExcept {
					field.SetString(function())
				}
			}
		} else {
			v, _ := strconv.Atoi(fakeType)
			fakeTypeInt := "Int"
			function := findFakeFunctionForInt(fakeTypeInt)
			switch fakeTypes {
			case "int", "int64":
				if field.CanSet() && fakeTypeInt != "" {
					field.SetInt(int64(function(v)))

				}

			case "uint", "uint8", "uint64":
				if field.CanSet() && fakeTypeInt != "" {
					field.SetUint(uint64(function(v)))
				}
			case "float64", "float32":
				// TODO
				if field.CanSet() && fakeTypeInt != "" {
					field.SetFloat(float64(function(v)))
				}

			}
		}

	}
}
