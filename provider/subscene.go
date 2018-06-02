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

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"github.com/xrash/smetrics"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// subsceneHost is the URL for subscene
const subsceneHost = "https://subscene.com"

// subsceneDelay is the delay between calls to subscene to prevent spamming
const subsceneDelay = 500 * time.Millisecond

var subsceneLock = new(sync.Mutex)

// lockSubscene is used to limit the number of calls to subscene to prevent spamming
func lockSubscene() {
	subsceneLock.Lock()
	go func() {
		time.Sleep(subsceneDelay)
		subsceneLock.Unlock()
	}()
}

// Subscene interfaces with subscene.com for downloading subtitles
func Subscene() types.Provider {
	return &subscene{}
}

type subscene struct{}

func (s *subscene) ResolveSubtitle(l types.Linker) (types.Downloadable, error) {
	return subsceneURL(l.Link()), nil
}

func (s *subscene) searchTerm(m types.Media) string {
	if movie, ok := m.TypeMovie(); ok {
		return s.searchTermMovie(movie)
	} else if episode, ok := m.TypeEpisode(); ok {
		return s.searchTermEpisode(episode)
	}
	return ""
}

func (s *subscene) searchTermMovie(movie types.Movie) string {
	return movie.MovieName()
}

func (s *subscene) searchTermEpisode(episode types.Episode) string {
	season := parse.PhoneticNumber(episode.Season())
	return fmt.Sprintf("%s - %s Season", episode.TVShow(), season)
}

func (s *subscene) filterTerm(m types.Media) string {
	if movie, ok := m.TypeMovie(); ok {
		return fmt.Sprintf("%s (%v)", movie.MovieName(), movie.Year())
	} else if episode, ok := m.TypeEpisode(); ok {
		return episode.TVShow()
	}
	return ""
}

type searchResult struct {
	Title string
	URL   string
}

var cleanRegexp = regexp.MustCompile(`\(\d{4}\)`)

func (s *subscene) cleanSearchTerm(search string) string {
	search = cleanRegexp.ReplaceAllString(search, "")
	return strings.TrimSpace(search)
}

// FindMediaURL retrieves the subscene.com URL for the given media item
func (s *subscene) FindMediaURL(media types.Media, retries int) (string, error) {
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
	res, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 409 && retries > 0 {
		// Busy, try again in two seconds
		log.WithField("media", media).WithField("status", 409).
			WithField("retries", retries).
			Debug("Retrying subscene.com")
		time.Sleep(1500 * time.Millisecond)
		return s.FindMediaURL(media, retries-1)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("subscene.com returned status code %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

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

	filter := s.filterTerm(media)

	if filter == "" {
		return "", errors.New("subscene unknown media")
	}

	var result string
	min := math.MaxInt32
	for url, name := range results {
		score := smetrics.WagnerFischer(filter, name, 1, 1, 2)
		if score < min {
			min = score
			result = url
		}
	}

	if len(result) == 0 {
		return "", errors.New("no media found on subscene.com")
	}

	return fmt.Sprintf("%s%s", subsceneHost, result), nil
}

// SearchSubtitles searches subscene.com for subtitles
func (s *subscene) SearchSubtitles(local types.LocalMedia) (subs []types.OnlineSubtitle, err error) {
	url, err := s.FindMediaURL(local, 3)

	if err != nil {
		return
	}

	lockSubscene()
	doc, err := goquery.NewDocument(url)

	if err != nil {
		return
	}

	subs = make([]types.OnlineSubtitle, 0)

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

		meta, err := media.NewFromString(name)

		if err != nil {
			return
		}

		langTag, err := parse.Language(lang)

		if err != nil {
			return
		}

		subs = append(subs, &subsceneSubtitle{
			Media:       meta,
			subsceneURL: subsceneURL(url),
			lang:        langTag,
			comment:     comm,
			hi:          hi,
		})
	})

	return
}

func newZipReader(file *os.File) (*zipReader, error) {
	data, err := zip.OpenReader(file.Name())

	if err != nil {
		return nil, err
	}

	var srt io.ReadCloser
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
		log.WithError(err).Error("Could not cleaup temporary zip file")
		return err
	}
	return nil
}

type subsceneURL string

func (uri subsceneURL) Link() string {
	return string(uri)
}

func (uri subsceneURL) Download() (io.ReadCloser, error) {
	fulluri := fmt.Sprintf("%s%s", subsceneHost, string(uri))

	suburl, err := url.ParseRequestURI(fulluri)

	if err != nil {
		return nil, err
	}

	lockSubscene()
	doc, err := goquery.NewDocument(suburl.String())

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

	download = fmt.Sprintf("%s%s", subsceneHost, download)

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
	subsceneURL
	lang    language.Tag
	comment string
	hi      bool
}

func (b *subsceneSubtitle) String() string {
	return display.English.Languages().Name(b.Language())
}

func (b *subsceneSubtitle) ForMedia() types.Media {
	return b.Media
}

func (b *subsceneSubtitle) Language() language.Tag {
	return b.lang
}

func (b *subsceneSubtitle) HearingImpaired() bool {
	return b.hi
}
