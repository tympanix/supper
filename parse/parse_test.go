package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	assert.Equal(t, "thisisatest", Identity("thìs is â tést"))
	assert.Equal(t, "vyzkousejtetentoretezec", Identity("vyzkoušejte tento řetězec"))
	assert.Equal(t, "abc123", Identity(`"?=_ä!'<b½c)#1,2...3`))
	assert.Equal(t, "这是一个测试", Identity("这是一个测试"))
}

func TestCleanName(t *testing.T) {
	assert.Equal(t, "this is a test", CleanName("this.is.a.test"))
	assert.Equal(t, "this is a (test)", CleanName("this?_=is#.,a_(test)"))
	assert.Equal(t, "abc ABC abc", CleanName("abc.A.B.C.abc"))
	assert.Equal(t, "abc ABC abc", CleanName("abc A B C abc"))
}
