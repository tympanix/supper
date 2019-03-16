package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, "group", Group("x264.1080p.group"))
	assert.Equal(t, "", Group("DivX.1080p"))
	assert.Equal(t, "", Group("x264"))
	assert.Equal(t, "", Group("DVDRip"))
	assert.Equal(t, "", Group("x264.1080p.s"))
}
