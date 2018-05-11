package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumbers(t *testing.T) {
	assert.Equal(t, "", PhoneticNumber(-1))
	assert.Equal(t, "Second", PhoneticNumber(2))
	assert.Equal(t, "Eleventh", PhoneticNumber(11))
	assert.Equal(t, "Thirtieth", PhoneticNumber(30))
	assert.Equal(t, "Twenty-Fourth", PhoneticNumber(24))
	assert.Equal(t, "", PhoneticNumber(1<<16)) /* number too big */
}
