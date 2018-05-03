package types

import (
	"time"

	"golang.org/x/text/language"

	"github.com/fatih/set"
)

// List interface describes all common properties of lists
type List interface {
	Len() int
}

// SubtitleList is a collection of subtitles
type SubtitleList interface {
	List
	Add(...Subtitle)
	Best() (Subtitle, float32)
	List() []Subtitle
	LanguageSet() set.Interface
	FilterLanguage(language.Tag) SubtitleList
	FilterScore(float32) SubtitleList
	HearingImpaired(bool) SubtitleList
}

// LocalMediaList is a list of media which can be manipulated
type LocalMediaList interface {
	List
	Add(LocalMedia)
	List() []LocalMedia
	FilterModified(time.Duration) LocalMediaList
	FilterVideo() VideoList
}

// VideoList is a list of video
type VideoList interface {
	List
	List() []Video
	FilterMissingSubs(set.Interface) (VideoList, error)
}
