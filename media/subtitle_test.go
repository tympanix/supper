package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tympanix/supper/meta/quality"
	"golang.org/x/text/language"
)

func TestSubtitles(t *testing.T) {
	s, err := NewSubtitle("Inception.2010.1080p.en")

	assert.NoError(t, err)
	assert.Equal(t, language.English, s.Language())

	_, ok := s.TypeSubtitle()
	assert.True(t, ok)

	m, ok := s.ForMedia().TypeMovie()
	assert.True(t, ok)
	assert.Equal(t, "Inception", m.MovieName())
	assert.Equal(t, 2010, m.Year())
	assert.Equal(t, quality.HD1080p, m.Quality())
}
