package source

type Source int

const (
	Remux Source = iota
	BluRay
	WEBDL
	WEBRip
	VODRip
	HDTV
	DVDRip
	R5
	Screener
	Telecine
	Workprint
	Telesync
	Cam
)
