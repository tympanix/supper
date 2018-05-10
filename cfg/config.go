package cfg

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
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

type mediaConfig struct {
	DirectoryX string `mapstructure:"directory"`
	TemplateX  string `mapstructure:"template"`
}

func (m mediaConfig) Interface() *Media {
	template := template.New(m.DirectoryX).Funcs(renameFuncs)
	template, err := template.Parse(cleanTemplate(m.TemplateX))
	if err != nil {
		log.WithError(err).Fatal("could not parse renaming template")
	}
	return &Media{
		directory: m.DirectoryX,
		template:  template,
	}
}

// Media is a configuration object for media collections
type Media struct {
	directory string
	template  *template.Template
}

// Directory returns the directory whe the media is located
func (m *Media) Directory() string {
	return m.directory
}

// Template returns the template for renaming the media
func (m *Media) Template() *template.Template {
	return m.template
}

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

var renameFuncs = template.FuncMap{
	"pad": func(d int) string {
		return fmt.Sprintf("%02d", d)
	},
}

var templateRegex = regexp.MustCompile(`[\r\n]`)
var seperatorRegex = regexp.MustCompile(`\s*/\s*`)

func cleanTemplate(template string) string {
	cleaned := templateRegex.ReplaceAllString(template, "")
	cleaned = seperatorRegex.ReplaceAllString(cleaned, "/")
	return strings.TrimSpace(cleaned)
}

type viperConfig struct {
	languages set.Interface
	modified  time.Duration
	delay     time.Duration
	plugins   []types.Plugin
	apikeys   map[string]string
	movies    types.MediaConfig
	tvshows   types.MediaConfig
	filters   int
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
		if p.PluginName == "" || p.Exec == "" {
			log.Fatal("Invalid plugin definitions, missing name and/or exec")
		}
		plugins = append(plugins, &p)
	}

	media := map[string]*Media{
		"movies":  &Media{},
		"tvshows": &Media{},
	}

	for k, v := range media {
		sub := viper.Sub(k)
		if sub == nil {
			continue
		}
		var media mediaConfig
		if err := sub.Unmarshal(&media); err != nil {
			log.WithError(err).Fatalf("invalid configuration for %v", k)
		}
		*v = *media.Interface()
	}

	var filters int
	for _, b := range []bool{
		viper.GetBool("filter-movies"),
		viper.GetBool("filter-tvshows"),
		viper.GetBool("filter-subtitles"),
	} {
		if b {
			filters++
		}
	}

	Default = viperConfig{
		languages: lang,
		modified:  modified,
		delay:     delay,
		plugins:   plugins,
		apikeys:   viper.GetStringMapString("apikeys"),
		movies:    media["movies"],
		tvshows:   media["tvshows"],
		filters:   filters,
	}
}

func (v viperConfig) MediaFilter() types.MediaFilter {
	if v.filters == 0 {
		return nil
	}

	return func(m types.Media) bool {
		if _, ok := m.TypeMovie(); ok && viper.GetBool("filter-movies") {
			return true
		}

		if _, ok := m.TypeEpisode(); ok && viper.GetBool("filter-tvshows") {
			return true
		}

		if _, ok := m.TypeSubtitle(); ok && viper.GetBool("filter-subtitles") {
			return true
		}

		return false
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

func (v viperConfig) Movies() types.MediaConfig {
	return v.movies
}

func (v viperConfig) TVShows() types.MediaConfig {
	return v.tvshows
}
