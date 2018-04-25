package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

var episodeRegexp = regexp.MustCompile(`^(.*?[\w)]+)[\W_]+?[Ss]?(\d{1,2})[Eex](\d{1,2})(?:[Ee]\d{1,2})?[\W_]*(.*)$`)

// EpisodeMeta represents an episode from a TV show
type EpisodeMeta struct {
	Metadata
	NameX        string
	EpisodeNameX string
	EpisodeX     int
	SeasonX      int
}

func (e *EpisodeMeta) MarshalJSON() (b []byte, err error) {
	type jsonEpisode struct {
		Meta    Metadata `json:"metadata"`
		Name    string   `json:"name"`
		Episode int      `json:"episode"`
		Seasion int      `json:"season"`
	}

	return json.Marshal(jsonEpisode{
		e.Metadata,
		e.TVShow(),
		e.Episode(),
		e.Season(),
	})
}

// NewEpisode parses media info from a file
func NewEpisode(filename string) (*EpisodeMeta, error) {
	groups := episodeRegexp.FindStringSubmatch(filename)

	if groups == nil {
		return nil, errors.New("could not parse media")
	}

	name := groups[1]
	season, err := strconv.Atoi(groups[2])

	if err != nil {
		return nil, err
	}

	episode, err := strconv.Atoi(groups[3])

	if err != nil {
		return nil, err
	}

	tags := groups[4]

	return &EpisodeMeta{
		Metadata: ParseMetadata(tags),
		NameX:    parse.CleanName(name),
		EpisodeX: episode,
		SeasonX:  season,
	}, nil
}

func (e *EpisodeMeta) String() string {
	return fmt.Sprintf("%s S%02dE%02d", e.TVShow(), e.Season(), e.Episode())
}

// TVShow is the name of the TV show
func (e *EpisodeMeta) TVShow() string {
	return e.NameX
}

// EpisodeName is the name of the episode
func (e *EpisodeMeta) EpisodeName() string {
	return e.EpisodeNameX
}

// Episode is the episode number in the season
func (e *EpisodeMeta) Episode() int {
	return e.EpisodeX
}

// Season is the season number of the show
func (e *EpisodeMeta) Season() int {
	return e.SeasonX
}

// Matches an episode against a subtitle
func (e *EpisodeMeta) Matches(types.Subtitle) bool {
	return true
}

// Score returns the likelyhood of the subtitle maching the episode
func (e *EpisodeMeta) Score(types.Subtitle) float32 {
	return 0.0
}
