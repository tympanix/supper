package extract

import (
	"io/ioutil"

	"github.com/nwaples/rardecode"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

// RarArchive is an struct used to extract media files from rar archives
type RarArchive struct {
	*rardecode.ReadCloser
}

// Next returns the next media item in the rar archive. If there are not more
// media files io.EOF is returned
func (r *RarArchive) Next() (types.MediaReadCloser, error) {
	file, err := r.ReadCloser.Next()

	if err != nil {
		return nil, err
	}

	med, err := media.NewFromString(parse.Filename(file.Name))

	if err != nil {
		return r.Next()
	}

	return &MediaFromArchive{
		Media:      med,
		ReadCloser: ioutil.NopCloser(r.ReadCloser),
	}, nil
}

// NewRarArchive creates a new rar archive object to extract media from
func NewRarArchive(path string) (types.MediaArchive, error) {
	r, err := rardecode.OpenReader(path, "")

	if err != nil {
		return nil, err
	}

	return &RarArchive{
		ReadCloser: r,
	}, nil
}
