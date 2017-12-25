package types

import (
	"time"

	"golang.org/x/text/language"
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
	FilterLanguage(language.Tag) SubtitleList
	HearingImpaired(bool) SubtitleList
}

// MediaList is a list of media which can be manipulated
type LocalMediaList interface {
	List
	Add(LocalMedia)
	List() []LocalMedia
	FilterModified(time.Duration) LocalMediaList
}
