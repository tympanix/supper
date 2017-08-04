package parse

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// Collection is a collection of tags
type Collection map[string]*regexp.Regexp

// FindTag find a tag from the collection in the string
func (t Collection) FindTag(str string) string {
	lower := strings.ToLower(str)
	for tag, reg := range t {
		if reg.MatchString(lower) {
			return tag
		}
	}
	return ""
}

// NewCollection creates a new collection of tags
func NewCollection(tags []string) Collection {
	regs := make(map[string]*regexp.Regexp)
	for _, tag := range tags {
		regexpstr := fmt.Sprintf("\\b%s\\b", strings.ToLower(tag))
		regs[tag] = regexp.MustCompile(regexpstr)
	}
	return regs
}

// Filename returns the filename of the file without extension
func Filename(filename string) string {
	f := path.Base(filename)
	return strings.TrimSuffix(f, path.Ext(f))
}

var abbreviationRegexp = regexp.MustCompile(`[A-Z](\.)`)
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

// Source parses the source from a filename
func Source(name string) string {
	return Sources.FindTag(name)
}

// Quality finds the quality of the media
func Quality(name string) string {
	return Qualities.FindTag(name)
}

// Codec parses the codec from a file name
func Codec(name string) string {
	return Codecs.FindTag(name)
}

var tagsRegexp = regexp.MustCompile(`[\W_]+`)

// Tags returns a string as tags split by non-word characters
func Tags(name string) []string {
	return tagsRegexp.Split(name, -1)
}
