package parse

import "github.com/tympanix/supper/media/meta/source"

var sourceMap = map[string]interface{}{
	"CAM.?Rip|CAM|HD.?CAM":                      source.Cam,
	"TS|HD.?TS|TELESYNC|PDVD":                   source.Telesync,
	"WP|WORKPRINT":                              source.Workprint,
	"TC|HD.?TC|TELECINE":                        source.Telecine,
	"SCR|SCREENER|DVD.?SCR|DVD.?SCREENER|DBSCR": source.Screener,
	"R5|R5.LINE":                                source.R5,
	"DVD.?Rip|DVD.?Mux":                         source.DVDRip,
	"DVD.?(R|5|9)(.?Full)?|ISO":                 source.DVDR,
	"(DSR|DS|SAT|DTH|DVB|TV|HDTV|PDTV)(-?Rip)?": source.HDTV,
	"VOD-?Rip|VODR":                             source.VODRip,
	"WEB.?DL(Rip)?|HD.?Rip|WEB":                 source.WEBDL,
	"WEB.?Rip":                                  source.WEBRip,
	"Blu.?Ray|(BD|BR).?Rip|BD.?(R|5|9)":         source.BluRay,
}

// Sources lists all prossible sources to parse
var Sources = makeMatcher(sourceMap)

// SourceIndex returns the source tag and the index in the string
func SourceIndex(str string) ([]int, source.Tag) {
	i, s := Sources.FindTagIndex(str)
	if s != nil {
		return i, s.(source.Tag)
	}
	return nil, source.None
}

// Source returns the source tag, if found, in a string
func Source(str string) source.Tag {
	_, s := SourceIndex(str)
	return s
}
