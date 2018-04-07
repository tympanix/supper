package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// DownloadSubtitles downloads subtitles for a whole list of mediafiles for every
// langauge in the language set
func (a *Application) DownloadSubtitles(media types.LocalMediaList, lang set.Interface) (int, error) {
	numsubs := 0

	// Iterate all media files found in each path
	for i, item := range media.List() {
		ctx := log.WithFields(log.Fields{
			"media": item,
			"item":  fmt.Sprintf("%v/%v", i+1, media.Len()),
		})

		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return -1, err
		}

		var missingLangs set.Interface
		if !a.Config().Dry() {
			missingLangs = set.Difference(lang, cursubs.LanguageSet())

			if missingLangs.Size() == 0 {
				continue
			}
		} else {
			missingLangs = lang
		}

		subs := list.RatedSubtitles(item)

		if !a.Config().Dry() {
			search, err := a.SearchSubtitles(item)
			if err != nil {
				return -1, err
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
				ctx.Warn("Subtitle not available")
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
				saved, err := item.SaveSubtitle(onl, onl.Language())
				if err != nil {
					ctx.WithError(err).Error("Subtitle error")
					if a.Config().Strict() {
						os.Exit(1)
					}
					continue
				}
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
			}

			var strscore string
			if best == 0.0 {
				strscore = "N/A"
			} else {
				strscore = fmt.Sprintf("%.0f%%", best*100.0)
			}

			ctx.WithField("score", strscore).Info("Subtitle downloaded")
		}
	}
	return numsubs, nil
}
