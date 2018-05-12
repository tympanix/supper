package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/misc"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
)

func TestForeignMovie(t *testing.T) {
	m, err := NewMovie("Den_utrolige_historie_om_den_kæmpestore_pære_2017")
	require.NoError(t, err)
	assert.Equal(t, "Den utrolige historie om den kæmpestore pære", m.MovieName())
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
}
