package parse

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/language"
)

var langs = map[string]language.Tag{
	"arabic":      language.Arabic,
	"brazillian":  language.Portuguese,
	"danish":      language.Danish,
	"dutch":       language.Dutch,
	"english":     language.English,
	"farsi":       language.Persian,
	"persian":     language.Persian,
	"finnish":     language.Finnish,
	"french":      language.French,
	"hebrew":      language.Hebrew,
	"indonesian":  language.Indonesian,
	"italian":     language.Italian,
	"malay":       language.Malay,
	"norwegian":   language.Norwegian,
	"romanian":    language.Romanian,
	"spanish":     language.Spanish,
	"swedish":     language.Swedish,
	"turkish":     language.Turkish,
	"vietnamese":  language.Vietnamese,
	"albanian":    language.Albanian,
	"armenian":    language.Armenian,
	"azerbaijani": language.Azerbaijani,
	/*"belarusian":  language.Belarusian,*/
	"bengali": language.Bengali,
	/*"bosnian":     language.Bosnian, */
	"bulgarian":  language.Bulgarian,
	"catalan":    language.Catalan,
	"chinese":    language.Chinese,
	"croatian":   language.Croatian,
	"czech":      language.Czech,
	"georgian":   language.Georgian,
	"german":     language.German,
	"greek":      language.Greek,
	"hindi":      language.Hindi,
	"hungarian":  language.Hungarian,
	"icelandic":  language.Icelandic,
	"japanese":   language.Japanese,
	"korean":     language.Korean,
	"latvian":    language.Latvian,
	"lithuanian": language.Lithuanian,
	"macedonian": language.Macedonian,
	"malayalam":  language.Malayalam,
	"mongolian":  language.Mongolian,
	"polish":     language.Polish,
	"portuguese": language.Portuguese,
	"russian":    language.Russian,
	"serbian":    language.Serbian,
	"slovak":     language.Slovak,
	"selovenian": language.Slovenian,
	/*"somalia":    language.Somalia,*/
	"swahili":   language.Swahili,
	"tamil":     language.Tamil,
	"telugu":    language.Telugu,
	"thai":      language.Thai,
	"ukrainian": language.Ukrainian,
	"urdu":      language.Urdu,
}

var langRegex = regexp.MustCompile(`[^A-Za-z]+`)

// Language returns a language taken when given the english word for a language
// (e.g. english). If the string is not a known language an error is returned
func Language(lang string) (language.Tag, error) {
	lang = strings.ToLower(lang)
	words := langRegex.Split(lang, -1)

	for _, word := range words {
		if tag, ok := langs[word]; ok {
			return tag, nil
		}
	}

	return language.Und, fmt.Errorf("Could not find language: %s", lang)
}
