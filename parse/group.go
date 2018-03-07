package parse

import (
	"regexp"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
)

var groupRegexp = regexp.MustCompile(`[a-zA-Z0-9]+$`)

// Group parses and returns the release group
func Group(name string) string {
	group := groupRegexp.FindString(name)

	if Quality(group) != quality.None {
		return ""
	} else if Source(group) != source.None {
		return ""
	} else if Codec(group) != codec.None {
		return ""
	} else if len(group) <= 1 {
		return ""
	}

	return group
}
