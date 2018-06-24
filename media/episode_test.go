package media

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/meta/quality"
)

func TestEpisode(t *testing.T) {
	e, err := NewEpisode("The.Office.US.S02E04.720p")
	_, ok := e.TypeEpisode()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.True(t, e.IsVideo())
	assert.Equal(t, quality.HD720p, e.Quality())
	assert.Equal(t, 2, e.Season())
	assert.Equal(t, 4, e.Episode())
	assert.Equal(t, "The Office US", e.TVShow())
	assert.Contains(t, e.String(), "S02E04")
	assert.Equal(t, "theofficeus:2:4", e.Identity())
}

func TestEpisodeNoGroup(t *testing.T) {
	e, err := NewEpisode("Friends.S01E01.The.One.Where.Monica.Gets.a.Roommate")
	require.NoError(t, err)
	assert.Equal(t, 1, e.Season())
	assert.Equal(t, 1, e.Episode())
	assert.Equal(t, "", e.Group())
	assert.Equal(t, "The One Where Monica Gets a Roommate", e.EpisodeName())
}

func TestEpisodeGroup(t *testing.T) {
	e, err := NewEpisode("Friends.S01E01.The.One.Where.Monica.Gets.a.Roommate.720p.GROUP")
	require.NoError(t, err)
	assert.Equal(t, 1, e.Season())
	assert.Equal(t, 1, e.Episode())
	assert.Equal(t, "GROUP", e.Group())
	assert.Equal(t, quality.HD720p, e.Quality())
	assert.Equal(t, "The One Where Monica Gets a Roommate", e.EpisodeName())
}

func TestEpisodeError(t *testing.T) {
	_, err := NewEpisode("blablatestest")
	assert.Error(t, err)
}

func TestEpisodeJSON(t *testing.T) {
	e, err := NewEpisode("Silicon.Valley.S02E05.720p.x264")
	require.NoError(t, err)

	j := struct {
		Name    string `json:"name"`
		Episode int    `json:"episode"`
		Season  int    `json:"season"`
	}{}

	data, err := json.Marshal(e)
	require.NoError(t, err)

	err = json.Unmarshal(data, &j)
	require.NoError(t, err)

	assert.Equal(t, "Silicon Valley", j.Name)
	assert.Equal(t, 2, j.Season)
	assert.Equal(t, 5, j.Episode)
}

func TestEpisodeMerge(t *testing.T) {
	e, err := NewEpisode("game of thrones 1x5 720p")
	c := &Episode{
		NameX:        "Game of Thrones",
		SeasonX:      1,
		EpisodeX:     5,
		EpisodeNameX: "The Wolf and the Lion",
	}

	require.NoError(t, err)

	err = e.Merge(c)
	require.NoError(t, err)
	assert.Equal(t, "Game of Thrones", e.TVShow())
	assert.Equal(t, quality.HD720p, e.Quality())
	assert.Equal(t, "The Wolf and the Lion", e.EpisodeName())
}

func TestEpisodeMergeEpisodeError(t *testing.T) {
	e1, err := NewEpisode("Game of Thrones 1x1")
	require.NoError(t, err)
	e2, err := NewEpisode("Game of Thrones 1x2")
	require.NoError(t, err)

	err = e1.Merge(e2)
	assert.Error(t, err)
}

func TestEpisodeMergeSeasonError(t *testing.T) {
	e1, err := NewEpisode("Arrow 1x1")
	require.NoError(t, err)
	e2, err := NewEpisode("Arrow 2x1")
	require.NoError(t, err)

	err = e1.Merge(e2)
	assert.Error(t, err)
}

func TestEpisodeMergeMediaError(t *testing.T) {
	e, err := NewEpisode("Community 1x1")
	require.NoError(t, err)
	m, err := NewMovie("Inception 2010")
	require.NoError(t, err)

	err = e.Merge(m)
	assert.Error(t, err)
}
