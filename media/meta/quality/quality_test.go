package quality

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestQualityString(t *testing.T) {
  assert.Equal(t, "1080p", HD1080p.String())
  assert.Equal(t, "720p", HD720p.String())
  assert.Equal(t, "UHD", UHD2160p.String())
}
