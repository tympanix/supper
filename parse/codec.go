package parse

import "github.com/tympanix/supper/meta/codec"

var codecMap = map[string]interface{}{
	"HEVC":      codec.HEVC,
	"AVC":       codec.AVC,
	"[xh].?265": codec.X265,
	"[xh].?264": codec.X264,
	"XviD":      codec.XviD,
	"DivX":      codec.DivX,
	"WMV":       codec.WMV,
}

// Codecs is a list of parseable codecs
var Codecs = makeMatcher(codecMap)

// CodecIndex returns the codec tag and the index in the string
func CodecIndex(str string) ([]int, codec.Tag) {
	i, c := Codecs.FindTagIndex(str)
	if c != nil {
		return i, c.(codec.Tag)
	}
	return nil, codec.None
}

// Codec returns the codec tag, if found, in a string
func Codec(str string) codec.Tag {
	_, c := CodecIndex(str)
	return c
}
