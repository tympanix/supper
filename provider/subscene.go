package provider

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nwaples/rardecode"
	"github.com/xrash/smetrics"

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"

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
	Path  string
}

var cleanRegexp = regexp.MustCompile(`\(\d{4}\)`)

func (s *subscene) cleanSearchTerm(search string) string {
	search = cleanRegexp.ReplaceAllString(search, "")
	return strings.TrimSpace(search)
}

// FindMediaURL retrieves the subscene.com URL for the given media item
func (s *subscene) FindMediaURL(media types.Media, retries int) ([]searchResult, error) {
	url, err := url.Parse("https://subscene.com/subtitles/title")

	if err != nil {
		return nil, err
	}

	search := s.searchTerm(media)

	if len(search) == 0 {
		return nil, fmt.Errorf("unable to search for media: %s", media)
	}

	log.WithField("query", search).Debug("Searching subscene.com")

	query := url.Query()
	query.Add("q", s.cleanSearchTerm(search))
	url.RawQuery = query.Encode()

	lockSubscene()
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("subscene.com returned status code %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, err
	}

	var results []searchResult

	doc.Find("div.search-result div.title a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url, ok := s.Attr("href")
		if len(title) > 0 && ok {
			results = append(results, searchResult{
				Title: title,
				Path:  url,
			})
		}
	})

	if len(results) == 0 {
		return nil, errors.New("no media found on subscene.com")
	}

	results = resultList(results).Filter(media)

	return results, nil
}

type resultList []searchResult

// Best returns the search result from the list with the highest probability
// of fitting the target media. Is achieves this with the wagner fsicher
// string similiarity measure
func (r resultList) Best(m types.Media) (*searchResult, error) {
	var keyword string
	if movie, ok := m.TypeMovie(); ok {
		keyword = fmt.Sprintf("%s (%v)", movie.MovieName(), movie.Year())
	} else if episode, ok := m.TypeEpisode(); ok {
		season := parse.PhoneticNumber(episode.Season())
		keyword = fmt.Sprintf("%s - %s Season", episode.TVShow(), season)
	}

	if keyword == "" {
		return nil, errors.New("subscene unknown media")
	}

	var result *searchResult
	min := math.MaxInt32
	for _, e := range r {
		score := smetrics.WagnerFischer(e.Title, keyword, 1, 1, 2)
		if score < min {
			min = score
			result = &e
		}
	}

	if result == nil {
		return nil, errors.New("could not find best subtitle from subscene.com")
	}

	return result, nil
}

// Filter rules out search results which is guaranteed not to match the target media
func (r resultList) Filter(m types.Media) resultList {
	var include string
	if mov, ok := m.TypeMovie(); ok {
		include = strconv.Itoa(mov.Year())
	} else if eps, ok := m.TypeEpisode(); ok {
		season := parse.PhoneticNumber(eps.Season())
		include = fmt.Sprintf("%s Season", season)
	}
	var result []searchResult
	include = strings.ToLower(include)
	for _, e := range r {
		if strings.Contains(strings.ToLower(e.Title), include) {
			result = append(result, e)
		}
	}
	return result
}

// SearchSubtitles searches subscene.com for subtitles
func (s *subscene) SearchSubtitles(local types.LocalMedia) (subs []types.OnlineSubtitle, err error) {
	search, err := s.FindMediaURL(local, 3)

	if err != nil {
		return
	}

	best, err := resultList(search).Best(local)

	if err != nil {
		return nil, err
	}

	url, err := url.Parse(subsceneHost)

	if err != nil {
		return nil, err
	}

	url.Path = best.Path

	lockSubscene()
	doc, err := goquery.NewDocument(url.String())

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
		return nil, errors.New("no srt file found in zip")
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

func newRarReader(file *os.File) (*rarReader, error) {
	r, err := rardecode.NewReader(file, "")

	if err != nil {
		return nil, err
	}

	var found bool
	h, err := r.Next()
	for err != io.EOF {
		if filepath.Ext(h.Name) == ".srt" {
			found = true
			break
		}
		h, err = r.Next()
	}

	if !found {
		return nil, errors.New("no subtitle found in rar archive")
	}

	return &rarReader{Reader: r, file: file}, nil
}

type rarReader struct {
	*rardecode.Reader
	file *os.File
}

func (r *rarReader) Close() error {
	r.file.Close()
	if err := os.Remove(r.file.Name()); err != nil {
		log.WithError(err).Error("Could not cleaup temporary zip file")
		return err
	}
	return nil
}

func newSrtReader(file *os.File) (*srtReader, error) {
	return &srtReader{file}, nil
}

type srtReader struct {
	*os.File
}

func (s *srtReader) Close() error {
	s.File.Close()
	if err := os.Remove(s.File.Name()); err != nil {
		log.WithError(err).Error("Could not cleaup temporary zip file")
		return err
	}
	return nil
}

type subsceneURL string

func (uri subsceneURL) Link() string {
	return string(uri)
}

func (uri subsceneURL) Download() (r io.ReadCloser, err error) {
	var file *os.File

	defer func() {
		if err != nil && file != nil {
			os.Remove(file.Name())
		}
	}()

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

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))

	if err != nil {
		return nil, err
	}

	filename, ok := params["filename"]

	if !ok {
		return nil, errors.New("unable to retrieve file format from subscene.com")
	}

	ext := filepath.Ext(filename)

	file, err = ioutil.TempFile("", "supper")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"uri":    string(uri),
		"format": strings.TrimPrefix(ext, "."),
	}).Debug("Downloading from subscene.com")

	if ext == ".zip" {
		return newZipReader(file)
	} else if ext == ".rar" {
		return newRarReader(file)
	} else if ext == ".srt" {
		return newSrtReader(file)
	}
	return nil, fmt.Errorf("unknown subtitle format %s from subscene.com", ext)
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
