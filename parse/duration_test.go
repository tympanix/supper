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
