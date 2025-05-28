package fako

import (
	"strings"
	"unicode"
)

// camelize converts a string to CamelCase
// Examples:
//
//	"hello_world" -> "HelloWorld"
//	"hello-world" -> "HelloWorld"
//	"hello world" -> "HelloWorld"
//	"helloWorld" -> "HelloWorld"
func camelize(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	var currentWord strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			// Delimiter found, process current word
			if currentWord.Len() > 0 {
				addCamelWord(&result, currentWord.String())
				currentWord.Reset()
			}
		} else if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) {
			// CamelCase boundary (lowercase followed by uppercase)
			if currentWord.Len() > 0 {
				addCamelWord(&result, currentWord.String())
				currentWord.Reset()
			}
			currentWord.WriteRune(r)
		} else {
			currentWord.WriteRune(r)
		}
	}

	// Process the last word
	if currentWord.Len() > 0 {
		addCamelWord(&result, currentWord.String())
	}

	return result.String()
}

// addCamelWord adds a word to the result with proper CamelCase formatting
func addCamelWord(result *strings.Builder, word string) {
	if word == "" {
		return
	}

	runes := []rune(strings.ToLower(word))
	result.WriteRune(unicode.ToUpper(runes[0]))
	for _, r := range runes[1:] {
		result.WriteRune(r)
	}
}
