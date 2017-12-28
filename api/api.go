package api

import (
	"encoding/json"
	"net/http"

	"github.com/tympanix/supper/types"
)

type API struct {
	types.App
	*http.ServeMux
}

type APIError interface {
	error
	Status() int
}

func New(app types.App) http.Handler {
	api := &API{
		app,
		http.NewServeMux(),
	}

	api.Handle("/media", apiHandler(api.media))

	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.ServeMux.ServeHTTP(w, r)
}

type apiError struct {
	error `json:"-"`
	Code  int `json:"-"`
}

func (e *apiError) Status() int {
	return e.Code
}

type apiHandler func(http.ResponseWriter, *http.Request) interface{}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		js, err := json.Marshal(e)
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
		http.Error(w, "No found", http.StatusNotFound)
	}
}

func Error(err error, status int) APIError {
	return &apiError{
		err,
		status,
	}
}
