package list

import (
	"io"
	"math/rand"
	"strconv"
	"testing"

	"github.com/fatih/set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

type fakesubtitles struct {
	fakelocal
	subs []language.Tag
}

func (s fakesubtitles) ExistingSubtitles() (types.SubtitleList, error) {
	var subs []types.Subtitle
	for _, l := range s.subs {
		subs = append(subs, subtitle{s.fakelocal, l, false})
	}
	return Subtitles(subs...), nil
}

func (fakesubtitles) SaveSubtitle(io.Reader, language.Tag) (types.LocalSubtitle, error) {
	return nil, errMock
}

var sampleLanguages = []language.Tag{
	language.English,
	language.German,
	language.Spanish,
	language.Italian,
	language.Portuguese,
	language.Chinese,
	language.Polish,
	language.Arabic,
}

type langSelector func(types.Media) []language.Tag

func genTestVideoSampleList(size int, l langSelector) types.VideoList {
	var video []types.Video
	for i := 0; i < size; i++ {
		fake := fakelocal{movie{name: strconv.Itoa(i), year: 1970 + i}}
		video = append(video, fakesubtitles{
			fake,
			l(fake),
		})
	}
	return NewVideo(video...)
}

var noLanguages = func(m types.Media) []language.Tag {
	return make([]language.Tag, 0)
}

func TestVideoNoMissingSubtitles(t *testing.T) {
	sampleSize := 128

	video := genTestVideoSampleList(sampleSize, func(m types.Media) []language.Tag {
		return []language.Tag{language.English}
	})

	require.Equal(t, sampleSize, video.Len())

	missing, err := video.FilterMissingSubs(set.New(language.English))
	require.NoError(t, err)
	assert.Equal(t, 0, missing.Len())
}

func TestVideoAllMissingSubtitles(t *testing.T) {
	sampleSize := 128

	video := genTestVideoSampleList(sampleSize, func(m types.Media) []language.Tag {
		return []language.Tag{language.Catalan}
	})

	require.Equal(t, sampleSize, video.Len())

	all, err := video.FilterMissingSubs(set.New(language.English))
	require.NoError(t, err)
	assert.Equal(t, sampleSize, all.Len())
}

func TestVideoListRandomPerm(t *testing.T) {
	sampleSize := 128

	r := rand.New(rand.NewSource(int64(1337)))

	video := genTestVideoSampleList(sampleSize, func(m types.Media) []language.Tag {
		return []language.Tag{
			sampleLanguages[r.Intn(len(sampleLanguages))],
		}
	})

	require.Equal(t, sampleSize, video.Len())

	var sum int

	for _, l := range sampleLanguages {
		f, err := video.FilterMissingSubs(set.New(l))
		require.NoError(t, err)
		sum += sampleSize - f.Len()
		for _, m := range f.List() {
			subs, err := m.ExistingSubtitles()
			require.NoError(t, err)
			for _, s := range subs.List() {
				assert.NotEqual(t, l, s.Language())
			}
		}
	}

	assert.Equal(t, sampleSize, sum)
}
