package extract

import (
	"archive/zip"
	"io"

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

// ZipArchive is an struct used to extract media files from zip archives
type ZipArchive struct {
	*zip.ReadCloser
	idx int
}

// Next returns the next media item in the zip archive. If there are not more
// media files io.EOF is returned
func (z *ZipArchive) Next() (types.MediaReadCloser, error) {
	for i := z.idx; i < len(z.File); i++ {
		z.idx = i + 1
		file := z.File[i]
		med, err := media.NewFromString(parse.Filename(file.Name))
		if err != nil {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		return &MediaFromArchive{
			FileInfo:   file.FileInfo(),
			Media:      med,
			ReadCloser: rc,
		}, nil
	}
	return nil, io.EOF
}

// NewZipArchive creates a new zip archive object to extract media from
func NewZipArchive(path string) (types.MediaArchive, error) {
	r, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &ZipArchive{
		ReadCloser: r,
		idx:        0,
	}, nil
}
