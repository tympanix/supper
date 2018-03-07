package parse

import (
	"testing"

	"github.com/tympanix/supper/meta/quality"
)

func TestQuality(t *testing.T) {
	if Qualities.FindTag("4k") != quality.UHD2160p {
		t.Error("should be 4K")
	}

	if Qualities.FindTag("1080P") != quality.HD1080p {
		t.Error("should be 1080p")
	}

	if Qualities.FindTag("this.media.is.720p.test") != quality.HD720p {
		t.Error("should be 820p")
	}
}
