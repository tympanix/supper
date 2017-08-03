package parse

import (
	"fmt"
	"os"
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
func Filename(file *os.File) string {
	f := path.Base(file.Name())
	return strings.TrimSuffix(f, path.Ext(f))
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
