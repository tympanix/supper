package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

var languages = []types.Subtitle{
	// inception
	subtitle{inception, language.English, false},
	subtitle{inception, language.German, true},
	subtitle{inception, language.Spanish, false},

	// fightclub
	subtitle{fightclub, language.German, false},
	subtitle{fightclub, language.French, true},

	// batman begins
	subtitle{batmanbegins, language.English, false},
	subtitle{batmanbegins, language.Italian, false},

	// the office
	subtitle{theoffice, language.Spanish, true},
	subtitle{theoffice, language.Italian, false},
	subtitle{theoffice, language.Chinese, false},

	// arrow
	subtitle{arrow, language.English, false},
	subtitle{arrow, language.Portuguese, false},

	// westworld
	subtitle{westworld, language.Polish, false},
}

func TestSubtitleList(t *testing.T) {
	subs := Subtitles(languages...)
	langs := subs.LanguageSet()

	// test subtitle languages
	for _, l := range []language.Tag{
		language.English,
		language.German,
		language.Italian,
		language.French,
		language.Spanish,
		language.Chinese,
		language.Portuguese,
		language.Polish,
	} {
		f := subs.FilterLanguage(l)
		require.NotEqual(t, 0, f.Len())
		assert.True(t, langs.Has(l))

		assert.True(t, f.LanguageSet().Has(l))
		assert.Equal(t, 1, f.LanguageSet().Size())

		for _, s := range f.List() {
			assert.Equal(t, l, s.Language())
		}
	}

	// test hearing impaired subtitles
	for _, b := range []bool{
		true,
		false,
	} {
		f := subs.HearingImpaired(b)
		require.NotEqual(t, 0, f.Len())
		for _, s := range f.List() {
			assert.Equal(t, b, s.HearingImpaired())
		}
	}

}

func TestSubtitlesFromInterface(t *testing.T) {
	list, err := NewSubtitlesFromInterface(languages)
	require.NoError(t, err)
	assert.Equal(t, len(languages), list.Len())
	for _, s := range languages {
		assert.Contains(t, list.List(), s)
	}
}

func TestSubtitlesFromInterfaceError(t *testing.T) {
	list, err := NewSubtitlesFromInterface(42)
	assert.Error(t, err)
	assert.Nil(t, list)
}

func TestSubtitleRateByMedia(t *testing.T) {
	list := Subtitles(languages...)
	rated := list.RateByMedia(inception)
	for _, s := range rated.List() {
		assert.Equal(t, inception, s.Subtitle().ForMedia())
		assert.True(t, s.Score() > 0.0)
		assert.True(t, s.Score() <= 1.1)
	}
}
