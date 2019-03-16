package parse

import (
	"regexp"

	"github.com/tympanix/supper/media/meta/codec"
	"github.com/tympanix/supper/media/meta/quality"
	"github.com/tympanix/supper/media/meta/source"
)

var groupRegexp = regexp.MustCompile(`[a-zA-Z0-9]+$`)

// GroupIndex returns the release group as well as the index for the match
func GroupIndex(str string) ([]int, string) {
	idx := groupRegexp.FindStringIndex(str)

	if idx == nil {
		return nil, ""
	}

	group := str[idx[0]:idx[1]]

	if Quality(group) != quality.None {
		return nil, ""
	} else if Source(group) != source.None {
		return nil, ""
	} else if Codec(group) != codec.None {
		return nil, ""
	} else if len(group) <= 1 {
		return nil, ""
	}

	return idx, group
}

// GroupAfter returns the release group of the string only occuring after
// the position of idx
func GroupAfter(idx int, str string) string {
	m, g := GroupIndex(str)
	if m == nil {
		return ""
	}
	if m[0] <= idx {
		return ""
	}
	return g
}

// Group parses and returns the release group
func Group(str string) string {
	_, g := GroupIndex(str)
	return g
}
