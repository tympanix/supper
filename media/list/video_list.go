package list

import (
	"github.com/fatih/set"
	"github.com/tympanix/supper/types"
)

// Video is a list of video media
type Video struct {
	video []types.Video
}

// NewVideo creates a new list of video media
func NewVideo(video ...types.Video) *Video {
	return &Video{
		video: video,
	}
}

// List returns the slice representation of the video list
func (l *Video) List() []types.Video {
	return l.video
}

// Len returns the length of the video list
func (l *Video) Len() int {
	return len(l.video)
}

// FilterMissingSubs returns a filtered list of video media which does not
// satisfy one or more of the subtitle languages in the input set. A language
// is satisfied if a subtitle with that language can be found on disk relative
// to the location of the video media
func (l *Video) FilterMissingSubs(lang set.Interface) (types.VideoList, error) {
	media := make([]types.Video, 0)
	for _, m := range l.List() {
		extsubs, err := m.ExistingSubtitles()
		if err != nil {
			return nil, err
		}
		missing := set.Difference(lang, extsubs.LanguageSet())
		if missing.Size() > 0 {
			media = append(media, m)
		}
	}
	return NewVideo(media...), nil
}
