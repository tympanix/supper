package plugin

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/media"
)

func TestPluginExecWindows(t *testing.T) {
	p := Plugin{
		PluginName: "test",
		Exec:       "echo %SUBTITLE% > test/exec-proof",
	}

	s, err := media.NewLocalSubtitle("test/Inception.2010.720p.en.srt")
	require.NoError(t, err)

	err = p.Run(s)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	file, err := os.Open("test/exec-proof")
	require.NoError(t, err)

	data, err := ioutil.ReadAll(file)
	require.NoError(t, err)

	assert.Equal(t, []byte("test/Inception.2010.720p.en.srt \r\n"), data)

	require.NoError(t, file.Close())

	assert.NoError(t, os.Remove("test/exec-proof"))
}
