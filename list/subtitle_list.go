package list

import (
	"fmt"
	"reflect"

	"github.com/fatih/set"
	"github.com/tympanix/supper/score"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

// Subtitles creates a new list of subtitles from variadic argument
func Subtitles(subs ...types.Subtitle) types.SubtitleList {
	list := subtitleList(subs)
	return &list
}

// NewSubtitlesFromInterface construct a subtitle list from interface values
func NewSubtitlesFromInterface(subs ...interface{}) (l types.SubtitleList, err error) {
	var list []types.Subtitle
	defer func() {
		if r := recover(); r != nil {
			l, err = nil, fmt.Errorf("%v", r)
		}
	}()
	for _, v := range subs {
		list = append(list, extractSubtitles(v)...)
	}
	lx := subtitleList(list)
	return &lx, nil
}

func extractSubtitles(e interface{}) []types.Subtitle {
	if s, ok := e.(types.Subtitle); ok {
		return []types.Subtitle{s}
	}
	t := reflect.TypeOf(e).Kind()
	if t == reflect.Slice || t == reflect.Array {
		var subs []types.Subtitle
		v := reflect.ValueOf(e)
		for j := 0; j < v.Len(); j++ {
			subs = append(subs, extractSubtitles(v.Index(j).Interface())...)
		}
		return subs
	}

	panic(fmt.Sprintf("Unknown subtitle format %v", reflect.TypeOf(e)))
}

type subtitleList []types.Subtitle

func (s *subtitleList) Len() int {
	return len(*s)
}

func (s *subtitleList) List() []types.Subtitle {
	return *s
}

func (s *subtitleList) Best() (types.Subtitle, float32) {
	return nil, -1
}

func (s *subtitleList) Add(sub ...types.Subtitle) {
	*s = append(*s, sub...)
}

// FilterLanguage returns a new subtitle collection including only the argument language
func (s *subtitleList) FilterLanguage(lang language.Tag) types.SubtitleList {
	_subs := make([]types.Subtitle, 0)
	for _, sub := range *s {
		if sub.Language() == lang {
			_subs = append(_subs, sub)
		}
	}
	list := subtitleList(_subs)
	return &list
}

func (s *subtitleList) FilterScore(score float32) types.SubtitleList {
	panic("unrated subtitle list does not support filtering by score")
}

// HearingImpaired returns a new subtitle collection where hearing impared subtitles has been filtered
func (s *subtitleList) HearingImpaired(hi bool) types.SubtitleList {
	_subs := make([]types.Subtitle, 0)
	for _, sub := range *s {
		if sub.HearingImpaired() == hi {
			_subs = append(_subs, sub)
		}
	}
	list := subtitleList(_subs)
	return &list
}

// RateByMedia returns a rated subtitle list, where every subtitle has been
// given a score according to how well it matches the argument media
func (s *subtitleList) RateByMedia(m types.Media) types.RatedSubtitleList {
	return NewRatedSubtitles(m, new(score.DefaultEvaluator), (*s)...)
}

func (s *subtitleList) LanguageSet() set.Interface {
	langs := set.New()
	for _, sub := range *s {
		langs.Add(sub.Language())
	}
	return langs
}
