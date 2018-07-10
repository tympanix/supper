package api

import (
	"errors"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type jsonConfig struct {
	Lang      []jsonLang `json:"languages"`
	Proxypath string     `json:"proxypath"`
}

type jsonLang struct {
	Code language.Tag `json:"code"`
	Lang string       `json:"language"`
}

func (a *API) config(w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method == "GET" {
		langs := make([]jsonLang, 0)
		for _, l := range a.Config().Languages().List() {
			if tag, ok := l.(language.Tag); ok {
				langs = append(langs, jsonLang{
					tag,
					display.English.Languages().Name(tag),
				})
			}
		}
		return jsonConfig{
			langs,
			a.Config().ProxyPath(),
		}
	}
	err := errors.New("Method not allowed")
	return NewError(err, http.StatusMethodNotAllowed)
}
