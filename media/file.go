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
	FilePath
}

// MarshalJSON returns the JSON represnetation of a media file
func (f *File) MarshalJSON() (b []byte, err error) {
	if js, ok := f.Media.Meta().(json.Marshaler); ok {
		return js.MarshalJSON()
	}
	return nil, errors.New("media could not be json encoded")
}

// String returns a string representation of the media in the file
func (f *File) String() string {
	return f.Media.String()
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
