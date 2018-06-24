package parse

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type regexMatcher map[*regexp.Regexp]interface{}

func makeMatcher(tags map[string]interface{}) regexMatcher {
	regs := make(map[*regexp.Regexp]interface{})
	for reg, tag := range tags {
		regexpstr := fmt.Sprintf("(?i)\\b(%s)\\b", reg)
		regs[regexp.MustCompile(regexpstr)] = tag
	}
	return regs
}

func (r regexMatcher) FindTagIndex(str string) ([]int, interface{}) {
	lower := strings.ToLower(str)
	for reg, tag := range r {
		if m := reg.FindStringIndex(lower); m != nil {
			return m, tag
		}
	}
	return nil, nil
}

func (r regexMatcher) FindTag(str string) interface{} {
	if _, t := r.FindTagIndex(str); t != nil {
		return t
	}
	return nil
}

func (r regexMatcher) FindAllTagsIndex(str string) ([]int, []interface{}) {
	var tags []interface{}
	var idx []int
	lower := strings.ToLower(str)
	for reg, tag := range r {
		if loc := reg.FindStringIndex(lower); loc != nil {
			idx = append(idx, loc...)
			tags = append(tags, tag)
		}
	}
	return idx, tags
}

// FindAll returns a list of all matched tags
func (r regexMatcher) FindAllTags(str string) []interface{} {
	_, tags := r.FindAllTagsIndex(str)
	return tags
}

// Filename returns the filename of the file without extension
func Filename(filename string) string {
	f := filepath.Base(filename)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

var abbreviationRegexp = regexp.MustCompile(`\s[A-Z]\s[A-Z](\s[A-Z])*\s`)
var illegalcharsRegexp = regexp.MustCompile(`[^\p{L}0-9\s&'_\(\)-]`)
var spaceReplaceRegexp = regexp.MustCompile(`[\.\s_]+`)

// CleanName returns the media name cleaned from punctuation
func CleanName(name string) string {
	name = spaceReplaceRegexp.ReplaceAllString(name, " ")
	name = illegalcharsRegexp.ReplaceAllString(name, "")

	name = abbreviationRegexp.ReplaceAllStringFunc(name, func(match string) string {
		return " " + strings.Replace(match, " ", "", -1) + " "
	})

	name = Capitalize(name)

	return strings.TrimSpace(name)
}

var illegalIdentity = regexp.MustCompile(`[^\p{L}0-9]`)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// Identity returns a string where special characters are removed. The returned
// string is suitable for use in an identity string
func Identity(str string) string {
	var err error
	var ident string
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	ident, _, err = transform.String(t, str)
	if err != nil {
		ident = str
	}
	ident = illegalIdentity.ReplaceAllString(ident, "")
	ident = strings.ToLower(ident)
	return ident
}

var tagsRegexp = regexp.MustCompile(`[^\p{L}0-9]+`)

// Tags returns a string as tags split by non-word characters
func Tags(name string) []string {
	return tagsRegexp.Split(name, -1)
}
