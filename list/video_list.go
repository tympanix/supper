package list

import (
	"github.com/fatih/set"
	"github.com/tympanix/supper/types"
)

// Video is a list of video media
type Video struct {
	video []types.Video
}

func NewVideo(video ...types.Video) *Video {
	return &Video{
		video: video,
	}
}

func (l *Video) List() []types.Video {
	return l.video
}

func (l *Video) Len() int {
	return len(l.video)
}

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
