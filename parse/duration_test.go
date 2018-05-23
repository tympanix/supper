package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	d, err := Duration("2h3m")
	require.NoError(t, err)
	assert.Equal(t, d, 2*time.Hour+3*time.Minute)
}

func TestDurationError(t *testing.T) {
	for _, str := range []string{
		"123",
		"12/#m23?s",
		"42q",
		"_pppp",
		"????",
	} {
		d, err := Duration(str)
		assert.Error(t, err)
		assert.Equal(t, time.Duration(0), d)
	}
}
