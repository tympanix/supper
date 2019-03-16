package list

import (
	"encoding/json"
	"time"

	"github.com/tympanix/supper/types"
)

// LocalMedia is a list of locally stored media
type LocalMedia struct {
	media []types.LocalMedia
}

// NewLocalMedia return a new local media list from its arguments
func NewLocalMedia(media ...types.LocalMedia) *LocalMedia {
	return &LocalMedia{
		media,
	}
}

// Add adds new local media to the list
func (l *LocalMedia) Add(m types.LocalMedia) {
	l.media = append(l.media, m)
}

// Len returns the number of media in the list
func (l *LocalMedia) Len() int {
	return len(l.media)
}

// List returns the list of localmedia as a plain slice
func (l *LocalMedia) List() []types.LocalMedia {
	return l.media
}

// Filter return the list of local media which satisfies some predicate
func (l *LocalMedia) Filter(p types.MediaFilter) types.LocalMediaList {
	filtered := make([]types.LocalMedia, 0)
	for _, media := range l.List() {
		if p(media) {
			filtered = append(filtered, media)
		}
	}
	return NewLocalMedia(filtered...)
}

// FilterModified returns only media which has been modified since some duration
func (l *LocalMedia) FilterModified(d time.Duration) types.LocalMediaList {
	t := time.Now().Local().Add(-1 * d)
	media := make([]types.LocalMedia, 0)
	for _, m := range l.List() {
		if m.ModTime().After(t) {
			media = append(media, m)
		}
	}
	return NewLocalMedia(media...)
}

// FilterVideo returns only media which is of type video (e.g. not subtites)
func (l *LocalMedia) FilterVideo() types.VideoList {
	video := make([]types.Video, 0)
	for _, media := range l.List() {
		if v, ok := media.(types.Video); ok {
			video = append(video, v)
		}
	}
	return NewVideo(video...)
}

// FilterMovies return only media which is of type movie
func (l *LocalMedia) FilterMovies() types.LocalMediaList {
	movies := make([]types.LocalMedia, 0)
	for _, media := range l.List() {
		if _, ok := media.TypeMovie(); ok {
			movies = append(movies, media)
		}
	}
	return NewLocalMedia(movies...)
}

// FilterEpisodes returns only media which is of type episode
func (l *LocalMedia) FilterEpisodes() types.LocalMediaList {
	episodes := make([]types.LocalMedia, 0)
	for _, media := range l.List() {
		if _, ok := media.TypeEpisode(); ok {
			episodes = append(episodes, media)
		}
	}
	return NewLocalMedia(episodes...)
}

// FilterSubtitles returns only media which is of type subtitles
func (l *LocalMedia) FilterSubtitles() types.LocalMediaList {
	subtitles := make([]types.LocalMedia, 0)
	for _, media := range l.List() {
		if _, ok := media.TypeSubtitle(); ok {
			subtitles = append(subtitles, media)
		}
	}
	return NewLocalMedia(subtitles...)
}

// MarshalJSON returns a JSON representation of the media list
func (l *LocalMedia) MarshalJSON() (b []byte, err error) {
	return json.Marshal(l.List())
}
