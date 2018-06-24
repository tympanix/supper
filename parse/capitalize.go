package parse

import (
	"regexp"
	"strings"
	"unicode"
)

var nonCapitalized = []string{
	// articles
	"a",
	"an",
	"the",

	// coordinating conjunctions
	"for",
	"and",
	"nor",
	"but",
	"or",
	"yet",
	"so",

	// prepositions (length < 5)
	"ago",
	"anti",
	"as",
	"at",
	"but",
	"by",
	"down",
	"for",
	"from",
	"in",
	"into",
	"like",
	"near",
	"of",
	"off",
	"on",
	"onto",
	//"over",
	"past",
	"per",
	"plus",
	"save",
	"than",
	"to",
	"up",
	//"upon",
	"via",
	"with",
}

var wordRegex = regexp.MustCompile(`[\p{L}0-9]+`)

func shouldCapitalize(str string) bool {
	lower := strings.ToLower(str)
	for _, v := range nonCapitalized {
		if lower == v {
			return false
		}
	}
	return true
}

var romanRegex = regexp.MustCompile(`(?i)^[IVXMC]+$`)

func isRoman(str string) bool {
	return romanRegex.MatchString(str)
}

var breakRegex = regexp.MustCompile(`[\.;:-]\s*[\p{L}0-9]+`)

func isUpper(str string) bool {
	var upper int
	for _, char := range str {
		if unicode.IsLower(char) {
			return false
		} else if unicode.IsUpper(char) {
			upper++
		}
	}
	return upper > 0
}

// Capitalize returns the string with proper english capitalization
func Capitalize(str string) string {
	if isUpper(str) && len(str) > 3 {
		str = strings.ToLower(str)
	}
	str = strings.Title(str)
	var i int
	str = wordRegex.ReplaceAllStringFunc(str, func(word string) string {
		defer func() {
			i++
		}()
		if i == 0 {
			return word
		}
		if !shouldCapitalize(word) {
			return strings.ToLower(word)
		}
		if isRoman(word) {
			return strings.ToUpper(word)
		}
		return word
	})
	str = breakRegex.ReplaceAllStringFunc(str, func(word string) string {
		return strings.Title(word)
	})
	str = strings.Replace(str, "'S", "'s", -1)
	return str
}
