package parse

import (
	"fmt"
	"testing"

	"github.com/tympanix/supper/meta/source"
)

func TestSource(t *testing.T) {
	if Sources.FindTag("WEB-DL") != source.WEBDL {
		fmt.Println(Sources.FindTag("WEB-DL"), source.WEBDL)
		t.Error("should be WEBDL")
	}

	if Sources.FindTag("BLURAY") != source.BluRay {
		t.Error("should be BluRay")
	}

	if Sources.FindTag("this.should.be.TS.test") != source.Telesync {
		t.Error("should be Telesync")
	}
}
