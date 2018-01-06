package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tympanix/supper/types"
)

type API struct {
	types.App
	*mux.Router
}

type APIError interface {
	error
	Status() int
}

func New(app types.App) http.Handler {
	api := &API{
		app,
		mux.NewRouter(),
	}

	api.Handle("/media", apiHandler(api.media))
	api.Handle("/config", apiHandler(api.config))
	apiSubs := api.PathPrefix("/subtitles").Subrouter()
	api.subtitleRouter(apiSubs)

	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

type apiError struct {
	error `json:"-"`
	Code  int `json:"-"`
}

func (e *apiError) Status() int {
	return e.Code
}

func (e *apiError) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: e.Error(),
	})
}

type apiHandler func(http.ResponseWriter, *http.Request) interface{}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		if err, ok := e.(error); ok {
			e = Error(err, http.StatusBadRequest)
		}
		js, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err, ok := e.(APIError); ok {
			w.WriteHeader(err.Status())
		}
		w.Write(js)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func Error(err error, status int) APIError {
	return &apiError{
		err,
		status,
	}
}
