package list

import (
	"time"

	"github.com/fatih/set"
	"github.com/tympanix/supper/types"
)

type LocalMedia struct {
	media []types.LocalMedia
}

func NewLocalMedia(media ...types.LocalMedia) *LocalMedia {
	return &LocalMedia{
		media,
	}
}

func (l *LocalMedia) Add(m types.LocalMedia) {
	l.media = append(l.media, m)
}

func (l *LocalMedia) Len() int {
	return len(l.media)
}

func (l *LocalMedia) List() []types.LocalMedia {
	return l.media
}

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

func (l *LocalMedia) FilterMissingSubs(lang set.Interface) (types.LocalMediaList, error) {
	media := make([]types.LocalMedia, 0)
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
	return NewLocalMedia(media...), nil
}
