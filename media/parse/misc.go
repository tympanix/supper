package parse

import "github.com/tympanix/supper/media/meta/misc"

var miscMap = map[string]interface{}{
	"3D":                              misc.Video3D,
	"HC|Hardcoded?":                   misc.HC,
	"DTS(.?HD)?":                      misc.DTS,
	"DD(\\+|5\\.1|P)?|TrueHD|Atmos":   misc.DolbyDigital,
	"Extended(.(Cut|Edition))?":       misc.Extended,
	"(DD[\\+P]?|TrueHD|MA|DTS)?5\\.1": misc.Surround5x1,
	"(DD[\\+P]?|TrueHD|MA|DTS)?7\\.1": misc.Surround7x1,
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
