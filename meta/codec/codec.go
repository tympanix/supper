package codec

// Tag is an enum representing media codecs
type Tag int

func (t Tag) String() string {
	s, ok := stringer[t]
	if ok {
		return s
	}
	panic("unknown codec tag")
}

const (
	None Tag = iota
	HEVC
	AVC
	X265
	X264
	XviD
	DivX
	WMV
)

var stringer = map[Tag]string{
	None: "Unknown",
	HEVC: "HEVC",
	AVC:  "AVC",
	X265: "x265",
	X264: "x264",
	XviD: "XviD",
	DivX: "DivX",
	WMV:  "WMV",
}
