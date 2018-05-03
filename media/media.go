package media

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

// videoExt contains the list of recognizable video extensions
var videoExt = []string{
	".avi", ".mkv", ".mp4", ".m4v", ".flv", ".mov", ".wmv", ".webm", ".mpg", ".mpeg",
}

// subExt contains the list of recognizable subtitle extensions
var subExt = []string{
	".srt",
}

func fileIsVideo(name string) bool {
	for _, ext := range videoExt {
		if ext == filepath.Ext(name) {
			return true
		}
	}
	return false
}

func fileIsSubtitle(name string) bool {
	for _, ext := range subExt {
		if ext == filepath.Ext(name) {
			return true
		}
	}
	return false
}

// NewLocalFile parses a filepath into a local media object. The path may be an
// absolute or relative path. The filename of the media must contain
// appropriate information to describe the media file.
func NewLocalFile(path string) (types.LocalMedia, error) {
	filename := filepath.Base(path)
	media, err := NewFromFilename(filename)

	if err != nil {
		return nil, err
	}

	file, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return &File{
		FileInfo: file,
		Media:    media,
		FilePath: FilePath(path),
	}, nil
}

// NewFromFilename parses the filename and returns a media object. The filename
// (with extenstion) may describe either some video material (.avi, .mkv, .mp4)
// or a subtitle (.srt).
func NewFromFilename(name string) (types.Media, error) {
	filename := parse.Filename(name)
	if fileIsVideo(name) {
		return NewFromString(filename)
	} else if fileIsSubtitle(name) {
		return NewSubtitle(filename)
	}
	return nil, errors.New("could not parse filename into media")
}

// NewFromString returns a media object parsed from a string describing the
// media. This could be the name of a file (without extension). It is assumed
// the string describes some video material (movie or episode)
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

// Open opens the file and returns a readcloser
func (p FilePath) Open() (io.ReadCloser, error) {
	return os.Open(string(p))
}
