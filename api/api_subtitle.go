package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/tympanix/supper/types"
)

var busyFolders = new(sync.Map)

type jsonMedia struct {
	jsonFolder
	Filename string `json:"filename"`
}

func (j jsonMedia) getPath(a types.App) (path string, err error) {
	folder, err := j.jsonFolder.getPath(a)
	if err != nil {
		return
	}
	path = filepath.Join(folder, j.Filename)
	if !filepath.HasPrefix(path, folder) {
		return "", errors.New("Illegal media path")
	}
	return
}

func (a *API) subtitleRouter(mux *mux.Router) {
	mux.Queries("action", "download").Methods("POST").
		Handler(apiHandler(a.downloadSubtitles))
	mux.Queries("action", "list").Methods("POST").
		Handler(apiHandler(a.getSubtitles))
}

func (a *API) subtitles(w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method == "POST" {
		return a.downloadSubtitles(w, r)
	} else {
		return Error(errors.New("Method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (a *API) getSubtitles(w http.ResponseWriter, r *http.Request) interface{} {
	var mediaFile jsonMedia
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mediaFile); err != nil {
		return err
	}
	path, err := mediaFile.getPath(a)
	if err != nil {
		return err
	}
	media, err := a.FindMedia(path)
	if err != nil {
		return err
	}
	if media.Len() > 1 {
		return errors.New("More than one media file found")
	}
	if media.Len() < 1 {
		return errors.New("No media found")
	}
	subs, err := a.SearchSubtitles(media.List()[0])
	if err != nil {
		return err
	}
	return subs
}

func (a *API) downloadSubtitles(w http.ResponseWriter, r *http.Request) interface{} {
	var folder jsonFolder
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&folder); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	path, err := folder.getPath(a)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	if _, busy := busyFolders.LoadOrStore(path, true); busy {
		return Error(errors.New("Folder is busy"), http.StatusTooManyRequests)
	}
	defer busyFolders.Delete(path)
	media, err := a.FindMedia(path)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	err = a.DownloadSubtitles(media, a.Languages(), ioutil.Discard)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	files, err := a.fileList(folder)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	return files
}
