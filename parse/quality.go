package parse

import "github.com/tympanix/supper/meta/quality"

var qualityMap = map[string]interface{}{
	"2160p|4K|UHD|Ultra.?HD": quality.UHD2160p,
	"QHD":   quality.QHD1440p,
	"1080p": quality.HD1080p,
	"720p":  quality.HD720p,
	"576p":  quality.SD576p,
	"480p":  quality.SD480p,
}

// Qualities list all possible qualities to parse
var Qualities = makeMatcher(qualityMap)
