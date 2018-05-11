package parse

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
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

func (r regexMatcher) FindTag(str string) interface{} {
	lower := strings.ToLower(str)
	for reg, tag := range r {
		if reg.MatchString(lower) {
			return tag
		}
	}
	return nil
}

// Filename returns the filename of the file without extension
func Filename(filename string) string {
	f := filepath.Base(filename)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

var abbreviationRegexp = regexp.MustCompile(`[A-Z]\s[A-Z](\s[A-Z])*`)
var illegalcharsRegexp = regexp.MustCompile(`[^\p{L}0-9\s&'_\(\)-]`)

// CleanName returns the media name cleaned from punctuation
func CleanName(name string) string {
	name = strings.Replace(name, ". ", " ", -1)
	name = strings.Replace(name, ".", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = illegalcharsRegexp.ReplaceAllString(name, "")

	name = abbreviationRegexp.ReplaceAllStringFunc(name, func(match string) string {
		return strings.Replace(match, " ", "", -1)
	})

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
	return ident
}

// Source parses the source from a filename
func Source(name string) source.Tag {
	s := Sources.FindTag(name)
	if s != nil {
		return s.(source.Tag)
	}
	return source.None
}

// Quality finds the quality of the media
func Quality(name string) quality.Tag {
	q := Qualities.FindTag(name)
	if q != nil {
		return q.(quality.Tag)
	}
	return quality.None
}

// Codec parses the codec from a file name
func Codec(name string) codec.Tag {
	c := Codecs.FindTag(name)
	if c != nil {
		return c.(codec.Tag)
	}
	return codec.None
}

var tagsRegexp = regexp.MustCompile(`[\W_]+`)

// Tags returns a string as tags split by non-word characters
func Tags(name string) []string {
	return tagsRegexp.Split(name, -1)
}
