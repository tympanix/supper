package app

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/provider"

	"github.com/fatih/set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
)

const (
	defaultMovieStr  = "{{ .Movie }} ({{ .Year }}) {{ .Quality }}"
	defaultTVShowStr = "{{ .TVShow }} S{{ .Season }}E{{ .Episode }}"
)

var (
	fakeDefaultMovieTemplate   = template.Must(template.New("movie").Parse(defaultMovieStr))
	fakeDefaultTVShowsTemplate = template.Must(template.New("tv").Parse(defaultTVShowStr))
)

var defaultConfig = struct {
	fakeConfig
	fakeTemplates
}{
	fakeConfig{
		action:    "copy",
		scrapers:  []types.Scraper{fakeScraper{}},
		providers: []types.Provider{fakeProvider{}},
	},
	fakeTemplates{
		output:         "out",
		movieTemplate:  fakeDefaultMovieTemplate,
		tvshowTemplate: fakeDefaultTVShowsTemplate,
	},
}

type fakeProvider struct {
	subs []types.OnlineSubtitle
}

func (f fakeProvider) ResolveSubtitle(link types.Linker) (types.Downloadable, error) {
	for _, l := range f.subs {
		if l.Link() == link.Link() {
			return l, nil
		}
	}
	return nil, errors.New("mocked subtitle not found")
}

func (f fakeProvider) SearchSubtitles(m types.LocalMedia) ([]types.OnlineSubtitle, error) {
	return f.subs, nil
}

type fakeConfig struct {
	action    string
	strict    bool
	force     bool
	dry       bool
	scrapers  []types.Scraper
	providers []types.Provider
	languages set.Interface
}

func (c fakeConfig) Languages() set.Interface       { return c.languages }
func (c fakeConfig) APIKeys() types.APIKeys         { return fakeAPIKeys{} }
func (c fakeConfig) Config() string                 { return "" }
func (c fakeConfig) Delay() time.Duration           { return 0 }
func (c fakeConfig) Dry() bool                      { return c.dry }
func (c fakeConfig) Force() bool                    { return c.force }
func (c fakeConfig) Impaired() bool                 { return false }
func (c fakeConfig) Limit() int                     { return -1 }
func (c fakeConfig) Logfile() string                { return "" }
func (c fakeConfig) MediaFilter() types.MediaFilter { return nil }
func (c fakeConfig) Modified() time.Duration        { return 0 }
func (c fakeConfig) Plugins() []types.Plugin        { return nil }
func (c fakeConfig) Score() int                     { return 0 }
func (c fakeConfig) Strict() bool                   { return c.strict }
func (c fakeConfig) Verbose() bool                  { return false }
func (c fakeConfig) Providers() []types.Provider    { return c.providers }
func (c fakeConfig) Scrapers() []types.Scraper      { return c.scrapers }
func (c fakeConfig) RenameAction() string           { return c.action }

type fakeTemplates struct {
	output         string
	movieTemplate  *template.Template
	tvshowTemplate *template.Template
}

func (c fakeTemplates) Movies() types.MediaConfig {
	return fakeMediaConfig{
		directory: c.output,
		template:  c.movieTemplate,
	}
}

func (c fakeTemplates) TVShows() types.MediaConfig {
	return fakeMediaConfig{
		directory: c.output,
		template:  c.tvshowTemplate,
	}
}

type fakeAPIKeys struct{}

func (k fakeAPIKeys) TheMovieDB() string { return "" }
func (k fakeAPIKeys) TheTVDB() string    { return "" }

type fakeMediaConfig struct {
	directory string
	template  *template.Template
}

func (m fakeMediaConfig) Directory() string            { return m.directory }
func (m fakeMediaConfig) Template() *template.Template { return m.template }

type fakeScraper []types.Media

func (s fakeScraper) Scrape(m types.Media) (types.Media, error) {
	if s, ok := m.TypeSubtitle(); ok {
		return s.ForMedia(), nil
	}
	return m, nil
}

var res = map[string]string{
	"Inception (2010) 720p.mkv":    "Inception.2010.720p.x264.mkv",
	"Inception (2010) 720p.en.srt": "Inception.2010.720p.x264.en.srt",
	"Game of Thrones S1E2.mp4":     "Game.of.Thrones.s01e02.mp4",
	"Game of Thrones S1E2.en.srt":  "Game.of.Thrones.s01e02.en.srt",
}

