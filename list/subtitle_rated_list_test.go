package list

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
)

type fakeRatedSubtitle struct {
	types.Media
	types.Subtitle
	score float32
}

func (s fakeRatedSubtitle) ForMedia() types.Media {
	return s
}

func (s fakeRatedSubtitle) String() string {
	return fmt.Sprintf("%v", s.score)
}

type fakeEvaluator struct{ *testing.T }

func (t fakeEvaluator) Evaluate(m types.Media, n types.Media) float32 {
	if f, ok := n.(fakeRatedSubtitle); ok {
		return f.score
	}
	t.Error("Expected fake subtitles")
	return 0.0
}

func genTestRatedSubtitleSampleList(size int) []types.Subtitle {
	var subs []types.Subtitle
	r := rand.New(rand.NewSource(1337))
	for i := 0; i < size; i++ {
		subs = append(subs, fakeRatedSubtitle{score: r.Float32()})
	}
	return subs
}

func testRatedSubtitlesDescending(t *testing.T, l types.RatedSubtitleList) {
	for i := 1; i < l.Len(); i++ {
		fst := l.List()[i-1]
		snd := l.List()[i]
		assert.True(t, fst.Score() >= snd.Score())
	}
}

func TestRatedSubtitleList(t *testing.T) {
	var sampleSize = 1024
	subs := genTestRatedSubtitleSampleList(sampleSize)
	rated := NewRatedSubtitles(inception, fakeEvaluator{t}, subs...)

	require.Equal(t, sampleSize, rated.Len())

	testRatedSubtitlesDescending(t, rated)
	assert.Equal(t, rated.List()[0], rated.Best())

	filtered := rated.FilterScore(0.5)

	for _, s := range filtered.List() {
		assert.True(t, s.Score() >= 0.5)
	}

	testRatedSubtitlesDescending(t, filtered)
}
