package types

import (
	"os"

	"golang.org/x/text/language"
)

// Subtitle can be downloaded
type Subtitle interface {
	Name() string
	Download() *os.File
	IsLang(language.Tag) bool
	IsHI() bool
}

// Provider interfaces with subtitles websites
type Provider interface {
	Search(string) ([]Subtitle, error)
}

// Media is an interface for movies and TV shows
type Media interface {
	Matches(Subtitle) bool
	Score(Subtitle) float32
}
