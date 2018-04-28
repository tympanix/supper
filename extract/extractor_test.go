package extract

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipArchive(t *testing.T) {
	media, err := Extract("test/movie_test.zip")

	assert.NoError(t, err)

	med, err := media.Next()
	assert.NoError(t, err)

	movie, ok := med.TypeMovie()
	assert.True(t, ok)
	assert.Equal(t, "Blade Runner", movie.MovieName())
	assert.Equal(t, 2017, movie.Year())

	med, err = media.Next()
	assert.Equal(t, io.EOF, err)
}
