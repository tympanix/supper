package source

// Tag is an enum type for media sources
type Tag int

func (t Tag) String() string {
	s, ok := stringer[t]
	if ok {
		return s
	}
	panic("unknown source tags")
}

const (
	None Tag = iota
	Remux
	BluRay
	WEBDL
	WEBRip
	VODRip
	HDTV
	DVDR
	DVDRip
	R5
	Screener
	Telecine
	Workprint
	Telesync
	Cam
)

var stringer = map[Tag]string{
	None:      "",
	Remux:     "Remux",
	BluRay:    "BluRay",
	WEBDL:     "WEB-DL",
	WEBRip:    "WEB-Rip",
	VODRip:    "VOD-Rip",
	HDTV:      "HDTV",
	DVDR:      "DVDR",
	DVDRip:    "DVD-Rip",
	R5:        "R5",
	Screener:  "Screener",
	Telecine:  "Telecine",
	Workprint: "Workprint",
	Telesync:  "Telesync",
	Cam:       "Cam",
}
