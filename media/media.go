package media

import (
	"errors"
	"os"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

type File struct {
	os.FileInfo
	types.Media
}

func NewFile(file os.FileInfo, meta types.Metadata) *File {
	return &File{file, NewType(meta)}
}

type Type struct {
	types.Metadata
}

func (m *Type) Meta() types.Metadata {
	return m.Metadata
}

func (m *Type) TypeMovie() (r types.Movie, ok bool) {
	r, ok = m.Metadata.(types.Movie)
	return
}

func (m *Type) TypeEpisode() (r types.Episode, ok bool) {
	r, ok = m.Metadata.(types.Episode)
	return
}

func NewType(m types.Metadata) *Type {
	return &Type{m}
}

// New parses a file into media attributes
func New(file os.FileInfo) (types.LocalMedia, error) {
	filename := parse.Filename(file.Name())
	media, err := NewMetadata(filename)

	if err != nil {
		return nil, err
	}

	if movie, ok := media.(types.Movie); ok {
		return NewFile(file, movie), nil
	} else if episode, ok := media.(types.Episode); ok {
		return NewFile(file, episode), nil
	} else {
		return nil, errors.New("Unknown media type")
	}
}

// NewMetadata returns a metadata object parsed from the string
func NewMetadata(str string) (types.Metadata, error) {
	if episodeRegexp.MatchString(str) {
		return NewEpisode(str)
	}
	return NewMovie(str)
}
