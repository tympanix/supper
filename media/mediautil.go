package media

import (
	"strings"

	"github.com/tympanix/supper/types"
)

// IsSample return wether the media is a sample video of the real media. The
// media is a sample if any of the tags contains "sample" (case insensitive)
func IsSample(m types.Media) bool {
	for _, t := range m.Meta().AllTags() {
		if strings.ToLower(t) == "sample" {
			return true
		}
	}
	if strings.HasPrefix(strings.ToLower(m.String()), "sample") {
		return true
	}
	return false
}

// TypeNone represents a media of unknown format. Can be used for embedding into
// other types to fill common mundane methods
type TypeNone struct{}

// TypeMovie return false, since TypeNone is of an unknown media type
func (t TypeNone) TypeMovie() (types.Movie, bool) {
	return nil, false
}

// TypeSubtitle return false, since TypeNone is of an unknown media type
func (t TypeNone) TypeSubtitle() (types.Subtitle, bool) {
	return nil, false
}

// TypeEpisode return false, since TypeNone is of an unknown media type
func (t TypeNone) TypeEpisode() (types.Episode, bool) {
	return nil, false
}
