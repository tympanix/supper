package cfg

import (
	"testing"
	"time"

	"github.com/tympanix/supper/score"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

func TestConfigDefaults(t *testing.T) {
	Initialize()

	assert.Equal(t, Default.Evaluator(), &score.DefaultEvaluator{})
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
