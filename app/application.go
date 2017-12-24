package application

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Tympanix/supper/media"
	"github.com/Tympanix/supper/types"
)

var filetypes = []string{
	".avi", ".mkv", ".mp4", ".m4v", ".flv", ".mov", ".wmv", ".webm", ".mpg", ".mpeg",
}

// Application is an configuration instance of the application
type Application struct {
	types.Provider
}

func fileIsMedia(f os.FileInfo) bool {
	for _, ext := range filetypes {
		if ext == path.Ext(f.Name()) {
			return true
		}
	}
	return false
}

// FindMedia searches for media files
func (a *Application) FindMedia(root string) ([]types.LocalMedia, error) {
	medialist := make([]types.LocalMedia, 0)

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
			return fmt.Errorf("Cound not parse file: %s", filepath)
		}
		medialist = append(medialist, _media)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return medialist, nil
}
