package extract

import (
	"io"
	"io/ioutil"

	"github.com/nwaples/rardecode"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

type RarArchive struct {
	*rardecode.ReadCloser
}

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

	return nil, io.EOF
}

func NewRarArchive(path string) (types.MediaArchive, error) {
	r, err := rardecode.OpenReader(path, "")

	if err != nil {
		return nil, err
	}

	return &RarArchive{
		ReadCloser: r,
	}, nil
}
