package types

import (
	"io"
	"os"

	"golang.org/x/text/language"
)

// Provider interfaces with subtitles websites
type Provider interface {
	SearchSubtitles(LocalMedia) (SubtitleCollection, error)
}

// Downloadable is an interface for media that can be downloaded from the internet
type Downloadable interface {
	Download() (io.ReadCloser, error)
}

// Evaluator determines how well the subtitle matches the media
type Evaluator interface {
	Evaluate(LocalMedia, Subtitle) float32
}

// Media is an interface for movies and TV shows
type Media interface {
	Meta() Metadata
	TypeMovie() (Movie, bool)
	TypeEpisode() (Episode, bool)
}

// Metadata is an interface metadata information
type Metadata interface {
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
	Path() string
	SaveSubtitle(Subtitle) error
}

// Movie interface is for movie type media material
type Movie interface {
	Metadata
	MovieName() string
	Year() int
}

// Episode interface is for TV show type material
type Episode interface {
	Metadata
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
