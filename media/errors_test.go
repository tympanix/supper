package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUnknown(t *testing.T) {
	err := NewUnknownErr()
	assert.Contains(t, err.Error(), "unknown")
	assert.True(t, IsUnknown(err))
	assert.False(t, IsUnknown(nil))
}

func TestErrExists(t *testing.T) {
	err := NewExistsErr()
	assert.Contains(t, err.Error(), "exists")
	assert.True(t, IsExistsErr(err))
	assert.False(t, IsExistsErr(nil))
}
