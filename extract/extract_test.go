package extract

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestZipArchive(t *testing.T) {
	media, err := OpenMediaArchive("test/movie_test.zip")

	assert.NoError(t, err)
	assert.False(t, IsNotArchive(err))

	defer media.Close()

	med, err := media.Next()
	assert.NoError(t, err)

	movie, ok := med.TypeMovie()
	assert.True(t, ok)
	assert.Equal(t, "Blade Runner 2049", movie.MovieName())
	assert.Equal(t, 2017, movie.Year())

	med, err = media.Next()
	assert.NoError(t, err)

	sub, ok := med.TypeSubtitle()
	require.True(t, ok)
	assert.Equal(t, language.English, sub.Language())

	med, err = media.Next()
	assert.Equal(t, io.EOF, err)
}

func TestRarArchive(t *testing.T) {
	media, err := OpenMediaArchive("test/movie_test.rar")

	assert.NoError(t, err)
	assert.False(t, IsNotArchive(err))

	defer media.Close()

	med, err := media.Next()
	assert.NoError(t, err)

	movie, ok := med.TypeMovie()
	assert.True(t, ok)
	assert.Equal(t, "Fight Club", movie.MovieName())
	assert.Equal(t, 1999, movie.Year())

	med, err = media.Next()
	assert.NoError(t, err)

	s, ok := med.TypeSubtitle()
	require.True(t, ok)
	assert.Equal(t, language.English, s.Language())

	med, err = media.Next()
	assert.Equal(t, io.EOF, err)
}

func TestNotArchive(t *testing.T) {
	_, err := OpenMediaArchive("not/an/archive.nope")
	assert.Error(t, err)
	assert.True(t, IsNotArchive(err))
}
