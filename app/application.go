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
	".avi", ".mkv", ".mp4", ".m4a", ".flv",
}

// Application is an configuration instance of the application
type Application struct {
	types.Provider
}

// FindMedia searches for media files
func (a *Application) FindMedia(root string) ([]types.Media, error) {
	medialist := make([]types.Media, 0)

	err := filepath.Walk(root, func(filepath string, f os.FileInfo, err error) error {
		for _, ext := range filetypes {
			if ext == path.Ext(filepath) {
				_media := media.New(f)
				if _media == nil {
					return fmt.Errorf("Cound not parse file: %s", filepath)
				}
				medialist = append(medialist, _media)
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return medialist, nil
}
