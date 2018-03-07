package codec

type Codec int

const (
	HEVC Codec = iota
	AVC
	x265
	x264
	XviD
	DivX
	WMV
)
