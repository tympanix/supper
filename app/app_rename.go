package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/provider"
	"github.com/tympanix/supper/types"
)

type renamer func(types.Local, string) error

// Rename is a wrapper function around a renamer which performs some sanity checks
func (r renamer) Rename(local types.Local, dest string, force bool) error {
	if err := ensurePath(dest, force); err != nil {
		return err
	}
	return r(local, dest)
}

func ensurePath(dest string, force bool) error {
	_, err := os.Stat(dest)
	if !force && err == nil {
		return media.NewExistsErr()
	}
	if err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
		log.WithField("path", dest).Debug("Removed existing media")
	}
	if err := os.MkdirAll(filepath.Dir(dest), os.ModeDir); err != nil {
		return err
	}
	return nil
}

func copyMedia(file io.Reader, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	log.WithField("path", dest).Debug("Media copied")
	return nil
}

func copyRenamer(local types.Local, dest string) error {
	file, err := os.Open(local.Path())
	if err != nil {
		return err
	}
	defer file.Close()
	if err := copyMedia(file, dest); err != nil {
		return err
	}
	return nil
}

func moveRenamer(local types.Local, dest string) error {
	if err := os.Rename(local.Path(), dest); err != nil {
		return err
	}
	log.WithField("path", dest).Debug("Media moved")
	return nil
}

func symlinkRenamer(local types.Local, dest string) error {
	if err := os.Symlink(local.Path(), dest); err != nil {
		return err
	}
	log.WithField("path", dest).Debug("Media symlinked")
	return nil
}

func hardlinkRenamer(local types.Local, dest string) error {
	if err := os.Link(local.Path(), dest); err != nil {
		return err
	}
	log.WithField("path", dest).Debug("Media hardlinked")
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

var multispaceRegex = regexp.MustCompile(`\s\s+`)
var illegalPostfixRegex = regexp.MustCompile(`[^\p{L}\)0-9]+$`)

// truncateSpaces replaces all consecutive space characters with a single space.
// Trailing non-characters (i.e. spaces, symbols ect.) are removed
func truncateSpaces(str string) string {
	str = multispaceRegex.ReplaceAllString(str, " ")
	str = illegalPostfixRegex.ReplaceAllString(str, "")
	return str
}

// RenameMedia traverses the local media list and renames the media
func (a *Application) RenameMedia(list types.LocalMediaList) error {

	renamer, ok := Renamers[a.Config().RenameAction()]

	if !ok {
		return fmt.Errorf("%s: unknown action", a.Config().RenameAction())
	}

	for _, m := range list.List() {
		ctx := log.WithField("media", m).WithField("action", a.Config().RenameAction())

		dest, err := a.scrapeAndRenameMedia(m, m)

		if err != nil {
			return err
		}

		if !a.Config().Dry() {
			err = renamer.Rename(m, dest, a.Config().Force())
		} else {
			ctx.WithField("reason", "dry-run").Info("Skip rename")
			continue
		}

		if err != nil {
			if media.IsExistsErr(err) {
				ctx.WithField("reason", "media already exists").Warn("Rename skipped")
			} else {
				ctx.WithError(err).Error("Rename failed")
				if a.Config().Strict() {
					os.Exit(1)
				}
			}
		} else {
			ctx.Info("Media renamed")
		}
	}
	return nil
}

func (a *Application) scrapeAndRenameMedia(info os.FileInfo, m types.Media) (string, error) {
	scraped, err := a.scrapeMedia(m)

	if err != nil {
		return "", err
	}

	if err = m.Merge(scraped); err != nil {
		return "", err
	}

	dest, err := a.renameMedia(info, m)

	if media.IsUnknown(err) && a.Config().Strict() {
		return "", err
	}

	if err != nil && !media.IsUnknown(err) {
		return "", err
	}

	return dest, nil
}

func (a *Application) renameMedia(info os.FileInfo, m types.Media) (dest string, err error) {
	if movie, ok := m.TypeMovie(); ok {
		return a.renameMovie(info, movie)
	} else if episode, ok := m.TypeEpisode(); ok {
		return a.renameEpisode(info, episode)
	} else if sub, ok := m.TypeSubtitle(); ok {
		return a.renameSubtitle(info, sub)
	}
	return "", media.NewUnknownErr()
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

func (a *Application) renameMovie(info os.FileInfo, m types.Movie) (string, error) {
	var buf bytes.Buffer
	template := a.Config().Movies().Template()
	if template == nil {
		return "", errors.New("missing template for movies")
	}
	data := struct {
		Movie   string
		Year    int
		Quality string
		Codec   string
		Source  string
		Group   string
	}{
		Movie:   cleanString(m.MovieName()),
		Year:    m.Year(),
		Quality: m.Quality().String(),
		Codec:   m.Codec().String(),
		Source:  m.Source().String(),
		Group:   cleanString(m.Group()),
	}
	if err := template.Execute(&buf, &data); err != nil {
		return "", err
	}
	filename := truncateSpaces(buf.String()) + filepath.Ext(info.Name())
	return filepath.Join(a.Config().Movies().Directory(), filename), nil
}

func (a *Application) renameEpisode(info os.FileInfo, e types.Episode) (string, error) {
	var buf bytes.Buffer
	template := a.Config().TVShows().Template()
	if template == nil {
		return "", errors.New("missing template for tvshows")
	}
	data := struct {
		TVShow  string
		Name    string
		Episode int
		Season  int
		Quality string
		Codec   string
		Source  string
		Group   string
	}{
		TVShow:  cleanString(e.TVShow()),
		Name:    cleanString(e.EpisodeName()),
		Episode: e.Episode(),
		Season:  e.Season(),
		Quality: e.Quality().String(),
		Codec:   e.Codec().String(),
		Source:  e.Source().String(),
		Group:   cleanString(e.Group()),
	}
	if err := template.Execute(&buf, &data); err != nil {
		return "", err
	}
	filename := truncateSpaces(buf.String()) + filepath.Ext(info.Name())
	return filepath.Join(a.Config().TVShows().Directory(), filename), nil
}

func (a *Application) renameSubtitle(info os.FileInfo, s types.Subtitle) (string, error) {
	dest, err := a.renameMedia(info, s.ForMedia())

	if err != nil {
		return "", err
	}

	ext := filepath.Ext(dest)
	base := strings.TrimSuffix(dest, ext)
	dest = base + "." + fmt.Sprintf("%v", s.Language()) + ext

	return dest, err
}
