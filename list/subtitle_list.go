package list

import (
	"github.com/fatih/set"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

// Subtitles creates a new list of subtitles from variadic argument
func Subtitles(subs ...types.Subtitle) types.SubtitleList {
	list := subtitleList(subs)
	return &list
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

// HearingImpaired returnes a new subtitle collection where hearing impared subtitles has been filtered
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

func (s *subtitleList) LanguageSet() set.Interface {
	langs := set.New()
	for _, sub := range *s {
		langs.Add(sub.Language())
	}
	return langs
}
