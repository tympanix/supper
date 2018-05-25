package app

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/tympanix/supper/media"

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
	action string
	strict bool
	force  bool
}

func (c fakeConfig) Languages() set.Interface       { return nil }
func (c fakeConfig) APIKeys() types.APIKeys         { return fakeAPIKeys{} }
func (c fakeConfig) Config() string                 { return "" }
func (c fakeConfig) Delay() time.Duration           { return 0 }
func (c fakeConfig) Dry() bool                      { return false }
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
func (c fakeConfig) Providers() []types.Provider    { return []types.Provider{fakeProvider{}} }
func (c fakeConfig) Scrapers() []types.Scraper      { return []types.Scraper{fakeScraper{}} }
func (c fakeConfig) RenameAction() string           { return c.action }

type fakeTemplates struct {
	output string
}

func (c fakeTemplates) Movies() types.MediaConfig {
	return fakeMediaConfig{
		directory: c.output,
		template:  "{{ .Movie }} ({{ .Year }}) {{ .Quality }}",
	}
}

func (c fakeTemplates) TVShows() types.MediaConfig {
	return fakeMediaConfig{
		directory: c.output,
		template:  "{{ .TVShow }} S{{ .Season }}E{{ .Episode }}",
	}
}

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
			config := struct {
				fakeConfig
				fakeTemplates
			}{
				fakeConfig{
					action: action,
					strict: true,
				},
				fakeTemplates{
					output: test.Output(),
				},
			}

			assert.NoError(t, performRenameTest(t, test, config))

			cleanRenameTest(t)
		})
	}
}

func TestRenameMediaOverride(t *testing.T) {
	var err error

	config := struct {
		fakeConfig
		fakeTemplates
	}{
		fakeConfig{
			action: "copy",
			strict: false,
		},
		fakeTemplates{
			output: "out",
		},
	}

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
	assert.Equal(t, len(res), len(files))

	for _, f := range files {
		org, ok := res[f.Name()]
		assert.True(t, ok, f.Name())
		src := filepath.Join(test.Input(), org)
		dst := filepath.Join(test.Output(), f.Name())
		test.Test(t, src, dst)
	}

	return nil
}

func cleanRenameTest(t *testing.T) {
	err := os.RemoveAll("out")
	require.NoError(t, err)
}

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
