package parse

import "regexp"

var groupRegexp = regexp.MustCompile(`[a-zA-Z0-9]+$`)

// Group parses and returns the release group
func Group(name string) string {
	group := groupRegexp.FindString(name)

	if len(Quality(group)) > 0 {
		return ""
	} else if len(Source(group)) > 0 {
		return ""
	} else if len(Codec(group)) > 0 {
		return ""
	} else if len(group) <= 1 {
		return ""
	}

	return group
}
