package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tympanix/supper/provider"
	"github.com/tympanix/supper/types"
)

type renamer func(types.Local, string) error

func copyRenamer(local types.Local, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir); err != nil {
		return err
	}
	file, err := os.Open(local.Path())
	if err != nil {
		return err
	}
	defer file.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

func moveRenamer(local types.Local, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir); err != nil {
		return err
	}
	if err := os.Rename(local.Path(), dest); err != nil {
		return err
	}
	return nil
}

func symlinkRenamer(local types.Local, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir); err != nil {
		return err
	}
	if err := os.Symlink(local.Path(), dest); err != nil {
		return err
	}
	return nil
}

func hardlinkRenamer(local types.Local, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir); err != nil {
		return err
	}
	if err := os.Link(local.Path(), dest); err != nil {
		return err
	}
	return nil
}

var Renamers = map[string]renamer{
	"copy":     renamer(copyRenamer),
	"move":     renamer(moveRenamer),
	"symlink":  renamer(symlinkRenamer),
	"hardlink": renamer(hardlinkRenamer),
}

// RenameMedia traverses the local media list and renames the media
func (a *Application) RenameMedia(list types.LocalMediaList) error {

	doRename, ok := Renamers[viper.GetString("action")]

	if !ok {
		return fmt.Errorf("%s: unknown action", viper.GetString("action"))
	}

	templates := struct {
		Movies  string `mapstructure:"movies"`
		TVShows string `mapstructure:"tvshows"`
	}{}

	if err := viper.UnmarshalKey("templates", &templates); err != nil {
		return err
	}

	for _, m := range list.List() {
		var scraped types.Media
		for _, s := range a.Scrapers() {
			var err error
			scraped, err = s.Scrape(m)

			if err != nil {
				if provider.IsErrMediaNotSupported(err) {
					continue
				}
				return err
			}
		}

		fmt.Println(scraped)

		if movie, ok := m.TypeMovie(); ok {
			return a.renameMovie(m, movie, doRename, templates.Movies)
		} else if episode, ok := m.TypeEpisode(); ok {
			return a.renameEpisode(m, episode, doRename, templates.TVShows)
		} else {
			return errors.New("unknown media format cannot rename")
		}
	}
	return nil
}

func (a *Application) renameMovie(local types.Local, m types.Movie, rename renamer, template string) error {
	folder := fmt.Sprintf("%v (%v)", m.MovieName(), m.Year())
	media := fmt.Sprintf("%v.%v", folder, filepath.Ext(local.Path()))
	return rename(local, filepath.Join(folder, media))
}

func (a *Application) renameEpisode(local types.Local, m types.Episode, rename renamer, template string) error {
	showFolder := fmt.Sprintf("%v", m.TVShow())
	seasonFolder := fmt.Sprintf("Season %02d", m.Season())
	mediaFile := fmt.Sprintf("%v S%02dE%02d.%v", m.TVShow(), m.Season(), m.Episode(), filepath.Ext(local.Path()))
	return rename(local, filepath.Join(showFolder, seasonFolder, mediaFile))
}
