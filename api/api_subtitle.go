package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

var busyFolders = new(sync.Map)

type jsonMedia struct {
	jsonFolder
	Filepath string `json:"filepath"`
}

func (m jsonMedia) getPath(a types.App) (path string, err error) {
	folder, err := m.jsonFolder.getPath(a)
	if err != nil {
		return
	}
	path = filepath.Join(folder, m.Filepath)
	if !filepath.HasPrefix(path, folder) {
		return "", errors.New("Illegal media path")
	}
	return
}

func (m jsonMedia) MediaFile(a types.App) (types.LocalMedia, error) {
	path, err := m.getPath(a)
	if err != nil {
		return nil, err
	}
	media, err := a.FindMedia(path)
	if err != nil {
		return nil, err
	}
	if media.Len() != 1 {
		return nil, errors.New("No single media file found")
	}
	return media.List()[0], nil
}

type jsonSubtitle struct {
	jsonMedia
	URL  string `json:"link"`
	Lang string `json:"language"`
}

func (s jsonSubtitle) Link() string {
	return s.URL
}

type jsonRatedSubtitle struct {
	types.RatedSubtitle
}

func (r jsonRatedSubtitle) MarshalJSON() ([]byte, error) {
	hash := sha1.New()
	s := r.Subtitle()

	dl, ok := s.(types.OnlineSubtitle)
	if !ok {
		return nil, errors.New("Could not marshal subtitle which is not online")
	}

	info := []string{
		dl.Link(),
		s.ForMedia().Meta().Codec().String(),
		s.ForMedia().Meta().Group(),
		s.ForMedia().Meta().Quality().String(),
		s.ForMedia().Meta().Source().String(),
	}

	hash.Write([]byte(strings.Join(info, "")))
	hashval := hash.Sum(nil)
	infohash := make([]byte, hex.EncodedLen(len(hashval)))
	hex.Encode(infohash, hashval)

	return json.Marshal(struct {
		Hash  string       `json:"hash"`
		Lang  language.Tag `json:"language"`
		Link  string       `json:"link"`
		Score float32      `json:"score"`
		HI    bool         `json:"hi"`
		Media types.Media  `json:"media"`
	}{
		string(infohash),
		s.Language(),
		dl.Link(),
		r.Score(),
		s.HearingImpaired(),
		s.ForMedia(),
	})
}

type jsonSubtitleList []types.RatedSubtitle

func (l jsonSubtitleList) MarshalJSON() ([]byte, error) {
	var subs []jsonRatedSubtitle
	for _, s := range l {
		subs = append(subs, jsonRatedSubtitle{s})
	}
	return json.Marshal(subs)
}

func (a *API) subtitleRouter(mux *mux.Router) {
	mux.Queries("action", "download").Methods("POST").
		Handler(apiHandler(a.downloadSubtitles))
	mux.Queries("action", "list").Methods("POST").
		Handler(apiHandler(a.listSubtitles))
	mux.Queries("action", "single").Methods("POST").
		Handler(apiHandler(a.singleSubtitle))
}

func (a *API) singleSubtitle(w http.ResponseWriter, r *http.Request) interface{} {
	var mediaSubtitle jsonSubtitle
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mediaSubtitle); err != nil {
		return err
	}
	sub, err := a.ResolveSubtitle(mediaSubtitle)
	if err != nil {
		return err
	}
	media, err := mediaSubtitle.MediaFile(a)
	if err != nil {
		return err
	}
	tag := language.Make(mediaSubtitle.Lang)
	if tag == language.Und {
		return errors.New("unknown subtitle language")
	}
	video, ok := media.(types.Video)
	if !ok {
		return errors.New("media is not video")
	}
	srt, err := sub.Download()
	if err != nil {
		return err
	}
	defer srt.Close()
	_, err = video.SaveSubtitle(srt, tag)
	if err != nil {
		return err
	}
	return struct {
		Message string `json:"message"`
	}{
		"ok",
	}
}

func (a *API) listSubtitles(w http.ResponseWriter, r *http.Request) interface{} {
	var mediaFile jsonMedia
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mediaFile); err != nil {
		return err
	}
	item, err := mediaFile.MediaFile(a)
	if err != nil {
		return err
	}
	search, err := a.SearchSubtitles(item)
	if err != nil {
		return err
	}
	sublist, err := list.NewSubtitlesFromInterface(search)
	if err != nil {
		return err
	}
	rated := sublist.RateByMedia(item, a.Config().Evaluator())
	return jsonSubtitleList(rated.List())
}

func (a *API) downloadSubtitles(w http.ResponseWriter, r *http.Request) interface{} {
	langs, err := a.queryLang(r)
	if err != nil {
		return errors.New("unknown language for subtitle")
	}
	var folder jsonFolder
	dec := json.NewDecoder(r.Body)
	if err = dec.Decode(&folder); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	path, err := folder.getPath(a)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if _, busy := busyFolders.LoadOrStore(path, true); busy {
		return Error(errors.New("folder is busy"), http.StatusTooManyRequests)
	}
	defer busyFolders.Delete(path)
	media, err := a.FindMedia(path)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if media.Len() <= 0 {
		return Error(errors.New("no media found"), http.StatusBadRequest)
	}
	v, err := media.FilterVideo().FilterMissingSubs(langs)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if v.Len() <= 0 {
		return Error(errors.New("subtitle already satisfied"), http.StatusAccepted)
	}
	subs, err := a.DownloadSubtitles(media, langs)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if len(subs) <= 0 {
		return Error(errors.New("no subtitles found"), http.StatusBadRequest)
	}
	files, err := a.fileList(folder)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	return files
}
