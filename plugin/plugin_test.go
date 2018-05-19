package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginFromMap(t *testing.T) {
	data := map[string]string{
		"name": "test_plugin",
		"exec": "echo hello world",
	}

	p, err := NewFromMap(data)
	require.NoError(t, err)

	assert.Equal(t, "test_plugin", p.PluginName)
	assert.Equal(t, "echo hello world", p.Exec)
}
