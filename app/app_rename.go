package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

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

// Renamers holds the available renaming actions of the application
var Renamers = map[string]renamer{
	"copy":     renamer(copyRenamer),
	"move":     renamer(moveRenamer),
	"symlink":  renamer(symlinkRenamer),
	"hardlink": renamer(hardlinkRenamer),
}

var pathRegex = regexp.MustCompile(`[%/\?\\\*:\|"<>\n\r]`)

// cleanString cleans the string for unwanted characters such that it can
// be used safely as a name in a file hierarchy. All path seperators are
// removed from the string.
func cleanString(str string) string {
	return pathRegex.ReplaceAllString(str, "")
}

// RenameMedia traverses the local media list and renames the media
func (a *Application) RenameMedia(list types.LocalMediaList) error {

	doRename, ok := Renamers[viper.GetString("action")]

	if !ok {
		return fmt.Errorf("%s: unknown action", viper.GetString("action"))
	}

	for _, m := range list.List() {
		scraped, err := a.scrapeMedia(m)

		if err != nil {
			return err
		}

		if err := m.Merge(scraped); err != nil {
			return err
		}

		if movie, ok := m.TypeMovie(); ok {
			return a.renameMovie(m, movie, doRename)
		} else if episode, ok := m.TypeEpisode(); ok {
			return a.renameEpisode(m, episode, doRename)
		} else {
			return errors.New("unknown media format cannot rename")
		}
	}
	return nil
}

func (a *Application) scrapeMedia(m types.Media) (types.Media, error) {
	for _, s := range a.Scrapers() {
		scraped, err := s.Scrape(m)

		if err != nil {
			if provider.IsErrMediaNotSupported(err) {
				continue
			}
			return nil, err
		}
		return scraped, nil
	}
	return nil, errors.New("no scrapers to use for media")
}

func (a *Application) renameMovie(local types.Local, m types.Movie, rename renamer) error {
	var buf bytes.Buffer
	template := a.Config().Templates().Movies()
	data := struct {
		Movie   string
		Year    int
		Quality string
		Codec   string
		Group   string
	}{
		Movie:   cleanString(m.MovieName()),
		Year:    m.Year(),
		Quality: m.Quality().String(),
		Codec:   m.Codec().String(),
		Group:   cleanString(m.Group()),
	}
	if err := template.Execute(&buf, &data); err != nil {
		return err
	}
	return rename(local, buf.String()+filepath.Ext(local.Name()))
}

func (a *Application) renameEpisode(local types.Local, e types.Episode, rename renamer) error {
	var buf bytes.Buffer
	template := a.Config().Templates().TVShows()
	data := struct {
		TVShow  string
		Name    string
		Episode int
		Season  int
		Quality string
		Codec   string
		Group   string
	}{
		TVShow:  cleanString(e.TVShow()),
		Name:    cleanString(e.EpisodeName()),
		Episode: e.Episode(),
		Season:  e.Season(),
		Quality: e.Quality().String(),
		Codec:   e.Codec().String(),
		Group:   cleanString(e.Group()),
	}
	if err := template.Execute(&buf, &data); err != nil {
		return err
	}
	return rename(local, buf.String()+filepath.Ext(local.Name()))
}
