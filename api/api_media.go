package api

import (
	"net/http"
)

func (a *API) media(w http.ResponseWriter, r *http.Request) interface{} {
	media, err := a.FindMedia(a.Args()...)
	if err != nil {
		return Error(err, 500)
	}
	return media.List()
}
