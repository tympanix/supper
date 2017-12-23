package types

import (
	"golang.org/x/text/language"
)

// SubtitleCollection is a collection of subtitles
type SubtitleCollection interface {
	RemoveNotHI()
	RemoveHI()
	Add(Subtitle)
	FilterLanguage(language.Tag)
	Best() Subtitle
}
