package app

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/fatih/set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
)

type fakeProvider struct{}

func (f fakeProvider) ResolveSubtitle(link types.Linker) (types.Downloadable, error) {
	return nil, nil
}
func (f fakeProvider) SearchSubtitles(m types.LocalMedia) ([]types.OnlineSubtitle, error) {
	return nil, nil
}

type fakeConfig struct {
	action   string
	scraper  types.Scraper
	provider types.Provider
	movies   types.MediaConfig
	tvshows  types.MediaConfig
}

func (c fakeConfig) Languages() set.Interface       { return nil }
func (c fakeConfig) APIKeys() types.APIKeys         { return fakeAPIKeys{} }
func (c fakeConfig) Config() string                 { return "" }
func (c fakeConfig) Delay() time.Duration           { return 0 }
func (c fakeConfig) Dry() bool                      { return false }
func (c fakeConfig) Force() bool                    { return false }
func (c fakeConfig) Impaired() bool                 { return false }
func (c fakeConfig) Limit() int                     { return -1 }
func (c fakeConfig) Logfile() string                { return "" }
func (c fakeConfig) MediaFilter() types.MediaFilter { return nil }
func (c fakeConfig) Modified() time.Duration        { return 0 }
func (c fakeConfig) Movies() types.MediaConfig      { return c.movies }
func (c fakeConfig) Plugins() []types.Plugin        { return nil }
func (c fakeConfig) Score() int                     { return 0 }
func (c fakeConfig) Strict() bool                   { return true }
func (c fakeConfig) TVShows() types.MediaConfig     { return c.tvshows }
func (c fakeConfig) Verbose() bool                  { return false }
func (c fakeConfig) Providers() []types.Provider    { return []types.Provider{c.provider} }
func (c fakeConfig) Scrapers() []types.Scraper      { return []types.Scraper{c.scraper} }
func (c fakeConfig) RenameAction() string           { return c.action }

type fakeAPIKeys struct{}

func (k fakeAPIKeys) TheMovieDB() string { return "" }
func (k fakeAPIKeys) TheTVDB() string    { return "" }

type fakeMediaConfig struct {
	directory string
	template  string
}

func (m fakeMediaConfig) Directory() string { return m.directory }
func (m fakeMediaConfig) Template() *template.Template {
	t, err := template.New("test").Parse(m.template)
	if err != nil {
		panic(err.Error())
	}
	return t
}

type fakeScraper []types.Media

func (s fakeScraper) Scrape(m types.Media) (types.Media, error) {
	if s, ok := m.TypeSubtitle(); ok {
		return s.ForMedia(), nil
	}
	return m, nil
}

func genTestConfig(action string, output string) types.Config {
	return fakeConfig{
		action:  action,
		scraper: fakeScraper{},
		movies: fakeMediaConfig{
			directory: output,
			template:  "{{ .Movie }} ({{ .Year }}) {{ .Quality }}",
		},
		tvshows: fakeMediaConfig{
			directory: output,
			template:  "{{ .TVShow }} S{{ .Season }}E{{ .Episode }}",
		},
	}
}

var res = map[string]string{
	"Inception (2010) 720p.mkv":    "Inception.2010.720p.x264.mkv",
	"Inception (2010) 720p.en.srt": "Inception.2010.720p.x264.en.srt",
	"Game of Thrones S1E2.mp4":     "Game.of.Thrones.s01e02.mp4",
	"Game of Thrones S1E2.en.srt":  "Game.of.Thrones.s01e02.en.srt",
}

type renameTester interface {
	Skip() bool
	Pre(*testing.T)
	Input() string
	Output() string
	Test(*testing.T, string, string)
}

func TestRenameMedia(t *testing.T) {
	for action, test := range map[string]renameTester{
		"copy":     copyTester{},
		"move":     moveTester{},
		"hardlink": copyTester{},
		"symlink":  symlinkTester{},
	} {
		config := genTestConfig(action, test.Output())
		app := New(config)

		if test.Skip() {
			continue
		}

		test.Pre(t)

		l, err := app.FindMedia(test.Input())
		require.NoError(t, err)
		assert.Equal(t, len(res), l.Len())

		err = app.RenameMedia(l)
		require.NoError(t, err)

		files, err := ioutil.ReadDir(test.Output())
		require.NoError(t, err)
		assert.Equal(t, len(res), len(files))

		for _, f := range files {
			org, ok := res[f.Name()]
			assert.True(t, ok, f.Name())
			src := filepath.Join(test.Input(), org)
			dst := filepath.Join(test.Output(), f.Name())
			test.Test(t, src, dst)
		}

		err = os.RemoveAll("test/out")
		require.NoError(t, err)
	}

}

type copyTester struct{}

func (copyTester) Skip() bool {
	return false
}

func (copyTester) Pre(t *testing.T) {}

func (copyTester) Input() string {
	return "test"
}

func (copyTester) Output() string {
	return "test/out"
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

type moveTester struct{}

func (moveTester) Skip() bool {
	return false
}

func (moveTester) Pre(t *testing.T) {
	files, err := ioutil.ReadDir("test")
	require.NoError(t, err)
	assert.Equal(t, len(res), len(files))
	err = os.MkdirAll("test/out/from", os.ModePerm)
	require.NoError(t, err)
	for _, f := range files {
		require.False(t, f.IsDir())
		require.NotEmpty(t, f.Name())
		fi, err := os.Create(filepath.Join("test", "out", "from", f.Name()))
		defer fi.Close()
		require.NoError(t, err)
	}
}

func (moveTester) Input() string {
	return "test/out/from"
}

func (moveTester) Output() string {
	return "test/out/to"
}

func (moveTester) Test(t *testing.T, src, dst string) {
	f, err := os.Lstat(dst)
	assert.NoError(t, err)
	assert.Zero(t, f.Mode()&os.ModeSymlink)
	assert.True(t, f.Mode().IsRegular())

	_, err = os.Lstat(src)
	assert.True(t, os.IsNotExist(err), "should not exist: %s", src)
}

type symlinkTester struct{}

func (symlinkTester) Skip() bool {
	if runtime.GOOS == "windows" {
		log.Println("skipping symlink test on windows (requires admin rights)")
		return true
	}
	return false
}

func (symlinkTester) Pre(t *testing.T) {}

func (symlinkTester) Input() string {
	return "test"
}

func (symlinkTester) Output() string {
	return "test/out"
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
