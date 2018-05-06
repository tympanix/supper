package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/types"
)

const thetvdbHost = "https://api.thetvdb.com/"

var thetvdbClient = NewAPIClient("TheTVDB", 35)

var thetvdbCache = make(map[string]types.Media)

// TheTVDB is a scraper for thetvdb.com
func TheTVDB(key string) types.Scraper {
	return &thetvdb{
		client: thetvdbClient,
		key:    key,
	}
}

type thetvdb struct {
	client *APIClient
	key    string
	token  string
}

func (t *thetvdb) Scrape(m types.Media) (src types.Media, err error) {
	if t.key == "" {
		return nil, errors.New("thetvdb: missing API key")
	}
	if m == nil {
		return nil, errors.New("thetvdb: can't scrape nil media")
	}
	if cache, ok := thetvdbCache[m.Identity()]; ok {
		return cache, nil
	}
	if e, ok := m.TypeEpisode(); ok {
		src, err = t.searchTV(e)
	} else if sub, ok := m.TypeSubtitle(); ok {
		src, err = t.Scrape(sub.ForMedia())
	} else {
		return nil, mediaNotSupported("thetvdb")
	}
	if err == nil {
		thetvdbCache[m.Identity()] = src
	}
	return
}

func (t *thetvdb) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	if t.token == "" {
		return nil, errors.New("thetvdb: not authenticated")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))

	return t.client.Do(req)
}

func (t *thetvdb) url(p string) (*url.URL, error) {
	url, err := url.Parse(thetvdbHost)

	url.Path = path.Join(url.Path, p)

	if err != nil {
		return nil, err
	}

	return url, nil
}

func (t *thetvdb) authenticate() error {
	url, err := t.url("/login")

	if err != nil {
		return err
	}

	post := struct {
		APIKey string `json:"apikey"`
	}{
		APIKey: t.key,
	}

	data, err := json.Marshal(&post)

	if err != nil {
		return err
	}

	resp, err := t.client.Post(url.String(), "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("tmdb could not authenticate %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	res := struct {
		Token string `json:"token"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	t.token = res.Token

	return nil
}

func (t *thetvdb) searchTV(e types.Episode) (types.Media, error) {
	if t.token == "" {
		if err := t.authenticate(); err != nil {
			return nil, err
		}
	}

	url, err := t.url("/search/series")

	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Set("name", e.TVShow())
	url.RawQuery = q.Encode()

	resp, err := t.Get(url.String())

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("thetvdb: api returned %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	type series struct {
		ID         int    `json:"id"`
		FirstAired string `json:"firstAired"`
		SeriesName string `json:"seriesName"`
	}

	seriesData := struct {
		Data []series `json:"data"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&seriesData); err != nil {
		return nil, err
	}

	if len(seriesData.Data) == 0 {
		return nil, errors.New("no media found on thetvdb")
	}

	url, err = t.url(fmt.Sprintf("/series/%v/episodes/query", seriesData.Data[0].ID))

	if err != nil {
		return nil, err
	}

	q = url.Query()
	q.Set("airedSeason", strconv.Itoa(e.Season()))
	q.Set("airedEpisode", strconv.Itoa(e.Episode()))
	url.RawQuery = q.Encode()

	resp, err = t.Get(url.String())

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("thetvdb: api returned %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	type episode struct {
		Season      int    `json:"airedSeason"`
		Episode     int    `json:"airedEpisodeNumber"`
		EpisodeName string `json:"episodeName"`
	}

	episodeData := struct {
		Data []episode `json:"data"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&episodeData); err != nil {
		return nil, err
	}

	scraped := media.Episode{
		NameX:        seriesData.Data[0].SeriesName,
		EpisodeNameX: episodeData.Data[0].EpisodeName,
		EpisodeX:     episodeData.Data[0].Episode,
		SeasonX:      episodeData.Data[0].Season,
	}

	return &scraped, nil
}
