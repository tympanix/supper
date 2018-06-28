package media

import (
	"encoding/json"
	"os"

	"github.com/tympanix/supper/types"
)

// File represents a local media file on disk
type File struct {
	os.FileInfo
	types.Media
	FilePath
}

// MarshalJSON returns the JSON represnetation of a media file
func (f *File) MarshalJSON() (b []byte, err error) {
	return json.Marshal(f.Media)
}

// String returns a string representation of the media in the file
func (f *File) String() string {
	return f.Media.String()
}
