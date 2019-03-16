package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tympanix/supper/media/list"
	"github.com/tympanix/supper/app/notify"

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

type onlineError struct {
	types.Subtitle
}

func (o onlineError) Download() (io.ReadCloser, error) {
	return nil, errors.New("test download subtitle error")
}

func (o onlineError) Link() string {
	return ""
}

type mockSaveSubtitleError struct {
	types.Video
}

func (mockSaveSubtitleError) SaveSubtitle(io.Reader, language.Tag) (types.LocalSubtitle, error) {
	return nil, errors.New("test save subtitle")
}

type fakeEvaluator func(types.Media, types.Media) float32

func (e fakeEvaluator) Evaluate(m types.Media, n types.Media) float32 {
	return e(m, n)
}

type subtitleTester interface {
	Pre(*testing.T, []types.LocalMedia)
	Input() string
	Mock(types.LocalMedia) types.LocalMedia
	Test(*testing.T, types.LocalSubtitle)
	Post(*testing.T, []types.Video, []types.LocalSubtitle)
}

type fakePlugin func(types.LocalSubtitle) error

func (p fakePlugin) Run(s types.LocalSubtitle) error {
	return p(s)
}

func (p fakePlugin) Name() string {
	return "fakeplugin"
}

type fakeProviderError struct{}

func (p fakeProviderError) SearchSubtitles(m types.LocalMedia) ([]types.OnlineSubtitle, error) {
	return nil, errors.New("test provider error")
}

func (p fakeProviderError) ResolveSubtitle(l types.Linker) (types.Downloadable, error) {
	return nil, errors.New("test provider does not support resolving subtitles")
}

func must(m types.LocalMedia, err error) types.LocalMedia {
	if err != nil {
		panic(err)
	}
	return m
}

type fakeTimedProvider struct {
	delay time.Duration
	fakeProvider
	last time.Time
}

func (p *fakeTimedProvider) SearchSubtitles(m types.LocalMedia) ([]types.OnlineSubtitle, error) {
	fmt.Println(time.Since(p.last))
	if time.Since(p.last) < p.delay {
		return nil, errors.New("expected delay to occur")
	}
	p.last = time.Now()
	fmt.Println(p.last)
	return p.fakeProvider.SearchSubtitles(m)
}

func (p *fakeTimedProvider) reset() {
	p.last = time.Unix(0, 0)
}

