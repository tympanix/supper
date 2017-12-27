package types

import (
	"time"

	"golang.org/x/text/language"

	"github.com/fatih/set"
)

type List interface {
	Len() int
}

// SubtitleList is a collection of subtitles
type SubtitleList interface {
	List
	Add(Subtitle)
	Best() Subtitle
	List() []Subtitle
	LanguageSet() set.Interface
	FilterLanguage(language.Tag) SubtitleList
	HearingImpaired(bool) SubtitleList
}

// MediaList is a list of media which can be manipulated
type LocalMediaList interface {
	List
	Add(LocalMedia)
	List() []LocalMedia
	FilterModified(time.Duration) LocalMediaList
	FilterMissingSubs(set.Interface) (LocalMediaList, error)
}
