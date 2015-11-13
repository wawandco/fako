package fako

import (
	"reflect"

	"github.com/serenize/snaker"
	"github.com/wawandco/fako/Godeps/_workspace/src/github.com/icrowley/fake"
)

var typeMapping = map[string]func() string{
	"Brand":                    fake.Brand,
	"Character":                fake.Character,
	"Characters":               fake.Characters,
	"City":                     fake.City,
	"Color":                    fake.Color,
	"Company":                  fake.Company,
	"Continent":                fake.Continent,
	"Country":                  fake.Country,
	"CreditCardType":           fake.CreditCardType,
	"Currency":                 fake.Currency,
	"CurrencyCode":             fake.CurrencyCode,
	"Digits":                   fake.Digits,
	"DomainName":               fake.DomainName,
	"DomainZone":               fake.DomainZone,
	"EmailAddress":             fake.EmailAddress,
	"EmailBody":                fake.EmailBody,
	"EmailSubject":             fake.EmailSubject,
	"FemaleFirstName":          fake.FemaleFirstName,
	"FemaleFullName":           fake.FemaleFullName,
	"FemaleFullNameWithPrefix": fake.FemaleFullNameWithPrefix,
	"FemaleFullNameWithSuffix": fake.FemaleFullNameWithSuffix,
	"FemaleLastName":           fake.FemaleLastName,
	"FemalePatronymic":         fake.FemalePatronymic,
	"FirstName":                fake.FirstName,
	"FullName":                 fake.FullName,
	"FullNameWithPrefix":       fake.FullNameWithPrefix,
	"FullNameWithSuffix":       fake.FullNameWithSuffix,
	"Gender":                   fake.Gender,
	"GenderAbbrev":             fake.GenderAbbrev,
	"HexColor":                 fake.HexColor,
	"HexColorShort":            fake.HexColorShort,
	"IPv4":                     fake.IPv4,
	"Industry":                 fake.Industry,
	"JobTitle":                 fake.JobTitle,
	"Language":                 fake.Language,
	"LastName":                 fake.LastName,
	"LatitudeDirection":        fake.LatitudeDirection,
	"LongitudeDirection":       fake.LongitudeDirection,
	"MaleFirstName":            fake.MaleFirstName,
	"MaleFullName":             fake.MaleFullName,
	"MaleFullNameWithPrefix":   fake.MaleFullNameWithPrefix,
	"MaleFullNameWithSuffix":   fake.MaleFullNameWithSuffix,
	"MaleLastName":             fake.MaleLastName,
	"MalePatronymic":           fake.MalePatronymic,
	"Model":                    fake.Model,
	"Month":                    fake.Month,
	"MonthShort":               fake.MonthShort,
	"Paragraph":                fake.Paragraph,
	"Paragraphs":               fake.Paragraphs,
	"Patronymic":               fake.Patronymic,
	"Phone":                    fake.Phone,
	"Product":                  fake.Product,
	"ProductName":              fake.ProductName,
	"Sentence":                 fake.Sentence,
	"Sentences":                fake.Sentences,
	"SimplePassword":           fake.SimplePassword,
	"State":                    fake.State,
	"StateAbbrev":              fake.StateAbbrev,
	"Street":                   fake.Street,
	"StreetAddress":            fake.StreetAddress,
	"Title":                    fake.Title,
	"TopLevelDomain":           fake.TopLevelDomain,
	"UserName":                 fake.UserName,
	"WeekDay":                  fake.WeekDay,
	"WeekDayShort":             fake.WeekDayShort,
	"Word":                     fake.Word,
	"Words":                    fake.Words,
	"Zip":                      fake.Zip,
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

		if fakeType != "" {
			fakeType = snaker.SnakeToCamel(fakeType)
			function := findFakeFunctionFor(fakeType)

			inOnly := len(only) == 0 || (len(only) > 0 && contains(only, fieldt.Name))
			notInExcept := len(except) == 0 || (len(except) > 0 && !contains(except, fieldt.Name))

			if field.CanSet() && fakeType != "" && inOnly && notInExcept {
				field.SetString(function())
			}
		}

	}
}

func findFakeFunctionFor(fako string) func() string {
	result := func() string { return "" }

	for kind, function := range typeMapping {
		if fako == kind {
			result = function
			break
		}
	}

	return result
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
