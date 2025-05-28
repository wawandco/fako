package fako

import (
	"reflect"
)

// Fill fills all the fields that have a fako: tag
// This function is safe for concurrent use.
func Fill(elems ...any) {
	for _, elem := range elems {
		FillElem(elem)
	}
}

// FillElem provides a way to fill a simple interface
// This function is safe for concurrent use.
func FillElem(strukt any) {
	fillWithDetails(strukt, []string{}, []string{})
}

// FillOnly fills fields that have a fako: tag and its name is on the second argument array
// This function is safe for concurrent use.
func FillOnly(strukt any, fields ...string) {
	fillWithDetails(strukt, fields, []string{})
}

// FillExcept fills fields that have a fako: tag and its name is not on the second argument array
// This function is safe for concurrent use.
func FillExcept(strukt any, fields ...string) {
	fillWithDetails(strukt, []string{}, fields)
}

func fillWithDetails(strukt any, only []string, except []string) {
	elem := reflect.ValueOf(strukt).Elem()
	elemT := reflect.TypeOf(strukt).Elem()

	for i := range elem.NumField() {
		field := elem.Field(i)
		fieldt := elemT.Field(i)

		fakerID := fieldt.Tag.Get("fako")

		// If no fako tag is specified, try to match field name with faker functions
		if fakerID == "" {
			// Convert field name to camelized format and check if a faker function exists
			fieldNameCamelized := camelize(fieldt.Name)
			_, exists := allGenerators()[fieldNameCamelized]
			if !exists {
				continue
			}

			fakerID = fieldNameCamelized
		} else {
			fakerID = camelize(fakerID)
		}

		function := findFakeFunctionFor(fakerID)
		inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
		notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

		if field.CanSet() && fakerID != "" && inOnly && notInExcept {
			field.SetString(function())
		}
	}
}
