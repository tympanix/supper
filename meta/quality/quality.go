package quality

// Tag is an enum for media quality (720p, 1080p ect.)
type Tag int

func (t Tag) String() string {
	q, ok := stringer[t]
	if ok {
		return q
	}
	panic("unknown quality tag")
}

const (
	None Tag = iota
	UHD2160p
	QHD1440p
	HD1080p
	HD720p
	SD576p
	SD480p
)

var stringer = map[Tag]string{
	None:     "",
	UHD2160p: "UHD",
	QHD1440p: "QHD",
	HD1080p:  "1080p",
	HD720p:   "720p",
	SD576p:   "576p",
	SD480p:   "480p",
}
