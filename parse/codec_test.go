package parse

import (
	"testing"

	"github.com/tympanix/supper/meta/codec"
)

func TestCodec(t *testing.T) {
	if Codecs.FindTag("x264") != codec.X264 {
		t.Error("should be x264")
	}

	if Codecs.FindTag("this.is.a.hevc-test") != codec.HEVC {
		t.Error("should be hevc")
	}
}
