package provider

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"github.com/xrash/smetrics"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/language"
)

// HOST is the URL for subscene
const HOST = "https://subscene.com"

// DELAY is the delay between calls to subscene to prevent spamming
const DELAY = 500 * time.Millisecond

var lock = new(sync.Mutex)

// lockSubscene is used to limit the number of calls to subscene to prevent spamming
func lockSubscene() {
	lock.Lock()
	go func() {
		time.Sleep(DELAY)
		lock.Unlock()
	}()
}

// Subscene interfaces with subscene.com for downloading subtitles
type Subscene struct{}

func (s *Subscene) searchTerm(m types.Media) string {
	if movie, ok := m.TypeMovie(); ok {
		return s.searchTermMovie(movie)
	} else if episode, ok := m.TypeEpisode(); ok {
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
		return "", fmt.Errorf("unable to search for media: %s", media)
	}

	query := url.Query()
	query.Add("q", s.cleanSearchTerm(search))
	url.RawQuery = query.Encode()

	lockSubscene()
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
		return "", errors.New("no media found on subscene.com")
	}

	return fmt.Sprintf("%s%s", "https://subscene.com", result), nil
}

// SearchSubtitles searches subscene.com for subtitles
func (s *Subscene) SearchSubtitles(local types.LocalMedia) (subs types.SubtitleList, err error) {
	url, err := s.FindMediaURL(local)

	if err != nil {
		return
	}

	lockSubscene()
	doc, err := goquery.NewDocument(url)

	if err != nil {
		return
	}

	subs = list.RatedSubtitles(local)

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

		meta, err := media.NewMetadata(name)

		if err != nil {
			return
		}

		langTag, err := parse.Language(lang)

		if err != nil {
			return
		}

		subs.Add(&subsceneSubtitle{
			Media:        media.NewType(meta),
			Downloadable: subsceneURL(url),
			lang:         langTag,
			comment:      comm,
			hi:           hi,
		})
	})

	return
}

func newZipReader(file *os.File) (*zipReader, error) {
	data, err := zip.OpenReader(file.Name())

	if err != nil {
		return nil, err
	}

	var srt io.ReadCloser = nil
	for _, f := range data.File {
		if filepath.Ext(f.Name) == ".srt" {
			if srt, err = f.Open(); err != nil {
				return nil, err
			}
			break
		}
	}

	if srt == nil {
		return nil, errors.New("could not read srt file from subscene")
	}

	return &zipReader{srt, data, file}, nil
}

type zipReader struct {
	io.ReadCloser
	zip  *zip.ReadCloser
	file *os.File
}

func (t *zipReader) Close() error {
	t.ReadCloser.Close()
	t.zip.Close()
	t.file.Close()
	if err := os.Remove(t.file.Name()); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type subsceneURL string

func (url subsceneURL) Download() (io.ReadCloser, error) {

	uri := fmt.Sprintf("%s%s", HOST, string(url))

	lockSubscene()
	doc, err := goquery.NewDocument(uri)

	if err != nil {
		return nil, err
	}

	sel := doc.Find("#downloadButton")

	if len(sel.Nodes) == 0 {
		return nil, errors.New("could not parse response from subscene")
	}

	download, exists := sel.First().Attr("href")

	if !exists {
		return nil, errors.New("could not find download link from subscene")
	}

	download = fmt.Sprintf("%s%s", HOST, download)

	resp, err := http.Get(download)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("subscene download subtitle (%v)", resp.StatusCode)
	}

	file, err := ioutil.TempFile("", "supper")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, err
	}

	return newZipReader(file)
}

type subsceneSubtitle struct {
	types.Media
	types.Downloadable
	lang    language.Tag
	comment string
	hi      bool
}

func (b *subsceneSubtitle) String() string {
	return fmt.Sprintf("%-15s %-s", b.lang, b.Media)
}

func (b *subsceneSubtitle) Language() language.Tag {
	return b.lang
}

func (b *subsceneSubtitle) IsLang(tag language.Tag) bool {
	return b.lang == tag
}

func (b *subsceneSubtitle) IsHI() bool {
	return b.hi
}
