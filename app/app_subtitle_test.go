package app

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/media"

	"github.com/fatih/set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

type subtitle struct {
	types.Media
	lang language.Tag
	hi   bool
}

func (s subtitle) String() string                       { return "Subtitle: " + s.Media.String() }
func (s subtitle) ForMedia() types.Media                { return s.Media }
func (s subtitle) HearingImpaired() bool                { return s.hi }
func (s subtitle) Language() language.Tag               { return s.lang }
func (s subtitle) TypeSubtitle() (types.Subtitle, bool) { return s, true }
func (s subtitle) TypeEpisode() (types.Episode, bool)   { return nil, false }
func (s subtitle) TypeMovie() (types.Movie, bool)       { return nil, false }

type online struct {
	types.Subtitle
	data []byte
}

func (o online) Download() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer(o.data)), nil
}

func (o online) Link() string {
	return ""
}

type subtitleTester interface {
	Pre(*testing.T)
	Input() string
	Test(*testing.T, types.LocalSubtitle)
	Post(*testing.T, []types.LocalSubtitle)
}

type fakeplugin func(types.LocalSubtitle) error

func (p fakeplugin) Run(s types.LocalSubtitle) error {
	return p(s)
}

func (p fakeplugin) Name() string {
	return "fakeplugin"
}

func must(m types.LocalMedia, err error) types.LocalMedia {
	if err != nil {
		panic(err)
	}
	return m
}

var inception = must(media.NewLocalFile("test/Inception.2010.720p.x264.mkv"))
var gameofthrones = must(media.NewLocalFile("test/Game.of.Thrones.s01e02.mp4"))

var subtitles = []types.OnlineSubtitle{
	online{subtitle{inception, language.German, false}, []byte("online_inception")},
	online{subtitle{gameofthrones, language.German, false}, []byte("online_gameofthrones")},
}

func TestDownloadSubtitles(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.languages = set.New(language.German)
	config.providers = []types.Provider{
		fakeProvider{
			subtitles,
		},
	}

	err := performSubtitleTest(t, successTester{}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestSubtitlePlugins(t *testing.T) {
	config := defaultConfig

	var results []types.LocalSubtitle

	config.plugins = []types.Plugin{
		fakeplugin(func(s types.LocalSubtitle) error {
			results = append(results, s)
			return nil
		}),
	}

	config.providers = []types.Provider{
		fakeProvider{
			subtitles,
		},
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, pluginTester{&results}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestSubtitlePluginError(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.plugins = []types.Plugin{
		fakeplugin(func(s types.LocalSubtitle) error {
			return errors.New("test plugin error")
		}),
	}

	config.providers = []types.Provider{
		fakeProvider{
			subtitles,
		},
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, successTester{}, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test plugin error")

	cleanRenameTest(t)
}

func TestSubtitleNoMedia(t *testing.T) {
	app := New(defaultConfig)

	_, err := app.DownloadSubtitles(nil, set.New(language.English))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no media")
}

func TestSubtitleNoLanguage(t *testing.T) {
	app := New(defaultConfig)

	_, err := app.DownloadSubtitles(list.NewLocalMedia(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no languages")
}

func TestSubtitleNoVideo(t *testing.T) {
	app := New(defaultConfig)

	_, err := app.DownloadSubtitles(list.NewLocalMedia(), set.New(language.English))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no video")
}

func copyTestFiles(src, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		i, err := os.Open(filepath.Join(src, f.Name()))
		if err != nil {
			return err
		}
		defer i.Close()
		o, err := os.Create(filepath.Join(dst, f.Name()))
		if err != nil {
			return err
		}
		defer o.Close()
		_, err = io.Copy(o, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func performSubtitleTest(t *testing.T, test subtitleTester, config types.Config) error {

	app := New(config)

	require.NoError(t, copyTestFiles(test.Input(), "out"))

	media, err := app.FindMedia("out")
	require.NoError(t, err)
	assert.Equal(t, len(res), media.Len())

	test.Pre(t)

	list, err := app.DownloadSubtitles(media, config.Languages())
	if err != nil {
		return err
	}

	for _, s := range list {
		test.Test(t, s)
	}

	test.Post(t, list)

	return nil
}

type successTester struct{}

func (successTester) Pre(t *testing.T) {
}

func (successTester) Input() string {
	return "test"
}

func (successTester) Test(t *testing.T, s types.LocalSubtitle) {

}

func (successTester) Post(t *testing.T, l []types.LocalSubtitle) {
	assert.Equal(t, len(subtitles), len(l))
}

type pluginTester struct {
	runs *[]types.LocalSubtitle
}

func (p pluginTester) Pre(t *testing.T) {
}

func (p pluginTester) Input() string {
	return "test"
}

func (p pluginTester) Test(t *testing.T, s types.LocalSubtitle) {
	assert.Contains(t, *p.runs, s)
}

func (p pluginTester) Post(t *testing.T, l []types.LocalSubtitle) {
	assert.Equal(t, len(subtitles), len(l))
}