type renameTester interface {
	Pre(*testing.T)
	Input() string
	Output() string
	Test(*testing.T, string, string)
	Post(*testing.T, []os.FileInfo)
}

func TestRenameMedia(t *testing.T) {
	testCases := map[string]renameTester{
		"copy":     copyTester{},
		"move":     moveTester{},
		"hardlink": copyTester{},
		"symlink":  symlinkTester{},
	}

	for action, test := range testCases {
		t.Run(action, func(t *testing.T) {
			config := defaultConfig

			config.strict = true
			config.action = action
			config.output = test.Output()

			assert.NoError(t, performRenameTest(t, test, config))

			cleanRenameTest(t)
		})
	}
}

func TestRenameMediaOverride(t *testing.T) {
	var err error
	config := defaultConfig

	// first run, rename media sucessfully
	err = performRenameTest(t, copyTester{}, config)
	require.NoError(t, err)

	// second run, skip all media (strict=false)
	err = performRenameTest(t, copyTester{}, config)
	require.NoError(t, err, "should skip all conflicts (non-strict mode)")

	config.strict = true

	// third run, return error (strict=true)
	err = performRenameTest(t, copyTester{}, config)
	require.Error(t, err, "should error on conflict (strict mode)")
	assert.True(t, media.IsExistsErr(err))

	config.force = true

	// forth run, overwrite with success (force=true)
	err = performRenameTest(t, copyTester{}, config)
	require.NoError(t, err, "should overwrite media (force mode)")

	cleanRenameTest(t)
}

