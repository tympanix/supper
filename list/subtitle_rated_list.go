package list

import (
	"sort"

	"github.com/tympanix/supper/types"
)

type subtitleEntry struct {
	subtitle types.Subtitle
	score    float32
}

// Score returns the score of the rated subtitle
func (s subtitleEntry) Score() float32 {
	return s.score
}

// Subtitle returns the underlying subtitles (i.e. for type assertion)
func (s subtitleEntry) Subtitle() types.Subtitle {
	return s.subtitle
}

// NewRatedSubtitles returns a new subtitles collection
func NewRatedSubtitles(media types.Media, e types.Evaluator, subs ...types.Subtitle) types.RatedSubtitleList {
	var rated []types.RatedSubtitle
	for _, s := range subs {
		score := e.Evaluate(media, s.ForMedia())
		if score > 0.0 {
			rated = append(rated, subtitleEntry{s, score})
		}
	}

	list := RatedSubtitles(rated)
	sort.Sort(sort.Reverse(list))
	return list
}

// RatedSubtitles is a subtitle which is rated by some score
type RatedSubtitles []types.RatedSubtitle

// List returns the list of subtitles as a slice
func (s RatedSubtitles) List() []types.RatedSubtitle {
	subs := make([]types.RatedSubtitle, len(s))
	for i, v := range s {
		subs[i] = v
	}
	return subs
}

// Best returns the best matching subtitle
func (s RatedSubtitles) Best() types.RatedSubtitle {
	if len(s) > 0 {
		return (s)[0]
	}
	return nil
}

// FilterScore return all subtitles with score greater than or equal to some value
func (s RatedSubtitles) FilterScore(score float32) types.RatedSubtitleList {
	_subs := make([]types.RatedSubtitle, 0)
	for _, sub := range s {
		if sub.Score() >= score {
			_subs = append(_subs, sub)
		}
	}
	return RatedSubtitles(_subs)
}

func (s RatedSubtitles) Len() int {
	return len(s)
}

func (s RatedSubtitles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RatedSubtitles) Less(i, j int) bool {
	return s[i].Score() < s[j].Score()
}
