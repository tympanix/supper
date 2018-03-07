package parse

import "github.com/tympanix/supper/meta/source"

var sourceMap = map[string]interface{}{
	"CAMRip|CAM":                                             source.Cam,
	"TS|HDTS|TELESYNC|PDVD":                                  source.Telesync,
	"WP|WORKPRINT":                                           source.Workprint,
	"TC|HDTC|TELECINE":                                       source.Telecine,
	"SCR|SCREENER|DVDSCR|DVDSCREENER|DBSCR":                  source.Screener,
	"R5|R5.LINE":                                             source.R5,
	"DVD-?Rip|DVDMux":                                        source.DVDRip,
	"DVD-?(R|5|9)|DVD-?Full|ISO":                             source.DVDR,
	"DSR|DSRip|SATRip|DTHRip|DVBRip|HDTV|PDTV|TVRip|HDTVRip": source.HDTV,
	"VOD-?Rip|VODR":                                          source.VODRip,
	"WEB.?DL(Rip)?|HDRip":                                    source.WEBDL,
	"WEB.?Rip":                                               source.WEBRip,
	"Blu.?Ray|BDRip|BRRip|BDMV|BDR|BD(5|9)":                  source.BluRay,
}

// Sources lists all prossible sources to parse
var Sources = makeMatcher(sourceMap)
