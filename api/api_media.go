package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/tympanix/supper/types"
)

const (
	typeMovie = "movie"
	typeShow  = "show"
)

type mediaFolder struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func findMediaFolders(t string, paths ...string) ([]*mediaFolder, error) {
	media := make([]*mediaFolder, 0)

	for _, path := range paths {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				media = append(media, &mediaFolder{
					Type: t,
					Name: file.Name(),
				})
			}
		}
	}

	return media, nil
}

func (a *API) media(w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method == "GET" {
		return a.allMedia(w, r)
	} else if r.Method == "POST" {
		return a.detailsMedia(w, r)
	} else {
		err := errors.New("Method not allowed")
		return Error(err, http.StatusMethodNotAllowed)
	}
}

func (a *API) allMedia(w http.ResponseWriter, r *http.Request) interface{} {
	movies, err := findMediaFolders(typeMovie, a.Context().String("movies"))
	if err != nil {
		return Error(err, 500)
	}
	tvshows, err := findMediaFolders(typeShow, a.Context().String("shows"))
	if err != nil {
		return Error(err, 500)
	}
	return append(movies, tvshows...)
}

func (a *API) detailsMedia(w http.ResponseWriter, r *http.Request) interface{} {
	var media mediaFolder
	var root string
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&media); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if media.Type == typeMovie {
		root = a.Context().String("movies")
	} else if media.Type == typeShow {
		root = a.Context().String("shows")
	} else {
		err := errors.New("Unknown media format")
		return Error(err, http.StatusBadRequest)
	}
	path := filepath.Join(root, media.Name)
	if filepath.Dir(path) != filepath.Clean(root) {
		err := errors.New("Illegal folder path")
		return Error(err, http.StatusBadRequest)
	}
	list, err := a.FindMedia(path)
	if err != nil {
		return Error(err, http.StatusInternalServerError)
	}
	medialist := make([]interface{}, 0)
	for _, m := range list.List() {
		subs, err := m.ExistingSubtitles()
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}
		var mtype string
		if _, ok := m.TypeEpisode(); ok {
			mtype = typeShow
		} else if _, ok := m.TypeMovie(); ok {
			mtype = typeMovie
		}
		medialist = append(medialist, struct {
			Type  string             `json:"type"`
			Media types.Media        `json:"media"`
			Subs  types.SubtitleList `json:"subtitles"`
		}{
			mtype,
			m,
			subs,
		})
	}
	return medialist
}
