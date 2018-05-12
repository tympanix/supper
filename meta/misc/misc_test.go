package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiscHas(t *testing.T) {
	list := List([]Tag{
		Video3D,
		DTS,
		HC,
	})

	assert.True(t, list.Has(Video3D))
	assert.True(t, list.Has(DTS))
	assert.True(t, list.Has(HC))

	assert.False(t, list.Has(AC3))
}
