package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodecString(t *testing.T) {
	assert.Equal(t, "x264", X264.String())
  assert.Equal(t, "x265", X265.String())
  assert.Equal(t, "DivX", DivX.String())
}
