package media

import (
	"fmt"
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
func NewMovie(file os.FileInfo) *Movie {
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
		name:    parse.CleanName(name),
		tags:    tags,
		year:    year,
		quality: parse.Quality(tags),
		codec:   parse.Codec(tags),
		source:  parse.Source(tags),
		group:   parse.Group(tags),
	}
}

func (m *Movie) String() string {
	return fmt.Sprintf("%-48.44q%-8d%-12s%-12s%-12s%-24s", m.name, m.year, m.source, m.quality, m.codec, m.group)
}

// Matches a movie against a subtitle
func (m *Movie) Matches(types.Subtitle) bool {
	return true
}

// Score returns the likelyhood of the subtitle maching the movie
func (m *Movie) Score(types.Subtitle) float32 {
	return 0.0
}
