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
	return false
}
