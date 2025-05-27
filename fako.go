package fako

import (
	"reflect"
)

//Fill fills all the fields that have a fako: tag
func Fill(elems ...any) {
	for _, elem := range elems {
		FillElem(elem)
	}
}

//FillElem provides a way to fill a simple interface
func FillElem(strukt any) {
	fillWithDetails(strukt, []string{}, []string{})
}

//FillOnly fills fields that have a fako: tag and its name is on the second argument array
func FillOnly(strukt any, fields ...string) {
	fillWithDetails(strukt, fields, []string{})
}

//FillExcept fills fields that have a fako: tag and its name is not on the second argument array
func FillExcept(strukt any, fields ...string) {
	fillWithDetails(strukt, []string{}, fields)
}

func fillWithDetails(strukt any, only []string, except []string) {
	elem := reflect.ValueOf(strukt).Elem()
	elemT := reflect.TypeOf(strukt).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldt := elemT.Field(i)
		fakeType := fieldt.Tag.Get("fako")

		if fakeType != "" {
			fakeType = camelize(fakeType)
			function := findFakeFunctionFor(fakeType)

			inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
			notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

			if field.CanSet() && fakeType != "" && inOnly && notInExcept {
				field.SetString(function())
			}
		}

	}
}
