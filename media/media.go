package media

import (
	"os"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

// New parses a file into media attributes
func New(file os.FileInfo) types.LocalMedia {
	media := NewFromFilename(file.Name())

	if movie, ok := media.(types.Movie); ok {
		return &MovieFile{
			file,
			movie,
		}
	} else if episode, ok := media.(types.Episode); ok {
		return &EpisodeFile{
			file,
			episode,
		}
	} else {
		return nil
	}
}

// NewFromFilename parses media attributes from a files name
func NewFromFilename(filename string) types.Media {
	filename = parse.Filename(filename)
	if episodeRegexp.MatchString(filename) {
		return NewEpisode(filename)
	}
	return NewMovie(filename)
}
