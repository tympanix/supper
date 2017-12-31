package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

const (
	typeMovie = "movie"
	typeShow  = "show"
)

type jsonFolder struct {
	Type   string `json:"type"`
	Folder string `json:"folder"`
	Name   string `json:"name"`
}

func (f jsonFolder) getPath(a types.App) (path string, err error) {
	var root string
	if f.Type == typeMovie {
		root = a.Context().String("movies")
	} else if f.Type == typeShow {
		root = a.Context().String("shows")
	} else {
		err = errors.New("Unknown media format")
		return
	}
	path = filepath.Join(root, f.Folder)
	if filepath.Dir(path) != filepath.Clean(root) {
		return "", errors.New("Illegal folder path")
	}
	return
}

func findMediaFolders(t string, paths ...string) ([]*jsonFolder, error) {
	media := make([]*jsonFolder, 0)

	for _, path := range paths {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				media = append(media, &jsonFolder{
					Type:   t,
					Folder: file.Name(),
					Name:   parse.CleanName(file.Name()),
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
	var media jsonFolder
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&media); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	files, err := a.fileList(media)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	return files
}

func (a *API) fileList(folder jsonFolder) (interface{}, error) {
	path, err := folder.getPath(a)
	if err != nil {
		return nil, Error(err, http.StatusBadRequest)
	}
	list, err := a.FindMedia(path)
	if err != nil {
		return nil, Error(err, http.StatusInternalServerError)
	}
	medialist := make([]interface{}, 0)
	for _, m := range list.List() {
		subs, err := m.ExistingSubtitles()
		if err != nil {
			return nil, Error(err, http.StatusInternalServerError)
		}
		var mtype string
		if _, ok := m.TypeEpisode(); ok {
			mtype = typeShow
		} else if _, ok := m.TypeMovie(); ok {
			mtype = typeMovie
		}
		medialist = append(medialist, struct {
			Type  string             `json:"type"`
			Path  string             `json:"filepath"`
			Media types.Media        `json:"media"`
			Subs  types.SubtitleList `json:"subtitles"`
		}{
			mtype,
			m.Name(),
			m,
			subs,
		})
	}
	return medialist, nil
}
