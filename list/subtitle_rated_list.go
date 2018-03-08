package list

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/language"

	"github.com/fatih/set"
	"github.com/tympanix/supper/score"
	"github.com/tympanix/supper/types"
)

type subtitleEntry struct {
	types.OnlineSubtitle
	score float32
}

func (s subtitleEntry) MarshalJSON() (b []byte, err error) {
	hash := sha1.New()

	info := []string{
		s.Link(),
		s.Meta().Codec().String(),
		s.Meta().Group(),
		s.Meta().Quality().String(),
		s.Meta().Source().String(),
	}

	hash.Write([]byte(strings.Join(info, "")))
	hashval := hash.Sum(nil)
	infohash := make([]byte, hex.EncodedLen(len(hashval)))
	hex.Encode(infohash, hashval)

	return json.Marshal(struct {
		Hash  string         `json:"hash"`
		Lang  language.Tag   `json:"language"`
		Link  string         `json:"link"`
		Score float32        `json:"score"`
		HI    bool           `json:"hi"`
		Media types.Metadata `json:"media"`
	}{
		string(infohash),
		s.Language(),
		s.Link(),
		s.score,
		s.IsHI(),
		s.Meta(),
	})
}

// RatedSubtitles returns a new subtitles collection
func RatedSubtitles(media types.LocalMedia) types.SubtitleList {
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

func (s *ratedSubtitles) MarshalJSON() (b []byte, err error) {
	return json.Marshal(s.List())
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
func (s *ratedSubtitles) Add(subs ...types.Subtitle) {
	for _, sub := range subs {
		sub, ok := sub.(types.OnlineSubtitle)
		if !ok {
			panic("rated subtitle list only supports online subtitles")
		}
		score := s.Evaluate(s.media, sub)
		if score > 0 {
			s.subs = append(s.subs, subtitleEntry{
				OnlineSubtitle: sub,
				score:          score,
			})
		}
	}
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
		buffer.WriteString(fmt.Sprintf("%-8.2f%v\n", sub.score, sub.OnlineSubtitle))
	}

	return buffer.String()
}
