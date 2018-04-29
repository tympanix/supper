package extract

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tympanix/supper/types"
)

// ErrNotArchive is an error which indicates the file was not an archive
type ErrNotArchive struct {
	error
}

// IsNotArchive returns true if the error is of ErrNotArchive type
func IsNotArchive(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrNotArchive)
	return ok
}

// MediaFromArchive representa a media item which is located inside a
// compressed archive
type MediaFromArchive struct {
	os.FileInfo
	types.Media
	io.ReadCloser
}

// OpenMediaArchive opens the file as an archive and exposes the media files
// within. If the file is not recognized as any of the known archive formats,
// an error is returned
func OpenMediaArchive(path string) (types.MediaArchive, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".zip":
		return NewZipArchive(path)
	case ".rar":
		return NewRarArchive(path)
	}

	return nil, &ErrNotArchive{
		fmt.Errorf("%s: not of any known archive formats", path),
	}
}
