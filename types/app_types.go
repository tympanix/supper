package types

import (
	"html/template"
	"net/http"
	"time"

	"github.com/fatih/set"
)

// App is the interface for the top level capabilities of the application.
// It is an HTTP handler, a provider (for subtitles) and a CLI application.
// It means App can both be used as a HTTP server and a CLI application.
type App interface {
	Provider
	http.Handler
	Config() Config
	Scrapers() []Scraper
	FindMedia(...string) (LocalMediaList, error)
	DownloadSubtitles(LocalMediaList, set.Interface) (int, error)
	RenameMedia(LocalMediaList) error
	FindArchives(...string) ([]MediaArchive, error)
	ExtractMedia(MediaReadCloser) error
}

// Config is the interface for application configuration
type Config interface {
	Languages() set.Interface
	Impaired() bool
	Limit() int
	Modified() time.Duration
	Dry() bool
	Score() int
	Delay() time.Duration
	Force() bool
	Config() string
	Logfile() string
	Verbose() bool
	Strict() bool
	Plugins() []Plugin
	APIKeys() APIKeys
	Movies() MediaConfig
	TVShows() MediaConfig
	MediaFilter() MediaFilter
}

// APIKeys is the interface for configuration of 3rd party APIs
type APIKeys interface {
	TheTVDB() string
	TheMovieDB() string
}

// MediaConfig is the configuration interface for media collections
type MediaConfig interface {
	Directory() string
	Template() *template.Template
}

// Plugin is an interface for external functionality
type Plugin interface {
	Name() string
	Run(LocalSubtitle) error
}
