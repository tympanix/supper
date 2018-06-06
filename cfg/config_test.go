package cfg

import (
	"bytes"
	"testing"
	"time"

	"github.com/tympanix/supper/provider"
	"github.com/tympanix/supper/score"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestConfigFromViper(t *testing.T) {
	viper.Set("force", true)
	viper.Set("verbose", true)
	viper.Set("dry", true)
	viper.Set("strict", true)
	viper.Set("impaired", true)

	viper.Set("modified", "2h30m")
	viper.Set("config", "/foo/bar")
	viper.Set("delay", "32s")
	viper.Set("logfile", "/foo/bar/baz/log")
	viper.Set("action", "move")

	viper.Set("limit", 64)
	viper.Set("score", 89)

	viper.Set("lang", []string{
		"en",
		"de",
		"es",
	})

	Initialize()

	assert.Equal(t, Default.Force(), true)
	assert.Equal(t, Default.Verbose(), true)
	assert.Equal(t, Default.Dry(), true)
	assert.Equal(t, Default.Strict(), true)
	assert.Equal(t, Default.Impaired(), true)

	assert.Equal(t, Default.Modified(), 2*time.Hour+30*time.Minute)
	assert.Equal(t, Default.Config(), "/foo/bar")
	assert.Equal(t, Default.Delay(), 32*time.Second)
	assert.Equal(t, Default.Logfile(), "/foo/bar/baz/log")
	assert.Equal(t, Default.RenameAction(), "move")

	assert.Equal(t, Default.Limit(), 64)
	assert.Equal(t, Default.Score(), 89)

	assert.Equal(t, Default.Languages().Size(), 3)
	assert.True(t, Default.Languages().Has(language.English))
	assert.True(t, Default.Languages().Has(language.German))
	assert.True(t, Default.Languages().Has(language.Spanish))
}

func TestConfigThirdParty(t *testing.T) {
	viper.Set("apikeys", map[string]string{
		"themoviedb": "tmdb_test_key",
		"thetvdb":    "tvdb_test_key",
	})

	Initialize()

	assert.Equal(t, "tmdb_test_key", Default.APIKeys().TheMovieDB())
	assert.Equal(t, "tvdb_test_key", Default.APIKeys().TheTVDB())

	assert.Equal(t, Default.Evaluator(), &score.DefaultEvaluator{})
	assert.Contains(t, Default.Providers(), provider.Subscene())
	assert.Contains(t, Default.Scrapers(), provider.TheMovieDB("tmdb_test_key"))
	assert.Contains(t, Default.Scrapers(), provider.TheTVDB("tvdb_test_key"))
}

func TestConfigPlugins(t *testing.T) {
	plugins := []map[string]string{
		{
			"name": "test_name",
			"exec": "test_exec",
		},
	}

	viper.Set("plugins", plugins)

	Initialize()

	assert.Equal(t, len(Default.Plugins()), 1)

	for i, p := range plugins {
		golden := Default.Plugins()[i]
		assert.Equal(t, golden.Name(), p["name"])
	}
}

func TestConfigMedia(t *testing.T) {

	viper.Set("movies", map[string]interface{}{
		"directory": "/foo/bar/movie",
		"template":  "test_template_movie",
	})

	viper.Set("tvshows", map[string]interface{}{
		"directory": "/foo/bar/tvshow",
		"template":  "test_template_tvshow",
	})

	Initialize()

	// test movies
	assert.Equal(t, "/foo/bar/movie", Default.Movies().Directory())
	var mbuf bytes.Buffer
	err := Default.Movies().Template().Execute(&mbuf, nil)
	require.NoError(t, err)
	assert.Equal(t, "test_template_movie", mbuf.String())

	// test tvshows
	assert.Equal(t, "/foo/bar/tvshow", Default.TVShows().Directory())
	var tbuf bytes.Buffer
	err = Default.TVShows().Template().Execute(&tbuf, nil)
	require.NoError(t, err)
	assert.Equal(t, "test_template_tvshow", tbuf.String())
}
