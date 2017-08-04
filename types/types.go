package types

import (
	"io"
	"os"

	"golang.org/x/text/language"
)

// Provider interfaces with subtitles websites
type Provider interface {
	SearchSubtitles(LocalMedia) ([]Subtitle, error)
}

// Downloadable is an interface for media that can be downloaded from the internet
type Downloadable interface {
	Download() io.Reader
}

// Media is an interface for movies and TV shows
type Media interface {
	Group() string
	Codec() string
	Quality() string
	Source() string
	AllTags() []string
}

// LocalMedia is an interface for media found locally on disk
type LocalMedia interface {
	os.FileInfo
	Media
}

// Movie interface is for movie type media material
type Movie interface {
	Media
	MovieName() string
	Year() int
}

// Episode interface is for TV show type material
type Episode interface {
	Media
	TVShow() string
	Episode() int
	Season() int
}

// Subtitle can be downloaded
type Subtitle interface {
	Media
	Downloadable
	IsLang(language.Tag) bool
	IsHI() bool
}
