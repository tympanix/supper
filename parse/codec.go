package parse

import "github.com/tympanix/supper/meta/codec"

var codecMap = map[string]interface{}{
	"HEVC":    codec.HEVC,
	"AVC":     codec.AVC,
	"[xh]265": codec.X265,
	"[xh]264": codec.X264,
	"XviD":    codec.XviD,
	"DivX":    codec.DivX,
	"WMV":     codec.WMV,
}

// Codecs is a list of parseable codecs
var Codecs = makeMatcher(codecMap)
