package parse

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
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

var abbreviationRegexp = regexp.MustCompile(`[A-Z](\s)[A-Z]`)
var illegalcharsRegexp = regexp.MustCompile(`[^\w\s&'_\(\)-]`)

// CleanName returns the movie name cleaned from punctuation
func CleanName(name string) string {
	name = abbreviationRegexp.ReplaceAllStringFunc(name, func(match string) string {
		return strings.Replace(match, ".", "", -1)
	})

	name = strings.Replace(name, ". ", " ", -1)
	name = strings.Replace(name, ".", " ", -1)
	name = illegalcharsRegexp.ReplaceAllString(name, "")

	return strings.TrimSpace(name)
}

var allowedIdentity = regexp.MustCompile(`[^A-Za-z0-9]`)

// Identity returns a string where special characters are removed. The returned
// string is suitable for use in an identity string
func Identity(str string) string {
	return allowedIdentity.ReplaceAllString(str, "")
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
