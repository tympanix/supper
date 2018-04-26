package cfg

import (
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/plugin"
	"github.com/tympanix/supper/types"

	"github.com/apex/log"
	"github.com/fatih/set"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

// Default holds the global application configuration instance
var Default types.Config

var homePath string

func init() {
	home, err := homedir.Dir()

	if err != nil {
		log.WithError(err).Fatal("Could not find user home directory")
	}

	homePath = home
}

type viperConfig struct {
	languages set.Interface
	modified  time.Duration
	delay     time.Duration
	plugins   []types.Plugin
	apikeys   map[string]string
}

// Initialize construct the default configuration object using viper.
// Ths function must be called once all CLI flags and configuration files
// has been parsed.
func Initialize() {
	// Parse all language flags into slice of tags
	lang := set.New()
	for _, tag := range viper.GetStringSlice("lang") {
		_lang, err := language.Parse(tag)
		if err != nil {
			log.WithField("language", tag).Fatal("Invalid language tag")
		}
		lang.Add(_lang)
	}

	// Parse modified flag
	modified, err := parse.Duration(viper.GetString("modified"))
	if err != nil {
		log.WithError(err).WithField("modified", viper.GetString("modified")).
			Fatal("Invalid duration")
	}

	// Parse delay flag
	delay, err := parse.Duration(viper.GetString("delay"))
	if err != nil {
		log.WithError(err).WithField("delay", viper.GetString("modified")).
			Fatal("Invalid duration")
	}

	// Parse plugins
	var _plugins []plugin.Plugin
	if err := viper.UnmarshalKey("plugins", &_plugins); err != nil {
		log.WithError(err).Fatal("Invalid plugin definition")
	}

	plugins := make([]types.Plugin, 0)
	for _, p := range _plugins {
		plugins = append(plugins, &p)
	}

	Default = viperConfig{
		languages: lang,
		modified:  modified,
		delay:     delay,
		plugins:   plugins,
		apikeys:   viper.GetStringMapString("apikeys"),
	}
}

func (v viperConfig) Languages() set.Interface {
	return v.languages
}

func (v viperConfig) Verbose() bool {
	return viper.GetBool("verbose")
}

func (v viperConfig) Dry() bool {
	return viper.GetBool("dry")
}

func (v viperConfig) Strict() bool {
	return viper.GetBool("strict")
}

func (v viperConfig) Modified() time.Duration {
	return v.modified
}

func (v viperConfig) Config() string {
	return viper.GetString("config")
}

func (v viperConfig) Delay() time.Duration {
	return v.delay
}

func (v viperConfig) Force() bool {
	return viper.GetBool("force")
}

func (v viperConfig) Impaired() bool {
	return viper.GetBool("impaired")
}

func (v viperConfig) Limit() int {
	return viper.GetInt("limit")
}

func (v viperConfig) Logfile() string {
	return viper.GetString("logfile")
}

func (v viperConfig) Score() int {
	return viper.GetInt("score")
}

func (v viperConfig) Plugins() []types.Plugin {
	return v.plugins
}

func (v viperConfig) APIKeys() types.APIKeys {
	return v
}

func (v viperConfig) TheTVDB() string {
	return v.apikeys["thetvdb"]
}

func (v viperConfig) TheMovieDB() string {
	return v.apikeys["themoviedb"]
}
