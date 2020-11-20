package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/types"
)

type extractTest interface {
	Pre(*testing.T, []types.MediaArchive)
	Input() string
}

func TestAppArchives(t *testing.T) {
	defer cleanRenameTest(t)

	config := defaultConfig

	// first run, should succeed
	err := performExtractionTest(t, extractOkTest{}, config)
	require.NoError(t, err)

	files, err := ioutil.ReadDir("out")
	require.NoError(t, err)
	assert.Equal(t, len(files), 4)

	// second run, should fail (no overwrite)
	err = performExtractionTest(t, extractOkTest{}, config)
	require.Error(t, err)
	assert.True(t, media.IsExistsErr(err))
}

func TestAppArchivesInvalidPath(t *testing.T) {
	config := defaultConfig

	app := New(config)

	_, err := app.FindArchives("deosnotexist")
	assert.Error(t, err)
}

func TestAppArchiveDryRun(t *testing.T) {
	defer cleanRenameTest(t)

	config := defaultConfig

	config.dry = true

	err := performExtractionTest(t, extractDryTest{}, config)
	require.NoError(t, err)

	_, err = ioutil.ReadDir("out")
	require.Error(t, err)
}

func performExtractionTest(t *testing.T, test extractTest, config types.Config) error {
	app := New(config)

	arch, err := app.FindArchives(test.Input())
	if err != nil {
		return err
	}

	test.Pre(t, arch)

	for _, a := range arch {
		f, err := a.Next()

		for err != io.EOF {
			fmt.Println(f.String())
			if err = app.ExtractMedia(f); err != nil {
				return err
			}
			f, err = a.Next()
		}
	}
	return nil
}

type extractOkTest struct{}

func (extractOkTest) Pre(t *testing.T, a []types.MediaArchive) {
	assert.Equal(t, len(a), 2)
}

func (extractOkTest) Input() string {
	return "../test/archives"
}

type extractDryTest struct{}

func (extractDryTest) Pre(t *testing.T, a []types.MediaArchive) {
	assert.Equal(t, len(a), 2)
}

func (extractDryTest) Input() string {
	return "../test/archives"
}
