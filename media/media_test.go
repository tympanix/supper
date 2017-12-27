package media_test

import (
	"testing"

	"github.com/tympanix/supper/media"
)

func TestTypes(t *testing.T) {
	m := media.NewType(nil)

	if _, ok := m.TypeMovie(); ok {
		t.Errorf("media is not a movie")
	}

	if _, ok := m.TypeEpisode(); ok {
		t.Errorf("media is not an episode")
	}
}
