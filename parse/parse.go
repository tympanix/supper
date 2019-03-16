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

// FindTagIndex finds an index contained in the regex map and, if found, returns
// the position of the tag in the string and the tag itself as an interface value
func (r regexMatcher) FindTagIndex(str string) ([]int, interface{}) {
	lower := strings.ToLower(str)
	for reg, tag := range r {
		if m := reg.FindStringIndex(lower); m != nil {
			return m, tag
		}
	}
	return nil, nil
}

// FindTag is a helper function for FindTagIndex which only returns the tag value itself
func (r regexMatcher) FindTag(str string) interface{} {
	if _, t := r.FindTagIndex(str); t != nil {
		return t
	}
	return nil
}

// FindAllTagsIndex find all tags in a string and returns two slices, the index positions
// of the tags and the tag values themselves
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

var abbreviationList = []string{
	"mr",
	"mrs",
	"dr",
	"vol",
}

func isAbbreviation(str string) bool {
	lower := strings.ToLower(str)
	for _, v := range abbreviationList {
		if lower == v {
			return true
		}
	}
	return false
}

var illegalcharsRegexp = regexp.MustCompile(`[^\p{L}0-9\s&'_\(\)\-,:]`)
var spaceReplaceRegexp = regexp.MustCompile(`[\.\s_]+`)
var websiteRegexp = regexp.MustCompile(`((https?|ftp|smtp):\/\/)?(www.)[a-z0-9]+\.[a-z]+(\/[a-zA-Z0-9#]+\/?)*`)
var trimRegexp = regexp.MustCompile(`^[^\p{L}0-9]*(.+?[\.\)]?)[^\p{L}0-9]*$`)

// CleanName returns the media name cleaned from punctuation
func CleanName(name string) string {
	name = websiteRegexp.ReplaceAllString(name, "")
	name = spaceReplaceRegexp.ReplaceAllString(name, " ")
	name = illegalcharsRegexp.ReplaceAllString(name, "")

	name = cleanAbbreviations(name)

	name = wordRegex.ReplaceAllStringFunc(name, func(match string) string {
		if isAbbreviation(match) {
			return match + "."
		}
		return match
	})

	name = Capitalize(name)

	match := trimRegexp.FindStringSubmatch(name)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}

var abbrevRegex = regexp.MustCompile(`(?:^|[\.\s])((?:\p{L})(?:[\.\s]\p{L})+)(?:[\.\s]|$)`)

func cleanAbbreviations(s string) string {
	g := abbrevRegex.FindAllStringSubmatchIndex(s, -1)
	if g == nil {
		return s
	}
	var res string
	i := 0
	for _, p := range g {
		abbrev := s[p[2]:p[3]]
		r := strings.Join(spaceReplaceRegexp.Split(abbrev, -1), ".")
		res = s[i:p[2]] + r + "."
		i = p[3]
	}
	res = res + s[i:]
	return res
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
