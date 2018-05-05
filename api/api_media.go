package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

const (
	typeMovie = "movie"
	typeShow  = "show"
)

type jsonFolder struct {
	os.FileInfo `json:"-"`
	Type        string `json:"type"`
	Folder      string `json:"folder"`
	Name        string `json:"name"`
}

type folderList []*jsonFolder

func (l folderList) Len() int {
	return len(l)
}

func (l folderList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l folderList) Less(i, j int) bool {
	return l[i].ModTime().After(l[j].ModTime())
}

func (f jsonFolder) getPath(a types.App) (path string, err error) {
	var root string
	if f.Type == typeMovie {
		root = a.Config().Movies().Directory()
	} else if f.Type == typeShow {
		root = a.Config().TVShows().Directory()
	} else {
		err = errors.New("unknown media format")
		return
	}
	path = filepath.Join(root, f.Folder)
	if filepath.Dir(path) != filepath.Clean(root) {
		return "", errors.New("illegal folder path")
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
					FileInfo: file,
					Type:     t,
					Folder:   file.Name(),
					Name:     parse.CleanName(file.Name()),
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
		err := errors.New("method not allowed")
		return Error(err, http.StatusMethodNotAllowed)
	}
}

func (a *API) allMedia(w http.ResponseWriter, r *http.Request) interface{} {
	media := make([]*jsonFolder, 0)

	for mtype, path := range map[string]string{
		typeMovie: a.Config().Movies().Directory(),
		typeShow:  a.Config().TVShows().Directory(),
	} {
		if path != "" {
			list, err := findMediaFolders(mtype, path)
			if err != nil {
				return Error(err, 500)
			}
			media = append(media, list...)
		}
	}

	sort.Sort(folderList(media))
	return media
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
	video := list.FilterVideo()
	medialist := make([]interface{}, 0)
	for _, m := range video.List() {
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
		relpath, err := filepath.Rel(path, m.Path())
		if err != nil {
			return nil, errors.New("interval path error for media file")
		}
		medialist = append(medialist, struct {
			Type  string             `json:"type"`
			Name  string             `json:"filename"`
			Path  string             `json:"filepath"`
			Media types.Media        `json:"media"`
			Subs  types.SubtitleList `json:"subtitles"`
		}{
			mtype,
			m.Name(),
			relpath,
			m,
			subs,
		})
	}
	return medialist, nil
}
