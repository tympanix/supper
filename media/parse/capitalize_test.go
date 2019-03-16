package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	for _, v := range testMovieTitles {
		assert.Equal(t, v, Capitalize(strings.ToLower(v)))
	}
}

func TestUpperLowerRatio(t *testing.T) {
	assert.Equal(t, true, isUpper("ABC"))
	assert.Equal(t, true, isUpper("ABC123"))
	assert.Equal(t, true, isUpper("ABC_ABC:,'%/(!¤%))"))
	assert.Equal(t, true, isUpper("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))

	assert.Equal(t, false, isUpper("abc"))
	assert.Equal(t, false, isUpper("ABCDEFGHIJKLmNOPQRSTUVWXYZ"))
	assert.Equal(t, false, isUpper("123"))
}
