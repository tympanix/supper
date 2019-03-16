package media

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/media/meta/quality"
	"golang.org/x/text/language"
)

func TestSubtitles(t *testing.T) {
	s, err := NewSubtitle("Inception.2010.1080p.en")

	assert.NoError(t, err)
	assert.Equal(t, language.English, s.Language())
	assert.False(t, s.HearingImpaired())

	_, ok := s.TypeSubtitle()
	assert.True(t, ok)

	m, ok := s.ForMedia().TypeMovie()
	assert.True(t, ok)
	assert.Equal(t, "Inception", m.MovieName())
	assert.Equal(t, 2010, m.Year())
	assert.Equal(t, quality.HD1080p, m.Quality())

	assert.Contains(t, s.Identity(), s.ForMedia().Identity())

}

func TestLocalSubtitleError(t *testing.T) {
	s, err := NewLocalSubtitle("test/Test.en.srt")
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestLocalSubtitleJSON(t *testing.T) {
	s, err := NewLocalSubtitle("test/Inception 2010 720p.en.srt")
	require.NoError(t, err)

	data, err := json.Marshal(s)
	require.NoError(t, err)

	j := struct {
		Filename string `json:"filename"`
		Code     string `json:"code"`
		Language string `json:"language"`
	}{}

	err = json.Unmarshal(data, &j)
	require.NoError(t, err)

	assert.Equal(t, "Inception 2010 720p.en.srt", j.Filename)
	assert.Equal(t, "English", j.Language)
}
