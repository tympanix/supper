package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/extract"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/types"
)

// FindArchives searches through filepaths and finds media archives
func (a *Application) FindArchives(roots ...string) ([]types.MediaArchive, error) {
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

// ExtractMedia reads a media file from a stream performs renaming/copying to
// its new destination
func (a *Application) ExtractMedia(m types.MediaReadCloser) error {
	ctx := log.WithField("media", m).WithField("action", "extract")

	if _, ok := m.TypeMovie(); !ok && viper.GetBool("movies") {
		return nil
	}

	if _, ok := m.TypeEpisode(); !ok && viper.GetBool("tvshows") {
		return nil
	}

	if _, ok := m.TypeSubtitle(); !ok && viper.GetBool("subtitles") {
		return nil
	}

	dest, err := a.scrapeAndRenameMedia(m, m)

	if err != nil {
		return err
	}

	if err = ensurePath(dest, a.Config().Force()); err != nil {
		if media.IsExistsErr(err) {
			ctx.WithField("reason", "media already exists").Warn("Extraction skipped")
		}
		return err
	}

	if !a.Config().Dry() {
		if err := copyMedia(m, dest); err != nil {
			ctx.WithError(err).Error("Extraction failed")
		} else {
			ctx.Info("Media extracted")
		}
	} else {
		ctx.WithField("reason", "dry-run").Info("Skip extraction")
	}

	return nil
}
