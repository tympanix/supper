package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *API) subtitle(w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method == "POST" {
		return a.saveSubtitle(w, r)
	} else {
		return Error(errors.New("Method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (a *API) saveSubtitle(w http.ResponseWriter, r *http.Request) interface{} {
	var folder jsonFolder
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&folder); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	path, err := folder.getPath(a)
	if err != nil {
		fmt.Println("Oh noes 3")
		return Error(err, http.StatusBadRequest)
	}
	media, err := a.FindMedia(path)
	if err != nil {
		fmt.Println("Oh noes 2")
		return Error(err, http.StatusBadRequest)
	}
	err = a.DownloadSubtitles(media, a.Languages(), ioutil.Discard)
	if err != nil {
		fmt.Println("Oh noes 1")
		return Error(err, http.StatusBadRequest)
	}
	files, err := a.fileList(folder)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	return files
}
