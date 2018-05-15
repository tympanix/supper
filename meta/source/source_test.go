package source

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestSourceString(t *testing.T) {
  assert.Equal(t, "BluRay", BluRay.String())
  assert.Equal(t, "DVDRip", DVDRip.String())
  assert.Equal(t, "Remux", Remux.String())
}
