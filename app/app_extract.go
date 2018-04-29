package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/tympanix/supper/extract"
	"github.com/tympanix/supper/types"
)

// ExtractMedia searches through filepaths and finds media archives
func (a *Application) ExtractMedia(roots ...string) ([]types.MediaArchive, error) {
	archives := make([]types.MediaArchive, 0)

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			return nil, err
		}

		err := filepath.Walk(root, func(filepath string, f os.FileInfo, err error) error {
			if f == nil {
				return errors.New("invalid file path")
			}
			if f.IsDir() {
				return nil
			}
			archive, err := extract.OpenMediaArchive(filepath)
			if extract.IsNotArchive(err) {
				return nil
			}
			if err != nil {
				return err
			}
			archives = append(archives, archive)
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return archives, nil
}
