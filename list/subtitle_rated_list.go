package list

import (
	"bytes"
	"fmt"
	"sort"

	"golang.org/x/text/language"

	"github.com/fatih/set"
	"github.com/tympanix/supper/score"
	"github.com/tympanix/supper/types"
)

type subtitleEntry struct {
	types.Subtitle
	score float32
}

// NewSubtitles returns a new subtitles collection
func RatedSubtitles(media types.LocalMedia) *ratedSubtitles {
	return &ratedSubtitles{
		Evaluator: new(score.DefaultEvaluator),
		media:     media,
		subs:      make([]subtitleEntry, 0),
	}
}

type ratedSubtitles struct {
	types.Evaluator
	media types.LocalMedia
	subs  []subtitleEntry
}

func (s *ratedSubtitles) clone(list []subtitleEntry) types.SubtitleList {
	return &ratedSubtitles{
		Evaluator: s.Evaluator,
		media:     s.media,
		subs:      list,
	}
}

// List returns the list of subtitles as a slice
func (s *ratedSubtitles) List() []types.Subtitle {
	subs := make([]types.Subtitle, len(s.subs))
	for i, v := range s.subs {
		subs[i] = v
	}
	return subs
}

// Best returns the best matching subtitle
func (s *ratedSubtitles) Best() types.Subtitle {
	if len(s.subs) > 0 {
		return s.subs[0]
	}
	return nil
}

// Add a subtitle to the collection
func (s *ratedSubtitles) Add(sub types.Subtitle) {
	if sub == nil || sub.Meta() == nil {
		return
	}
	s.subs = append(s.subs, subtitleEntry{
		Subtitle: sub,
		score:    s.Evaluate(s.media, sub),
	})
	sort.Sort(sort.Reverse(s))
}

// FilterLanguage returns a new subtitle collection including only the argument language
func (s *ratedSubtitles) FilterLanguage(lang language.Tag) types.SubtitleList {
	_subs := make([]subtitleEntry, 0)
	for _, sub := range s.subs {
		if sub.IsLang(lang) {
			_subs = append(_subs, sub)
		}
	}
	return s.clone(_subs)
}

// HearingImpaired returnes a new subtitle collection where hearing impared subtitles has been filtered
func (s *ratedSubtitles) HearingImpaired(hi bool) types.SubtitleList {
	_subs := make([]subtitleEntry, 0)
	for _, sub := range s.subs {
		if sub.IsHI() == hi {
			_subs = append(_subs, sub)
		}
	}
	return s.clone(_subs)
}

func (s *ratedSubtitles) LanguageSet() set.Interface {
	langs := set.New()
	for _, sub := range s.subs {
		langs.Add(sub.Language())
	}
	return langs
}

func (s *ratedSubtitles) Len() int {
	return len(s.subs)
}

func (s *ratedSubtitles) Swap(i, j int) {
	s.subs[i], s.subs[j] = s.subs[j], s.subs[i]
}

func (s *ratedSubtitles) Less(i, j int) bool {
	return s.subs[i].score < s.subs[j].score
}

func (s *ratedSubtitles) String() string {
	var buffer bytes.Buffer

	for _, sub := range s.subs {
		buffer.WriteString(fmt.Sprintf("%-8.2f%v\n", sub.score, sub.Subtitle))
	}

	return buffer.String()
}
