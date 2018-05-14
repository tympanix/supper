package media

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
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

func TestSaveSubtitle(t *testing.T) {
	f, err := NewLocalFile("test/Inception 2010 720p.mp4")
	require.NoError(t, err)
	v, ok := f.(types.Video)
	require.True(t, ok)

	sample := []byte("this is a test")

	buf := bytes.NewBuffer(sample)
	s, err := v.SaveSubtitle(buf, language.German)
	require.NoError(t, err)
	defer func() {
		os.Remove(s.Path())
	}()

	assert.Equal(t, language.German, s.Language())
	assert.True(t, strings.HasSuffix(s.Path(), ".de.srt"))
	assert.Contains(t, s.Path(), "Inception 2010 720p")

	in, err := os.Open(s.Path())
	require.NoError(t, err)

	defer in.Close()
	data, err := ioutil.ReadAll(in)
	require.NoError(t, err)
	assert.Equal(t, data, sample)
}
