package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fatih/set"
	"github.com/gorilla/mux"
	"github.com/tympanix/supper/notify"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
)

// API exposes endpoints for the webapp to perform HTTP RESTFull actions
type API struct {
	types.App
	*Hub
	*mux.Router
}

// Error is an error occurring in an API endpoint
type Error interface {
	error
	Status() int
}

// New creates a new API handler
func New(app types.App) http.Handler {
	api := &API{
		App:    app,
		Hub:    newHub(),
		Router: mux.NewRouter(),
	}

	api.Handle("/media", apiHandler(api.media))
	api.Handle("/config", apiHandler(api.config))
	api.HandleFunc("/ws", api.serveWebsocket)
	apiSubs := api.PathPrefix("/subtitles").Subrouter()
	api.subtitleRouter(apiSubs)

	go api.Hub.run()

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

func (a *API) queryLang(r *http.Request) (set.Interface, error) {
	v := r.URL.Query()
	lang := v.Get("lang")
	if lang == "" {
		return a.Config().Languages(), nil
	}
	l := language.Make(lang)
	if l == language.Und {
		return set.New(), errors.New("unknown language")
	}
	return set.New(l), nil
}

type apiHandler func(http.ResponseWriter, *http.Request) interface{}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		if err, ok := e.(error); ok {
			e = fn.handleError(w, err)
		}
		js, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (fn apiHandler) handleError(w http.ResponseWriter, err error) error {
	var apiError Error
	if e, ok := err.(Error); ok {
		apiError = e
	} else {
		apiError = NewError(err, http.StatusBadRequest)
	}
	w.WriteHeader(apiError.Status())
	return apiError
}

// NewError returns a new error
func NewError(err error, status int) Error {
	return &apiError{
		err,
		status,
	}
}

func (a *API) sendToWebsocket() chan<- *notify.Entry {
	c := make(chan *notify.Entry)
	go func() {
		for v := range c {
			fmt.Println(v)
		}
	}()
	return c
}
