package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tympanix/supper/meta/misc"
)

func TestMiscDTS(t *testing.T) {
	l := Miscellaneous("this.is.a.DTS.test")
	assert.True(t, l.Has(misc.DTS))
}

func TestMiscHC(t *testing.T) {
	l := Miscellaneous("this.is.a.hc.video")
	assert.True(t, l.Has(misc.HC))
}

func TestMiscMultiple(t *testing.T) {
	l := Miscellaneous("this.has.both.dts.and.ac3.and.hc.in.string")
	assert.True(t, l.Has(misc.DTS))
	assert.True(t, l.Has(misc.AC3))
	assert.True(t, l.Has(misc.HC))
}

func TestNoMisc(t *testing.T) {
	l := Miscellaneous("this string has no misc tags")
	assert.Len(t, l, 0)
}
