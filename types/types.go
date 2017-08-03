package types

import "golang.org/x/text/language"

// Subtitle can be downloaded
type Subtitle interface {
	Name() string
	IsLang(language.Tag) bool
	IsHI() bool
}

// Provider interfaces with subtitles websites
type Provider interface {
	Search(Media) ([]Subtitle, error)
}

// Media is an interface for movies and TV shows
type Media interface {
	Matches(Subtitle) bool
	Score(Subtitle) float32
}
