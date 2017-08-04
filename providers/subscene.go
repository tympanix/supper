package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tympanix/supper/media"
	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Subscene interfaces with subscene.com for downloading subtitles
type Subscene struct{}

func (s *Subscene) searchTerm(m types.Media) string {
	if movie, ok := m.(*media.Movie); ok {
		return s.searchTermMovie(movie)
	} else if episode, ok := m.(*media.Episode); ok {
		return s.searchTermEpisode(episode)
	}
	return ""
}

func (s *Subscene) searchTermMovie(movie *media.Movie) string {
	return fmt.Sprintf("%s (%d)", movie.Name(), movie.Year())
}

func (s *Subscene) searchTermEpisode(episode *media.Episode) string {
	season := parse.PhoneticNumber(episode.Season())
	return fmt.Sprintf("%s - %s Season", episode.Name(), season)
}

// Search searches subscene.com for subtitles
func (s *Subscene) Search(media types.Media) (subs []types.Subtitle, err error) {
	doc, err := goquery.NewDocument("https://subscene.com/subtitles/guardians-of-the-galaxy")

	if err != nil {
		return
	}

	subs = make([]types.Subtitle, 0)

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		a1 := s.Find(".a1")

		url, exists := a1.Find("a").Attr("href")

		if !exists {
			return
		}

		spans := a1.Find("a span")

		lang := strings.TrimSpace(spans.First().Text())
		name := strings.TrimSpace(spans.Next().Text())
		comm := strings.TrimSpace(s.Find(".a6 div").Text())

		hi := s.Find("td.a41").Length() > 0

		subs = append(subs, &BasicSubtitle{
			name,
			url,
			lang,
			comm,
			hi,
		})
	})

	return
}

type BasicSubtitle struct {
	filename string
	url      string
	lang     string
	comment  string
	hi       bool
}

func (b *BasicSubtitle) Name() string {
	return b.filename
}

func (b *BasicSubtitle) Download() *os.File {
	return nil
}

func (b *BasicSubtitle) String() string {
	return fmt.Sprintf("%-15s %-s", b.lang, b.filename)
}

func (b *BasicSubtitle) IsLang(tag language.Tag) bool {
	return strings.Contains(b.lang, display.Self.Name(tag))
}

func (b *BasicSubtitle) IsHI() bool {
	return b.hi
}
