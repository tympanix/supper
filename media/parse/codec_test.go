package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tympanix/supper/media/meta/codec"
)

func TestCodec(t *testing.T) {
	assert.Equal(t, codec.X264, Codec("x264"))
	assert.Equal(t, codec.HEVC, Codec("this.is.a.hevc-test"))
}
