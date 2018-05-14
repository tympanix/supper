package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// DownloadSubtitles downloads subtitles for a whole list of mediafiles for every
// langauge in the language set
func (a *Application) DownloadSubtitles(media types.LocalMediaList, lang set.Interface) (int, error) {
	numsubs := 0

	if media == nil {
		return -1, errors.New("no media supplied for subtitles")
	}

	if lang == nil {
		return -1, errors.New("no languages supplied for subtitles")
	}

	video := media.FilterVideo()

	video, err := video.FilterMissingSubs(lang)

	if err != nil {
		return -1, nil
	}

	// Iterate all media files in the list
	for i, item := range video.List() {
		ctx := log.WithFields(log.Fields{
			"media": item,
			"item":  fmt.Sprintf("%v/%v", i+1, video.Len()),
		})

		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return -1, err
		}

		missingLangs := set.Difference(lang, cursubs.LanguageSet())

		if missingLangs.Size() == 0 {
			continue
		}

		subs := list.RatedSubtitles(item)

		if !a.Config().Dry() {
			var search []types.OnlineSubtitle
			search, err = a.SearchSubtitles(item)
			if err != nil {
				ctx.WithError(err).Error("Subtitle failed")
				if a.Config().Strict() {
					return -1, err
				}
				continue
			}
			for _, s := range search {
				subs.Add(s)
			}
		}

		subs = subs.HearingImpaired(a.Config().Dry())

		// Download subtitle for each language
		for _, l := range missingLangs.List() {
			ctx = ctx.WithField("lang", display.English.Languages().Name(l))

			if a.Config().Delay() > 0 {
				time.Sleep(a.Config().Delay())
			}

			l, ok := l.(language.Tag)

			if !ok {
				return -1, err
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !a.Config().Dry() {
				ctx.Warn("No subtitle available")
				continue
			}

			var best float32

			if !a.Config().Dry() {
				var sub types.Subtitle
				sub, best = langsubs.Best()
				if best < (float32(a.Config().Score()) / 100.0) {
					ctx.Warnf("Score too low %.0f%%", best*100.0)
					continue
				}
				onl, ok := sub.(types.OnlineSubtitle)
				if !ok {
					ctx.Fatal("Subtitle could not be cast to online subtitle")
				}
				srt, err := onl.Download()
				if err != nil {
					ctx.WithError(err).Error("Could not download subtitle")
					if a.Config().Strict() {
						os.Exit(1)
					}
				}
				defer srt.Close()
				saved, err := item.SaveSubtitle(srt, onl.Language())
				if err != nil {
					ctx.WithError(err).Error("Subtitle error")
					if a.Config().Strict() {
						os.Exit(1)
					}
					continue
				}

				var strscore string
				if best == 0.0 {
					strscore = "N/A"
				} else {
					strscore = fmt.Sprintf("%.0f%%", best*100.0)
				}

				ctx.WithField("score", strscore).Info("Subtitle downloaded")

				for _, plugin := range a.Config().Plugins() {
					err := plugin.Run(saved)
					if err != nil {
						ctx.WithField("plugin", plugin.Name()).Error("Plugin failed")
						if a.Config().Strict() {
							os.Exit(1)
						}
					} else {
						ctx.WithField("plugin", plugin.Name()).Info("Plugin finished")
					}
				}
				numsubs++
			} else {
				numsubs++
				ctx.WithField("reason", "dry-run").Info("Skip download")
			}
		}
	}
	return numsubs, nil
}
