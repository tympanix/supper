package media

import (
	"strings"

	"github.com/tympanix/supper/types"
)

func IsSample(m types.Media) bool {
	for _, t := range m.Meta().AllTags() {
		if strings.ToLower(t) == "sample" {
			return true
		}
	}
	return false
}
