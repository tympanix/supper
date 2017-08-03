package media

import (
	"os"

	"github.com/Tympanix/supper/types"
)

// New parses a file into media attributes
func New(file *os.File) types.Media {
	if episodeRegexp.MatchString(file.Name()) {
		return NewEpisode(file)
	}
	return NewMovie(file)
}
