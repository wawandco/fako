package fako

import (
	"reflect"
	"strconv"

	"math/rand"
	"strings"
	"time"

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

		switch field.Kind() {
		case reflect.String:
			if fakeType != "" {
				fakeType = snaker.SnakeToCamel(fakeType)
				function := findFakeFunctionFor(fakeType)

				inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
				notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

				if field.CanSet() && fakeType != "" && inOnly && notInExcept {
					field.SetString(function())
				}

				continue
			}

		case reflect.Int, reflect.Int32, reflect.Int64:
			min, max := extractMinMaxInt64(fakeType)
			rand.Seed(time.Now().Unix())
			val := rand.Int63n(max) + min
			field.SetInt(val)
		}

	}
}

/** extractMinMax function extracts min and max from the fako tag **/
func extractMinMaxInt64(base string) (int64, int64) {
	min, max := 0, 0
	parts := strings.Split(base, ";")
	for _, part := range parts {
		content := strings.Split(part, "=")
		if len(content) == 2 {
			switch content[0] {
			case "min":
				min, _ = strconv.Atoi(content[1])
			case "max":
				max, _ = strconv.Atoi(content[1])
			}
		}
	}

	if min > max {
		tmp := max
		max = min
		min = tmp
	}

	return int64(min), int64(max)
}
