package extract

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

type ArchivedMedia struct {
	os.FileInfo
	types.Media
	Opener
}

type Opener func() (io.ReadCloser, error)

func (o Opener) Open() (io.ReadCloser, error) {
	return o()
}

func (a *ArchivedMedia) ExistingSubtitles() (types.SubtitleList, error) {
	return list.Subtitles(), nil
}

func (a *ArchivedMedia) SaveSubtitle(types.Downloadable, language.Tag) (types.LocalSubtitle, error) {
	return nil, errors.New("cannot save subtitle for archived media")
}

func Extract(path string) (types.MediaArchive, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".zip":
		return NewZipArchive(path)
	case ".rar":
		return NewRarArchive(path)
	}

	return nil, errors.New("file is not of any known archive formats")
}