func TestRenameActionError(t *testing.T) {
	config := defaultConfig
	config.action = "invalid_test_action"

	err := performRenameTest(t, copyTester{}, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_test_action")
}

func TestRenameDryRun(t *testing.T) {
	config := defaultConfig
	config.dry = true

	assert.NoError(t, os.Mkdir("out", os.ModePerm))

	err := performRenameTest(t, dryTester{}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

func TestRenameTemplateEmpty(t *testing.T) {
	config := defaultConfig

	for _, sp := range []**template.Template{
		&config.movieTemplate,
		&config.tvshowTemplate,
	} {
		tmp := *sp
		*sp = template.Must(template.New("test").Parse(""))

		err := performRenameTest(t, copyTester{}, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
		assert.Contains(t, err.Error(), "template")

		cleanRenameTest(t)

		*sp = tmp
	}
}

func TestRenameTemplateError(t *testing.T) {
	config := defaultConfig

	for sp, templ := range map[**template.Template]string{
		&config.movieTemplate:  "{{ .Movie }} {{ .InvalidTestField }}",
		&config.tvshowTemplate: "{{ .TVShow }} {{ .InvalidTestField }}",
	} {
		tmp := *sp
		*sp = template.Must(template.New("test").Parse(templ))

		err := performRenameTest(t, copyTester{}, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "InvalidTestField")

		cleanRenameTest(t)

		*sp = tmp
	}
}

func TestRenameTemplateNil(t *testing.T) {
	config := defaultConfig

	for _, tp := range []**template.Template{
		&config.movieTemplate,
		&config.tvshowTemplate,
	} {
		tmp := *tp
		*tp = nil

		err := performRenameTest(t, copyTester{}, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing template")

		cleanRenameTest(t)

		*tp = tmp
	}
}

type fakeUnsupportedScraper struct{}

func (fakeUnsupportedScraper) Scrape(m types.Media) (types.Media, error) {
	return nil, provider.ErrMediaNotSupported{}
}

func TestRenameUnsupportedScrapers(t *testing.T) {
	config := defaultConfig
	config.scrapers = []types.Scraper{
		fakeUnsupportedScraper{},
		fakeUnsupportedScraper{},
		fakeUnsupportedScraper{},
		fakeScraper{},
		fakeUnsupportedScraper{},
	}

	err := performRenameTest(t, copyTester{}, config)
	assert.NoError(t, err)

	cleanRenameTest(t)
}

type fakeErrorScraper struct{}

func (fakeErrorScraper) Scrape(m types.Media) (types.Media, error) {
	return nil, errors.New("mocked error")
}

func TestRenameScrapeError(t *testing.T) {
	config := defaultConfig
	config.scrapers = []types.Scraper{
		fakeUnsupportedScraper{},
		fakeErrorScraper{},
		fakeScraper{},
		fakeUnsupportedScraper{},
	}

	err := performRenameTest(t, copyTester{}, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mocked error")

	cleanRenameTest(t)
}

func TestRenameNoScrapers(t *testing.T) {
	config := defaultConfig
	config.scrapers = []types.Scraper{
		fakeUnsupportedScraper{},
		fakeUnsupportedScraper{},
		fakeUnsupportedScraper{},
	}

	err := performRenameTest(t, copyTester{}, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no scrapers")

	cleanRenameTest(t)
}

func performRenameTest(t *testing.T, test renameTester, cfg types.Config) error {
	app := New(cfg)

	test.Pre(t)

	l, err := app.FindMedia(test.Input())
	require.NoError(t, err)
	assert.Equal(t, len(res), l.Len())

	err = app.RenameMedia(l)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(test.Output())
	require.NoError(t, err)

	for _, f := range files {
		org, ok := res[f.Name()]
		assert.True(t, ok, f.Name())
		src := filepath.Join(test.Input(), org)
		dst := filepath.Join(test.Output(), f.Name())
		test.Test(t, src, dst)
	}

	test.Post(t, files)

	return nil
}

func cleanRenameTest(t *testing.T) {
	err := os.RemoveAll("out")
	require.NoError(t, err)
}

// copyTester is used for testing copying of files
type copyTester struct{}

func (copyTester) Pre(t *testing.T) {}

func (copyTester) Input() string {
	return "test"
}

func (copyTester) Output() string {
	return "out"
}

func (copyTester) Test(t *testing.T, src, dst string) {
	f, err := os.Lstat(dst)
	assert.NoError(t, err)
	assert.Zero(t, f.Mode()&os.ModeSymlink)
	assert.True(t, f.Mode().IsRegular())
	o, err := os.Stat(src)
	assert.NoError(t, err)
	assert.Equal(t, f.Mode(), o.Mode())
}

func (copyTester) Post(t *testing.T, files []os.FileInfo) {
	assert.Equal(t, len(res), len(files))
}

// moveTester is used to test moving of files
type moveTester struct{}

func (moveTester) Pre(t *testing.T) {
	files, err := ioutil.ReadDir("test")
	require.NoError(t, err)
	assert.Equal(t, len(res), len(files))
	err = os.MkdirAll("out/from", os.ModePerm)
	require.NoError(t, err)
	for _, f := range files {
		require.False(t, f.IsDir())
		require.NotEmpty(t, f.Name())
		fi, err := os.Create(filepath.Join("out", "from", f.Name()))
		defer fi.Close()
		require.NoError(t, err)
	}
}

func (moveTester) Input() string {
	return "out/from"
}

func (moveTester) Output() string {
	return "out/to"
}

func (moveTester) Test(t *testing.T, src, dst string) {
	f, err := os.Lstat(dst)
	assert.NoError(t, err)
	assert.Zero(t, f.Mode()&os.ModeSymlink)
	assert.True(t, f.Mode().IsRegular())

	_, err = os.Lstat(src)
	assert.True(t, os.IsNotExist(err), "should not exist: %s", src)
}

func (moveTester) Post(t *testing.T, files []os.FileInfo) {
	assert.Equal(t, len(res), len(files))
}

// symlinkTester is used to test symlinking of files
type symlinkTester struct{}

func (symlinkTester) Pre(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping symlink test on windows (requires admin rights)")
	}
}

func (symlinkTester) Input() string {
	return "test"
}

func (symlinkTester) Output() string {
	return "out"
}

func (symlinkTester) Test(t *testing.T, src, dst string) {
	f, err := os.Lstat(dst)
	assert.NoError(t, err)
	assert.NotZero(t, f.Mode()&os.ModeSymlink)
	link, err := os.Readlink(dst)
	require.NoError(t, err)
	abssrc, err := filepath.Abs(src)
	assert.NoError(t, err)
	assert.Equal(t, abssrc, link)
	assert.False(t, f.Mode().IsRegular())
}

func (symlinkTester) Post(t *testing.T, files []os.FileInfo) {
	assert.Equal(t, len(res), len(files))
}

// dryTester is used to test no-op renaming of files with dry flag
type dryTester struct{}

func (dryTester) Pre(t *testing.T) {}

func (dryTester) Input() string {
	return "test"
}

func (dryTester) Output() string {
	return "out"
}

func (dryTester) Test(t *testing.T, src, dst string) {
	assert.Fail(t, "dry run should not rename files")
}

func (dryTester) Post(t *testing.T, files []os.FileInfo) {
	assert.Equal(t, 0, len(files))
}
