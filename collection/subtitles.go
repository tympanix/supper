package collection

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/Tympanix/supper/types"
)

type subtitleEntry struct {
	types.Subtitle
	score float32
}

// NewSubtitles returns a new subtitles collection
func NewSubtitles(media types.LocalMedia) *Subtitles {
	return &Subtitles{
		Evaluator: new(DefaultEvaluator),
		media:     media,
		subs:      make([]subtitleEntry, 0),
	}
}

// Subtitles is a sortable and filterable collection of subtitles
type Subtitles struct {
	types.Evaluator
	media types.LocalMedia
	subs  []subtitleEntry
}

// Add a subtitle to the collection
func (s *Subtitles) Add(sub types.Subtitle) {
	if sub == nil || sub.Meta() == nil {
		return
	}
	s.subs = append(s.subs, subtitleEntry{
		Subtitle: sub,
		score:    0, //s.Evaluate(s.media, sub),
	})
	sort.Sort(s)
}

// RemoveHI removes all HI subtitle from the collection
func (s *Subtitles) RemoveHI() {
	for i, sub := range s.subs {
		if sub.IsHI() {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
		}
	}
}

// RemoveNotHI removes all normal subtitles from the collection
func (s *Subtitles) RemoveNotHI() {
	for i, sub := range s.subs {
		if !sub.IsHI() {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
		}
	}
}

func (s *Subtitles) Len() int {
	return len(s.subs)
}

func (s *Subtitles) Swap(i, j int) {
	s.subs[i], s.subs[j] = s.subs[j], s.subs[i]
}

func (s *Subtitles) Less(i, j int) bool {
	return s.subs[i].score < s.subs[j].score
}

func (s *Subtitles) String() string {
	var buffer bytes.Buffer

	for _, sub := range s.subs {
		buffer.WriteString(fmt.Sprintf("%-8.2f%v\n", sub.score, sub.Subtitle))
	}

	return buffer.String()
}
