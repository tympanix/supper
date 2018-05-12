package parse

import "github.com/tympanix/supper/meta/misc"

var miscMap = map[string]interface{}{
	"3D":  misc.Video3D,
	"HC":  misc.HC,
	"DTS": misc.DTS,
	"AC3": misc.AC3,
}

var miscMatcher = makeMatcher(miscMap)

// Miscellaneous returns a list of miscellaneous media tags from a string
func Miscellaneous(name string) misc.List {
	var list misc.List
	for _, l := range miscMatcher.FindAll(name) {
		list = append(list, l.(misc.Tag))
	}
	return list
}
