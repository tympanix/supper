package app

import (
	"testing"

	"github.com/tympanix/supper/app/cfg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppFindMediaInvalidPath(t *testing.T) {
	config := defaultConfig

	app := New(config)

	_, err := app.FindMedia("doesnotexist")
	assert.Error(t, err)
}

func TestAppFindMedia(t *testing.T) {
	config := defaultConfig

	app := New(config)

	media, err := app.FindMedia("../test/find")
	require.NoError(t, err)

	assert.Equal(t, media.Len(), 1)
}

func TestAppFromDefault(t *testing.T) {
	cfg.Initialize()

	app := NewFromDefault()
	assert.Equal(t, app.Config(), cfg.Default)
}
