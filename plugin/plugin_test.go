package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tympanix/supper/media"
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

func TestPluginFromMapError(t *testing.T) {
	data := make(map[string]interface{})

	_, err := NewFromMap(data)
	assert.Error(t, err)

	data = map[string]interface{}{"name": 42}

	_, err = NewFromMap(data)
	assert.Error(t, err)
}

func TestPluginInvalid(t *testing.T) {
	p := Plugin{
		PluginName: "test",
		Exec:       "",
	}

	assert.Error(t, p.valid())

	p = Plugin{
		PluginName: "",
		Exec:       "echo hello world",
	}

	assert.Error(t, p.valid())
}

func TestPluginExecError(t *testing.T) {
	p := Plugin{
		PluginName: "error",
		Exec:       "blablablablablablabla",
	}

	s, err := media.NewLocalSubtitle("test/Inception.2010.720p.en.srt")
	require.NoError(t, err)

	err = p.Run(s)
	require.Error(t, err)

}
