package fako

import (
	"reflect"
	"slices"
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

// fillWithDetails fills fields that have a fako: tag and its name is on the second argument array or not on the third argument array
// This function is safe for concurrent use.
func fillWithDetails(strukt any, only []string, except []string) {
	elem := reflect.ValueOf(strukt).Elem()
	elemT := reflect.TypeOf(strukt).Elem()

	for i := range elem.NumField() {
		field := elem.Field(i)
		fieldt := elemT.Field(i)

		fakerID := fieldt.Tag.Get("fako")
		fakerID = camelize(fakerID)

		// If no fako tag is specified, try to match field name with faker functions
		if fakerID == "" {
			// Convert field name to camelized format and check if a faker function exists
			fakerID = camelize(fieldt.Name)
			_, exists := allGenerators()[fakerID]
			if !exists {
				continue
			}
		}

		if fakerID == "" {
			continue
		}

		inOnly := len(only) == 0 || (len(only) > 0 && slices.Contains(only, fieldt.Name))
		inExcept := len(except) != 0 && slices.Contains(except, fieldt.Name)
		if !inOnly || inExcept {
			continue
		}

		if !field.CanSet() {
			continue
		}

		function := findFakeFunctionFor(fakerID)
		field.SetString(function())
	}
}
