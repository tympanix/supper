package media

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
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
func NewMovie(filename string) (*MovieMeta, error) {
	groups := movieRegexp.FindStringSubmatch(filename)

	if groups == nil {
		return nil, errors.New("Could not parse media")
	}

	name := groups[1]
	year, err := strconv.Atoi(groups[2])
	tags := groups[3]

	if err != nil {
		return nil, err
	}

	return &MovieMeta{
		Metadata: ParseMetadata(tags),
		name:     parse.CleanName(name),
		tags:     tags,
		year:     year,
	}, nil
}

func (m *MovieMeta) String() string {
	return fmt.Sprintf("%s (%v)", m.MovieName(), m.Year())
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
