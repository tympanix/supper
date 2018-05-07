package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
