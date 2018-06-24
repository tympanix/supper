package parse

import "github.com/tympanix/supper/meta/quality"

var qualityMap = map[string]interface{}{
	"2160p|4K|UHD|Ultra.?HD": quality.UHD2160p,
	"1440p|QHD":              quality.QHD1440p,
	"1080p|Full.?HD":         quality.HD1080p,
	"720p":                   quality.HD720p,
	"576p":                   quality.SD576p,
	"480p":                   quality.SD480p,
}

// Qualities list all possible qualities to parse
var Qualities = makeMatcher(qualityMap)

// QualityIndex returns the quality tag and the index in the string
func QualityIndex(str string) ([]int, quality.Tag) {
	i, q := Qualities.FindTagIndex(str)
	if q != nil {
		return i, q.(quality.Tag)
	}
	return nil, quality.None
}

// Quality returns the quality tag, if found, in the string
func Quality(str string) quality.Tag {
	_, q := QualityIndex(str)
	return q
}
