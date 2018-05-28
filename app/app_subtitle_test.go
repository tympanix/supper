package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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
	Post(*testing.T)
}

var inception, _ = media.NewLocalFile("test/Inception.2010.720p.x264.mkv")
var gameofthrones, _ = media.NewLocalFile("test/Game.of.Thrones.s01e02.mp4")

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

	_, err = app.DownloadSubtitles(media, config.Languages())
	if err != nil {
		return err
	}

	test.Post(t)

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

func (successTester) Post(t *testing.T) {

}
