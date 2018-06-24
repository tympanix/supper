package parse

import "github.com/tympanix/supper/meta/misc"

var miscMap = map[string]interface{}{
	"3D":                        misc.Video3D,
	"HC|Hardcoded?":             misc.HC,
	"DTS(.?HD)?":                misc.DTS,
	"DD(\\+|5\\.1|P)?|TrueHD":   misc.DolbyDigital,
	"Extended(.(Cut|Edition))?": misc.Extended,
	"AC3": misc.AC3,
}

var miscMatcher = makeMatcher(miscMap)

// MiscellaneousIndex returns all miscellaneous tags with all their indexes
func MiscellaneousIndex(str string) ([]int, misc.List) {
	var list misc.List
	var idx []int
	idx, lx := miscMatcher.FindAllTagsIndex(str)
	for _, l := range lx {
		list = append(list, l.(misc.Tag))
	}
	return idx, list
}

// Miscellaneous returns a list of miscellaneous media tags from a string
func Miscellaneous(str string) misc.List {
	_, misc := MiscellaneousIndex(str)
	return misc
}
