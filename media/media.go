package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

// File represents a local media file on disk
type File struct {
	os.FileInfo
	types.Media
	path string
}

// MarshalJSON returns the JSON represnetation of a media file
func (f *File) MarshalJSON() (b []byte, err error) {
	if js, ok := f.Media.Meta().(json.Marshaler); ok {
		return js.MarshalJSON()
	}
	return nil, errors.New("media could not be json encoded")
}

func (f *File) String() string {
	return f.Meta().String()
}

// Path returns the path to the media file
func (f *File) Path() string {
	return f.path
}

// ExistingSubtitles returns a list of existing subtitles for the media
func (f *File) ExistingSubtitles() (types.SubtitleList, error) {
	folder := filepath.Dir(f.Path())
	name := parse.Filename(f.Path())

	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	subtitles := make([]types.Subtitle, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasPrefix(file.Name(), name) {
			continue
		}
		sub, err := NewLocalSubtitle(file)
		if err != nil {
			continue
		}
		subtitles = append(subtitles, sub)
	}
	return list.Subtitles(subtitles...), nil
}

// SaveSubtitle saves the subtitle for the given media to disk
func (f *File) SaveSubtitle(s types.Downloadable, lang language.Tag) (types.LocalSubtitle, error) {
	if s == nil {
		return nil, errors.New("invalid subtitle nil")
	}

	srt, err := s.Download()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer srt.Close()

	name := fmt.Sprintf("%s.%s.%s", parse.Filename(f.Path()), lang, "srt")
	folder := filepath.Dir(f.Path())
	srtpath := filepath.Join(folder, name)

	file, err := os.Create(srtpath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	_, err = io.Copy(file, srt)

	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	sub := struct {
		FilePath
		os.FileInfo
		types.Subtitle
	}{
		FilePath(srtpath),
		info,
		s.(types.Subtitle),
	}

	return sub, nil
}

// NewFromPath parses a filepath into a media object. The path may be an
// absolute or relative path. The filename of the media must contain
// appropriate informations to describe the media file.
func NewFromPath(path string) (types.LocalMedia, error) {
	filename := filepath.Base(path)
	basename := strings.TrimSuffix(filename, filepath.Ext(filename))
	media, err := NewFromString(basename)

	if err != nil {
		return nil, err
	}

	file, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return &File{file, media, path}, nil
}

// NewFromString returns a media object parsed from a string describing the
// media. This could be the name of a file (without extension)
func NewFromString(str string) (types.Media, error) {
	if episodeRegexp.MatchString(str) {
		return NewEpisode(str)
	}
	return NewMovie(str)
}

// FilePath is a string describing a path to a file
type FilePath string

// Path return the path to a file
func (p FilePath) Path() string {
	return string(p)
}
