package media

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

// MovieMeta represents a movie file
type MovieMeta struct {
	Metadata
	name string
	tags string
	year int
}

// MovieFile is a local movie file on disk
type MovieFile struct {
	os.FileInfo
	types.Movie
}

var movieRegexp = regexp.MustCompile(`^(.+?)[\W_]+(19\d\d|20\d\d)[\W_]+(.*)$`)

// NewMovie parses media info from a file
func NewMovie(filename string) *MovieMeta {
	groups := movieRegexp.FindStringSubmatch(filename)

	if groups == nil {
		return nil
	}

	name := groups[1]
	year, err := strconv.Atoi(groups[2])
	tags := groups[3]

	if err != nil {
		return nil
	}

	return &MovieMeta{
		Metadata: ParseMetadata(tags),
		name:     parse.CleanName(name),
		tags:     tags,
		year:     year,
	}
}

func (m *MovieMeta) String() string {
	return fmt.Sprintf("%-48.44q%-8d%-12s%-12s%-12s%-24s", m.name, m.year, m.source, m.quality, m.codec, m.group)
}

// MovieName is the name of the movie
func (m *MovieMeta) MovieName() string {
	return m.name
}

// Year is the release year of the movie
func (m *MovieMeta) Year() int {
	return m.year
}

// Matches a movie against a subtitle
func (m *MovieMeta) Matches(types.Subtitle) bool {
	return true
}

// Score returns the likelyhood of the subtitle maching the movie
func (m *MovieMeta) Score(types.Subtitle) float32 {
	return 0.0
}
