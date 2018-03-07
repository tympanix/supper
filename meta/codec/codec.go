package codec

type Tag int

const (
	HEVC Tag = iota
	AVC
	X265
	X264
	XviD
	DivX
	WMV
)
