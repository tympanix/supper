package api

import (
	"encoding/json"
	"errors"
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
	var media jsonFolder
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&media); err != nil {
		return Error(err, http.StatusBadRequest)
	}
	_, err := media.getPath(a)
	if err != nil {
		return Error(err, http.StatusBadRequest)
	}
	return nil
}
