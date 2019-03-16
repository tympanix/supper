package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tympanix/supper/media/meta/quality"
)

func TestQuality(t *testing.T) {
	assert.Equal(t, quality.UHD2160p, Quality("4k"))
	assert.Equal(t, quality.HD1080p, Quality("1080P"))
	assert.Equal(t, quality.HD720p, Quality("this.media.is.720p.test"))
}
