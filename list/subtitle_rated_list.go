package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tympanix/supper/score"
	"github.com/tympanix/supper/types"
)

type subtitleEntry struct {
	types.Subtitle
	score float32
}

// Score returns the score of the rated subtitle
func (s subtitleEntry) Score() float32 {
	return s.score
}

// NewRatedSubtitles returns a new subtitles collection
func NewRatedSubtitles(media types.Media, subs ...types.Subtitle) types.RatedSubtitleList {
	list := &RatedSubtitles{
		Evaluator: new(score.DefaultEvaluator),
		media:     media,
	}

	var rated []types.RatedSubtitle
	for _, s := range subs {
		score := list.Evaluate(media, s.ForMedia())
		if score > 0.0 {
			rated = append(rated, subtitleEntry{s, score})
		}
	}

	list.subs = rated
	sort.Sort(sort.Reverse(list))
	return list
}

// RatedSubtitles is a subtitle which is rated by some score
type RatedSubtitles struct {
	types.Evaluator
	media types.Media
	subs  []types.RatedSubtitle
}

// MarshalJSON returns a JSON representation of the rated subtitle list
func (s *RatedSubtitles) MarshalJSON() (b []byte, err error) {
	return json.Marshal(s.List())
}

// List returns the list of subtitles as a slice
func (s *RatedSubtitles) List() []types.RatedSubtitle {
	subs := make([]types.RatedSubtitle, len(s.subs))
	for i, v := range s.subs {
		subs[i] = v
	}
	return subs
}

// Best returns the best matching subtitle
func (s *RatedSubtitles) Best() types.RatedSubtitle {
	if len(s.subs) > 0 {
		return s.subs[0]
	}
	return nil
}

// FilterScore return all subtitles with score greater than or equal to some value
func (s *RatedSubtitles) FilterScore(score float32) types.RatedSubtitleList {
	_subs := make([]types.RatedSubtitle, 0)
	for _, sub := range s.subs {
		if sub.Score() >= score {
			_subs = append(_subs, sub)
		}
	}
	return &RatedSubtitles{
		Evaluator: s.Evaluator,
		media:     s.media,
		subs:      _subs,
	}
}

func (s *RatedSubtitles) Len() int {
	return len(s.subs)
}

func (s *RatedSubtitles) Swap(i, j int) {
	s.subs[i], s.subs[j] = s.subs[j], s.subs[i]
}

func (s *RatedSubtitles) Less(i, j int) bool {
	return s.subs[i].Score() < s.subs[j].Score()
}

func (s *RatedSubtitles) String() string {
	var buffer bytes.Buffer

	for _, sub := range s.subs {
		buffer.WriteString(fmt.Sprintf("%-8.2f%v\n", sub.Score(), sub))
	}

	return buffer.String()
}
