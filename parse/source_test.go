package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tympanix/supper/meta/source"
)

func TestSource(t *testing.T) {
	assert.Equal(t, source.WEBDL, Source("WEB-DL"))
	assert.Equal(t, source.BluRay, Source("BLURAY"))
	assert.Equal(t, source.Telesync, Source("this.should.be.TS.test"))
}
