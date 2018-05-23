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
	List() []Subtitle
	LanguageSet() set.Interface
	FilterLanguage(language.Tag) SubtitleList
	HearingImpaired(bool) SubtitleList
	RateByMedia(Media) RatedSubtitleList
}

// RatedSubtitleList is a collection of subtitle ordered by rating
type RatedSubtitleList interface {
	List
	List() []RatedSubtitle
	Best() RatedSubtitle
	FilterScore(float32) RatedSubtitleList
}

// LocalMediaList is a list of media which can be manipulated
type LocalMediaList interface {
	List
	Add(LocalMedia)
	List() []LocalMedia
	Filter(MediaFilter) LocalMediaList
	FilterModified(time.Duration) LocalMediaList
	FilterVideo() VideoList
	FilterMovies() LocalMediaList
	FilterEpisodes() LocalMediaList
	FilterSubtitles() LocalMediaList
}

// MediaFilter is used to filter out local media
type MediaFilter func(Media) bool

// VideoList is a list of video
type VideoList interface {
	List
	List() []Video
	FilterMissingSubs(set.Interface) (VideoList, error)
}
