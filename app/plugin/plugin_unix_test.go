// +build !windows

package plugin

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tympanix/supper/media"
)

func TestPluginExecUnix(t *testing.T) {
	p := Plugin{
		PluginName: "test",
		Exec:       "echo $SUBTITLE > test/exec-proof",
	}

	s, err := media.NewLocalSubtitle("test/Inception.2010.720p.en.srt")
	require.NoError(t, err)

	err = p.Run(s)
	require.NoError(t, err)

	file, err := os.Open("test/exec-proof")
	require.NoError(t, err)

	data, err := ioutil.ReadAll(file)
	require.NoError(t, err)

	assert.Equal(t, []byte("test/Inception.2010.720p.en.srt\n"), data)

	require.NoError(t, file.Close())

	assert.NoError(t, os.Remove("test/exec-proof"))
}
