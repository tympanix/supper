package media

import (
	"os"
	"regexp"
	"strconv"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

// Movie represents a movie file
type Movie struct {
	name    string
	tags    string
	year    int
	quality string
	group   string
	codec   string
	source  string
}

var movieRegexp = regexp.MustCompile(`^(.+?)[\W_]?\((\d{4})\)[\W_]?(.*)$`)

// NewMovie parses media info from a file
func NewMovie(file *os.File) *Movie {
	groups := movieRegexp.FindStringSubmatch(parse.Filename(file))

	if groups == nil {
		return nil
	}

	name := groups[1]
	year, err := strconv.Atoi(groups[2])
	tags := groups[3]

	if err != nil {
		return nil
	}

	return &Movie{
		name:    name,
		tags:    tags,
		year:    year,
		quality: parse.Quality(tags),
		codec:   parse.Codec(tags),
		source:  parse.Source(tags),
	}
}

// Matches a movie against a subtitle
func (m *Movie) Matches(types.Subtitle) bool {
	return true
}

// Score returns the likelyhood of the subtitle maching the movie
func (m *Movie) Score(types.Subtitle) float32 {
	return 0.0
}
