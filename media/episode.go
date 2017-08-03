package media

import (
	"os"
	"regexp"
	"strconv"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

var episodeRegexp = regexp.MustCompile(`^(.*?[\w)]+)[\W_]+?[Ss]?(\d{1,2})[Eex](\d{1,2})[\W_](.*)$`)

// Episode represents an episode from a TV show
type Episode struct {
	name    string
	episode int
	season  int
	tags    string
	quality string
	codec   string
	source  string
	group   string
}

// NewEpisode parses media info from a file
func NewEpisode(file os.FileInfo) *Episode {
	groups := episodeRegexp.FindStringSubmatch(parse.Filename(file))

	if groups == nil {
		return nil
	}

	name := groups[1]
	season, err := strconv.Atoi(groups[2])

	if err != nil {
		return nil
	}

	episode, err := strconv.Atoi(groups[3])

	if err != nil {
		return nil
	}

	tags := groups[4]

	return &Episode{
		name:    parse.CleanName(name),
		episode: episode,
		season:  season,
		tags:    tags,
		quality: parse.Quality(tags),
		codec:   parse.Codec(tags),
		source:  parse.Source(tags),
		group:   parse.Group(tags),
	}
}

// Matches an episode against a subtitle
func (e *Episode) Matches(types.Subtitle) bool {
	return true
}

// Score returns the likelyhood of the subtitle maching the episode
func (e *Episode) Score(types.Subtitle) float32 {
	return 0.0
}
