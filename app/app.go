package app

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tympanix/supper/api"
	"github.com/tympanix/supper/cfg"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/provider"
	"github.com/tympanix/supper/types"
)

var filetypes = []string{
	".avi", ".mkv", ".mp4", ".m4v", ".flv", ".mov", ".wmv", ".webm", ".mpg", ".mpeg",
}

// Application is an configuration instance of the application
type Application struct {
	types.Provider
	*http.ServeMux
	cfg types.Config
}

// New returns a new application from the cli context
func New(cfg types.Config) types.App {
	app := &Application{
		Provider: provider.Subscene(),
		cfg:      cfg,
		ServeMux: http.NewServeMux(),
	}

	static := viper.GetString("static")

	api := api.New(app)
	app.ServeMux.Handle("/api/", http.StripPrefix("/api", api))

	fs := WebAppHandler(static)
	app.ServeMux.Handle("/", fs)

	return app
}

// NewFromDefault construct an application using the default config
func NewFromDefault() types.App {
	return New(cfg.Default)
}

// Config returns the configuration for the application
func (a *Application) Config() types.Config {
	return a.cfg
}

// WebAppHandler serves a single-page web application
func WebAppHandler(path string) http.Handler {
	files := http.FileServer(http.Dir(path))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := filepath.Join(path, r.URL.Path)
		if _, err := os.Stat(uri); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(path, "index.html"))
		} else {
			files.ServeHTTP(w, r)
		}
	})
}

func fileIsMedia(f os.FileInfo) bool {
	for _, ext := range filetypes {
		if ext == filepath.Ext(f.Name()) {
			return true
		}
	}
	return false
}

// FindMedia searches for media files
func (a *Application) FindMedia(roots ...string) (types.LocalMediaList, error) {
	medialist := make([]types.LocalMedia, 0)

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			return nil, err
		}

		err := filepath.Walk(root, func(filepath string, f os.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}
			if !fileIsMedia(f) {
				return nil
			}
			_media, err := media.New(filepath)
			if err != nil {
				return nil
			}
			if media.IsSample(_media) {
				return nil
			}
			medialist = append(medialist, _media)
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return list.NewLocalMedia(medialist...), nil
}
