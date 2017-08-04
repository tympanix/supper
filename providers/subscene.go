package provider

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/Tympanix/supper/media"
	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
	"github.com/xrash/smetrics"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Subscene interfaces with subscene.com for downloading subtitles
type Subscene struct{}

func (s *Subscene) searchTerm(m types.Media) string {
	if movie, ok := m.(types.Movie); ok {
		return s.searchTermMovie(movie)
	} else if episode, ok := m.(types.Episode); ok {
		return s.searchTermEpisode(episode)
	}
	return ""
}

func (s *Subscene) searchTermMovie(movie types.Movie) string {
	return movie.MovieName()
}

func (s *Subscene) searchTermEpisode(episode types.Episode) string {
	season := parse.PhoneticNumber(episode.Season())
	return fmt.Sprintf("%s - %s Season", episode.TVShow(), season)
}

type searchResult struct {
	Title string
	URL   string
}

var cleanRegexp = regexp.MustCompile(`\(\d{4}\)`)

func (s *Subscene) cleanSearchTerm(search string) string {
	search = cleanRegexp.ReplaceAllString(search, "")
	return strings.TrimSpace(search)
}

// FindMediaURL retrieves the subscene.com URL for the given media item
func (s *Subscene) FindMediaURL(media types.Media) (string, error) {
	url, err := url.Parse("https://subscene.com/subtitles/title")

	if err != nil {
		return "", err
	}

	search := s.searchTerm(media)

	if len(search) == 0 {
		return "", fmt.Errorf("Unable to search for media: %s", media)
	}

	query := url.Query()
	query.Add("q", s.cleanSearchTerm(search))
	url.RawQuery = query.Encode()

	doc, err := goquery.NewDocument(url.String())

	if err != nil {
		return "", err
	}

	results := make(map[string]string, 0)

	doc.Find("div.search-result div.title a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url, ok := s.Attr("href")
		if len(title) > 0 && ok {
			results[url] = title
		}
	})

	var result string
	min := math.MaxInt32
	for url, name := range results {
		score := smetrics.WagnerFischer(search, name, 1, 1, 2)
		if score < min {
			min = score
			result = url
		}
	}

	if len(result) == 0 {
		return "", errors.New("No media found on subscene.com")
	}

	log.Println(result)
	return fmt.Sprintf("%s%s", "https://subscene.com", result), nil
}

// SearchSubtitles searches subscene.com for subtitles
func (s *Subscene) SearchSubtitles(local types.LocalMedia) (subs []types.Subtitle, err error) {
	url, err := s.FindMediaURL(local)

	if err != nil {
		return
	}

	doc, err := goquery.NewDocument(url)

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

		subs = append(subs, &subsceneSubtitle{
			Media:        media.NewFromFilename(name),
			Downloadable: subsceneDownloader(url),
			lang:         lang,
			comment:      comm,
			hi:           hi,
		})
	})

	return
}

type subsceneDownloader string

func (s subsceneDownloader) Download() io.Reader {
	return nil
}

type subsceneSubtitle struct {
	types.Media
	types.Downloadable
	lang    string
	comment string
	hi      bool
}

func (b *subsceneSubtitle) String() string {
	return fmt.Sprintf("%-15s %-s", b.lang, b.Media)
}

func (b *subsceneSubtitle) IsLang(tag language.Tag) bool {
	return strings.Contains(b.lang, display.Self.Name(tag))
}

func (b *subsceneSubtitle) IsHI() bool {
	return b.hi
}
