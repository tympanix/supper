package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestParseLanguage(t *testing.T) {
	l, err := Language("English")
	assert.NoError(t, err)
	assert.Equal(t, language.English, l)
}

func TestParseLanguageError(t *testing.T) {
	_, err := Language("blablabla")
	assert.Error(t, err)
}
