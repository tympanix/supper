package api

import (
	"io/ioutil"
	"net/http"
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
	movies, err := findMediaFolders("movie", a.Context().String("movies"))
	if err != nil {
		return Error(err, 500)
	}
	tvshows, err := findMediaFolders("show", a.Context().String("shows"))
	if err != nil {
		return Error(err, 500)
	}
	return append(movies, tvshows...)
}
