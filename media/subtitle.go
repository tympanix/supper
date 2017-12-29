package media

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

func NewLocalSubtitle(file os.FileInfo) (types.Subtitle, error) {
	if filepath.Ext(file.Name()) != ".srt" {
		return nil, errors.New("parsing non subtitle file as subtitle")
	}

	parts := strings.Split(file.Name(), ".")

	if len(parts) < 2 {
		return nil, errors.New("error parsing subtitle file")
	}

	tag := language.Make(parts[len(parts)-2])

	return &LocalSubtitle{
		file,
		tag,
	}, nil
}

type LocalSubtitle struct {
	os.FileInfo
	lang language.Tag
}

func (l *LocalSubtitle) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		File string       `json:"filename"`
		Lang language.Tag `json:"language"`
	}{
		l.Name(),
		l.Language(),
	})
}

func (l *LocalSubtitle) IsHI() bool {
	return false
}

func (l *LocalSubtitle) Download() (io.ReadCloser, error) {
	return nil, errors.New("local subtitle can't be downloaded")
}

func (l *LocalSubtitle) Language() language.Tag {
	return l.lang
}

func (l *LocalSubtitle) IsLang(tag language.Tag) bool {
	return l.lang == tag
}

func (l *LocalSubtitle) Meta() types.Metadata {
	return nil
}

func (l *LocalSubtitle) TypeMovie() (types.Movie, bool) {
	return nil, false
}

func (l *LocalSubtitle) TypeEpisode() (types.Episode, bool) {
	return nil, false
}
