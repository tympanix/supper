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
)

type File struct {
	os.FileInfo
	types.Media
	path string
}

func (f *File) MarshalJSON() (b []byte, err error) {
	if js, ok := f.Media.Meta().(json.Marshaler); ok {
		return js.MarshalJSON()
	}
	return nil, errors.New("media could not be json encoded")
}

func (f *File) String() string {
	return f.Meta().String()
}

func (f *File) Path() string {
	return f.path
}

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
func (f *File) SaveSubtitle(s types.OnlineSubtitle) error {
	if s == nil {
		return errors.New("invalid subtitle nil")
	}

	srt, err := s.Download()
	defer srt.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	name := fmt.Sprintf("%s.%s.%s", parse.Filename(f.Path()), s.Language(), "srt")
	folder := filepath.Dir(f.Path())
	srtpath := filepath.Join(folder, name)

	file, err := os.Create(srtpath)

	if err != nil {
		return err
	}

	defer file.Close()
	_, err = io.Copy(file, srt)

	if err != nil {
		return err
	}

	return nil
}

func NewFile(file os.FileInfo, meta types.Metadata, path string) *File {
	return &File{file, NewType(meta), path}
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
func New(path string) (types.LocalMedia, error) {
	filename := parse.Filename(path)
	media, err := NewMetadata(filename)

	if err != nil {
		return nil, err
	}

	file, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if movie, ok := media.(types.Movie); ok {
		return NewFile(file, movie, path), nil
	} else if episode, ok := media.(types.Episode); ok {
		return NewFile(file, episode, path), nil
	} else {
		return nil, errors.New("unknown media type")
	}
}

// NewMetadata returns a metadata object parsed from the string
func NewMetadata(str string) (types.Metadata, error) {
	if episodeRegexp.MatchString(str) {
		return NewEpisode(str)
	}
	return NewMovie(str)
}
