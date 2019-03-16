package media

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/media/meta/codec"
	"github.com/tympanix/supper/media/meta/misc"
	"github.com/tympanix/supper/media/meta/quality"
	"github.com/tympanix/supper/media/meta/source"
)

func TestMovie(t *testing.T) {
	m, err := NewMovie("Inception.2010.BluRay.1080p.x264.3D-GROUP")
	require.NoError(t, err)
	assert.Equal(t, "Inception", m.MovieName())
	assert.Equal(t, 2010, m.Year())
	assert.Equal(t, quality.HD1080p, m.Quality())
	assert.Equal(t, source.BluRay, m.Source())
	assert.Equal(t, "GROUP", m.Group())
	assert.Equal(t, codec.X264, m.Codec())
	assert.True(t, m.Misc().Has(misc.Video3D))

	assert.Contains(t, m.String(), "Inception")
	assert.Contains(t, m.String(), "2010")

	assert.Contains(t, m.Identity(), "inception")
	assert.Contains(t, m.Identity(), "2010")
}

func TestMovieError(t *testing.T) {
	m, err := NewMovie("blablatesttest")
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestMovieMerge(t *testing.T) {
	m1, err := NewMovie("iron.man.2008.BluRay.x264")
	require.NoError(t, err)

	m2 := &Movie{
		NameX: "Iron Man",
		YearX: 2008,
	}

	err = m1.Merge(m2)
	require.NoError(t, err)

	assert.Equal(t, "Iron Man", m1.MovieName())
	assert.Equal(t, 2008, m1.Year())
}

func TestMovieMergeError(t *testing.T) {
	m, err := NewMovie("iron.man.2008")
	require.NoError(t, err)
	e, err := NewEpisode("the.office.us.s01e03")
	require.NoError(t, err)

	err = m.Merge(e)
	assert.Error(t, err)
}

func TestMovieMergeOneYearDiff(t *testing.T) {
	m1, err := NewMovie("iron.man.2008")
	require.NoError(t, err)
	m2, err := NewMovie("iron.man.2009")
	require.NoError(t, err)

	err = m1.Merge(m2)
	assert.NoError(t, err)
}

func TestMovieMergeTwoYearDiff(t *testing.T) {
	m1, err := NewMovie("iron.man.2008")
	require.NoError(t, err)
	m2, err := NewMovie("iron.man.2010")
	require.NoError(t, err)

	err = m1.Merge(m2)
	assert.Error(t, err)
}

func TestForeignMovie(t *testing.T) {
	m, err := NewMovie("Den_utrolige_historie_om_den_kæmpestore_pære_2017")
	require.NoError(t, err)
	assert.Equal(t, "Den Utrolige Historie Om Den Kæmpestore Pære", m.MovieName())
	assert.Equal(t, 2017, m.Year())

	m, err = NewMovie("Les.Misérables.2012")
	require.NoError(t, err)
	assert.Equal(t, "Les Misérables", m.MovieName())
	assert.Equal(t, 2012, m.Year())

	m, err = NewMovie("Die.Fälscher.2007.1080p")
	require.NoError(t, err)
	assert.Equal(t, "Die Fälscher", m.MovieName())
	assert.Equal(t, 2007, m.Year())

	m, err = NewMovie("海底总动员_2003_720p")
	require.NoError(t, err)
	assert.Equal(t, "海底总动员", m.MovieName())
	assert.Equal(t, 2003, m.Year())
}

func TestMovieWebsite(t *testing.T) {
	m, err := NewMovie("[www.example.com] Inception (2010) x264 1080p")
	require.NoError(t, err)
	assert.Equal(t, "Inception", m.MovieName())
	assert.Equal(t, 2010, m.Year())
}

func TestMovieJSON(t *testing.T) {
	m, err := NewMovie("Inception.2010.720p")
	require.NoError(t, err)

	data, err := json.Marshal(m)
	require.NoError(t, err)

	j := struct {
		Name string `json:"name"`
		Year int    `json:"year"`
	}{}

	err = json.Unmarshal(data, &j)
	require.NoError(t, err)

	assert.Equal(t, 2010, m.Year())
	assert.Equal(t, "Inception", m.MovieName())
}
