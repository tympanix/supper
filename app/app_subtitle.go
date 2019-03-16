package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/set"
	"github.com/tympanix/supper/app/logutil"
	"github.com/tympanix/supper/app/notify"
	"github.com/tympanix/supper/media/list"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// DownloadSubtitles downloads subtitles for a whole list of mediafiles for every
// langauge in the language set
func (a *Application) DownloadSubtitles(input types.LocalMediaList, lang set.Interface, c chan<- *notify.Entry) ([]types.LocalSubtitle, error) {
	var result []types.LocalSubtitle

	if input == nil {
		return nil, errors.New("no media supplied for subtitles")
	}

	if lang == nil {
		return nil, errors.New("no languages supplied for subtitles")
	}

	video := input.FilterVideo()

	if video.Len() == 0 {
		return nil, errors.New("no video media found in path")
	}

	video, err := video.FilterMissingSubs(lang)

	if err != nil {
		return nil, err
	}

	// Iterate all media files in the list
	for i, item := range video.List() {
		ctx := notify.WithFields(notify.Fields{
			"media": item,
			"item":  fmt.Sprintf("%v/%v", i+1, video.Len()),
		})

		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return nil, err
		}

		missingLangs := set.Difference(lang, cursubs.LanguageSet())

		if missingLangs.Size() == 0 {
			continue
		}

		var subs = list.Subtitles()

		if !a.Config().Dry() {
			var search []types.OnlineSubtitle
			search, err = a.SearchSubtitles(item)
			if err != nil {
				c <- ctx.WithError(err).Error("Subtitle failed")
				if a.Config().Strict() {
					return nil, err
				}
				continue
			}
			subs, err = list.NewSubtitlesFromInterface(search)
			if err != nil {
				ctx.WithError(err).Fatal("Subtitle error")
			}
		}

		subs = subs.HearingImpaired(a.Config().Impaired())

		// Download subtitle for each language
		for _, v := range missingLangs.List() {
			l, ok := v.(language.Tag)
			if !ok {
				return nil, logutil.Errorf("unknown language %v", v)
			}

			ctx = ctx.WithField("lang", display.English.Languages().Name(l))

			if a.Config().Delay() > 0 {
				time.Sleep(a.Config().Delay())
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !a.Config().Dry() {
				c <- ctx.Warn("No subtitle available")
				continue
			}

			if !a.Config().Dry() {
				rated := langsubs.RateByMedia(item, a.Config().Evaluator())
				sub, err := a.downloadBestSubtitle(ctx, item, rated, 3, c)
				if err != nil {
					if a.Config().Strict() {
						return nil, err
					}
					c <- ctx.WithError(err).Error("Could not download subtitle")
					continue
				}
				result = append(result, sub)
			} else {
				c <- ctx.WithField("reason", "dry-run").Info("Skip download")
			}
		}
	}
	return result, nil
}

func (a *Application) downloadBestSubtitle(ctx notify.Context, m types.Video, l types.RatedSubtitleList, retries int, c chan<- *notify.Entry) (types.LocalSubtitle, error) {
	if l.Len() == 0 {
		return nil, ctx.Warn("No subtitles satisfied media")
	}
	sub := l.Best()
	if sub.Score() < (float32(a.Config().Score()) / 100.0) {
		return nil, ctx.Warn("Score too low %.0f%%", sub.Score()*100.0)
	}
	onl, ok := sub.Subtitle().(types.OnlineSubtitle)
	if !ok {
		ctx.Fatal("Subtitle could not be cast to online subtitle")
	}
	srt, err := onl.Download()
	if err != nil && retries > 0 {
		p := list.RatedSubtitles(l.List()[1:])
		if p.Len() <= 0 {
			return nil, ctx.Error(err.Error())
		}
		c <- ctx.WithError(err).Debug("Retrying subtitle")
		return a.downloadBestSubtitle(ctx, m, p, retries-1, c)
	}
	if err != nil {
		return nil, ctx.Error(err.Error())
	}
	defer srt.Close()
	saved, err := m.SaveSubtitle(srt, onl.Language())
	if err != nil {
		return nil, ctx.Error(err.Error())
	}

	score := fmt.Sprintf("%.0f%%", sub.Score()*100.0)
	c <- ctx.WithField("score", score).WithExtra("sub", saved).Info("Subtitle downloaded")

	if err := a.execPluginsOnSubtitle(ctx, saved, c); err != nil {
		return nil, err
	}
	return saved, nil
}

func (a *Application) execPluginsOnSubtitle(ctx notify.Context, s types.LocalSubtitle, c chan<- *notify.Entry) error {
	for _, plugin := range a.Config().Plugins() {
		ctx = ctx.WithField("plugin", plugin.Name())
		if err := plugin.Run(s); err != nil {
			c <- ctx.Error("Plugin failed")
			return err
		}
		c <- ctx.Info("Plugin finished")
	}
	return nil
}