func TestDownloadSubtitles(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestSubtitlePlugins(t *testing.T) {
	config := defaultConfig

	var results []types.LocalSubtitle

	config.plugins = []types.Plugin{
		fakePlugin(func(s types.LocalSubtitle) error {
			results = append(results, s)
			return nil
		}),
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
		fakePlugin(func(s types.LocalSubtitle) error {
			return errors.New("test plugin error")
		}),
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test plugin error")

	cleanRenameTest(t)
}

func TestSubtitleNoMedia(t *testing.T) {
	app := New(defaultConfig)

	c := notify.AsyncDiscard()
	defer close(c)

	_, err := app.DownloadSubtitles(nil, set.New(language.English), c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no media")
}

func TestSubtitleNoLanguage(t *testing.T) {
	app := New(defaultConfig)

	c := notify.AsyncDiscard()
	defer close(c)

	_, err := app.DownloadSubtitles(list.NewLocalMedia(), nil, c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no languages")
}

func TestSubtitleNoVideo(t *testing.T) {
	app := New(defaultConfig)

	c := notify.AsyncDiscard()
	defer close(c)

	_, err := app.DownloadSubtitles(list.NewLocalMedia(), set.New(language.English), c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no video")
}

func TestSubtitleScore(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.evaluator = fakeEvaluator(func(m types.Media, n types.Media) float32 {
		return 0.01
	})

	config.languages = set.New(language.German)

	config.score = 100

	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Score too low")
	assert.Contains(t, err.Error(), "1%")

	cleanRenameTest(t)
}

func TestSubtitleProviderError(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.providers = []types.Provider{
		fakeProviderError{},
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test provider error")

	cleanRenameTest(t)
}

func TestSubtitleNoneAvailable(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.languages = set.New(language.Russian)

	err := performSubtitleTest(t, skipSubtitlesTest{}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestSubtitleDryRun(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.dry = true

	// when dry run, provider should never be called
	config.providers = []types.Provider{
		fakeProviderError{},
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, skipSubtitlesTest{}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestSubtitleUnsatisfied(t *testing.T) {
	config := defaultConfig
	config.strict = true

	// media is always unsatisfied (score = 0.0)
	config.evaluator = fakeEvaluator(func(m types.Media, n types.Media) float32 {
		return 0.0
	})

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, skipSubtitlesTest{}, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No subtitles satisfied media")

	cleanRenameTest(t)
}

func TestSubtitleDownloadError(t *testing.T) {
	config := defaultConfig
	config.strict = true

	config.providers = []types.Provider{
		fakeProviderDownloadError{[]language.Tag{
			language.English,
			language.German,
			language.Spanish,
		}},
	}

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test download subtitle error")

	cleanRenameTest(t)
}

func TestSubtitleSaveError(t *testing.T) {
	defer cleanRenameTest(t)

	config := defaultConfig
	config.strict = true

	config.languages = set.New(language.German)

	err := performSubtitleTest(t, saveErrorSubtitleTest{}, config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "test save subtitle")
}

func TestSubtitleDelay(t *testing.T) {
	defer cleanRenameTest(t)

	const testDelay = 250 * time.Millisecond

	config := defaultConfig

	config.strict = true
	config.languages = set.New(language.German)

	timedProvider := &fakeTimedProvider{
		testDelay,
		fakeProvider{[]language.Tag{
			language.German,
		}},
		time.Unix(0, 0),
	}

	config.providers = []types.Provider{
		timedProvider,
	}

	// first run, without delay, should return error
	err := performSubtitleTest(t, subtitleLangTester(language.German), config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected delay to occur")

	cleanRenameTest(t)

	// second run, with delay, should not return error
	timedProvider.reset()
	config.delay = time.Duration(testDelay)
	err = performSubtitleTest(t, subtitleLangTester(language.German), config)
	assert.NoError(t, err)
}

func TestSubtitleInvalidLanguages(t *testing.T) {
	defer cleanRenameTest(t)

	config := defaultConfig
	config.strict = true

	config.languages = set.New(42)

	err := performSubtitleTest(t, skipSubtitlesTest{}, config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown language")
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

	mocked := media.List()
	for i, m := range mocked {
		mocked[i] = test.Mock(m)
	}
	media = list.NewLocalMedia(mocked...)

	test.Pre(t, media.List())

	c := notify.AsyncDiscard()
	defer close(c)

	subs, err := app.DownloadSubtitles(media, config.Languages(), c)
	if err != nil {
		return err
	}

	for _, s := range subs {
		test.Test(t, s)
	}

	test.Post(t, media.FilterVideo().List(), subs)

	return nil
}

type subtitleLangTester language.Tag

func (subtitleLangTester) Pre(t *testing.T, l []types.LocalMedia) {
	assert.Equal(t, len(res), len(l))
}

func (subtitleLangTester) Mock(m types.LocalMedia) types.LocalMedia {
	return m
}

func (subtitleLangTester) Input() string {
	return "test"
}

func (l subtitleLangTester) Test(t *testing.T, s types.LocalSubtitle) {
	assert.Equal(t, s.Language(), language.Tag(l))
}

func (subtitleLangTester) Post(t *testing.T, m []types.Video, l []types.LocalSubtitle) {
	assert.Equal(t, len(m), len(l))
}

type pluginTester struct {
	runs *[]types.LocalSubtitle
}

func (p pluginTester) Pre(t *testing.T, l []types.LocalMedia) {
	assert.Equal(t, len(res), len(l))
}

func (p pluginTester) Input() string {
	return "test"
}

func (p pluginTester) Mock(m types.LocalMedia) types.LocalMedia {
	return m
}

func (p pluginTester) Test(t *testing.T, s types.LocalSubtitle) {
	assert.Contains(t, *p.runs, s)
}

func (p pluginTester) Post(t *testing.T, m []types.Video, l []types.LocalSubtitle) {
	//assert.Equal(t, len(m), len(l))
}

type skipSubtitlesTest struct{}

func (p skipSubtitlesTest) Pre(t *testing.T, l []types.LocalMedia) {
}

func (p skipSubtitlesTest) Input() string {
	return "test"
}

func (p skipSubtitlesTest) Mock(m types.LocalMedia) types.LocalMedia {
	return m
}

func (p skipSubtitlesTest) Test(t *testing.T, s types.LocalSubtitle) {
	assert.Fail(t, "should skip all subtitles")
}

func (p skipSubtitlesTest) Post(t *testing.T, m []types.Video, l []types.LocalSubtitle) {
	assert.Equal(t, 0, len(l))
}

type saveErrorSubtitleTest struct {
	skipSubtitlesTest
}

func (saveErrorSubtitleTest) Mock(m types.LocalMedia) types.LocalMedia {
	if v, ok := m.(types.Video); ok {
		return mockSaveSubtitleError{v}
	}
	return m
}
