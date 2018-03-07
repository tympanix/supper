package quality

// Quality is an enum for media quality (720p, 1080p ect.)
type Tag int

const (
	UHD2160p Tag = iota
	QHD1440p
	HD1080p
	HD720p
	SD576p
	SD480p
)
