package parse

import "github.com/tympanix/supper/meta/source"

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
