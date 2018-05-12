package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

func TestVideo(t *testing.T) {
	f, err := NewLocalFile("test/Inception 2010 720p.mp4")
	require.NoError(t, err)

	_, ok := f.TypeMovie()
	assert.True(t, ok)

	v, ok := f.(types.Video)
	require.True(t, ok)

	s, err := v.ExistingSubtitles()
	require.NoError(t, err)

	assert.Equal(t, 1, s.Len())
	assert.Equal(t, language.English, s.List()[0].Language())
}

func TestNoVideo(t *testing.T) {
	f, err := NewLocalFile("test/Inception 2010 720p.en.srt")
	require.NoError(t, err)

	_, ok := f.TypeSubtitle()
	assert.True(t, ok)

	_, ok = f.(types.Video)
	assert.False(t, ok)
}
