package media

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/media/meta/codec"
	"github.com/tympanix/supper/media/meta/quality"
	"github.com/tympanix/supper/media/meta/source"
)

func TestMetadata(t *testing.T) {
	m := ParseMetadata("DivX 720p BluRay")
	assert.Equal(t, quality.HD720p, m.Quality())
	assert.Equal(t, codec.DivX, m.Codec())
	assert.Equal(t, source.BluRay, m.Source())

	assert.Contains(t, m.String(), quality.HD720p.String())
	assert.Contains(t, m.String(), source.BluRay.String())
	assert.Contains(t, m.String(), codec.DivX.String())

	assert.Contains(t, m.AllTags(), "DivX")
	assert.Contains(t, m.AllTags(), "720p")
	assert.Contains(t, m.AllTags(), "BluRay")
}

func TestMetadataJSON(t *testing.T) {
	m := ParseMetadata("x264 1080p DVD-Rip GROUP")

	data, err := json.Marshal(m)
	require.NoError(t, err)

	j := struct {
		Quality string `json:"quality"`
		Source  string `json:"source"`
		Codec   string `json:"codec"`
		Group   string `json:"group"`
	}{}

	err = json.Unmarshal(data, &j)
	require.NoError(t, err)

	assert.Equal(t, j.Quality, quality.HD1080p.String())
	assert.Equal(t, j.Source, source.DVDRip.String())
	assert.Equal(t, j.Codec, codec.X264.String())
	assert.Equal(t, j.Group, "GROUP")
}
