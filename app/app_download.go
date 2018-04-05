package application

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/apex/log"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// DownloadSubtitles downloads subtitles for a whole list of mediafiles for every
// langauge in the language set. Any output is written to the ourpur writer
func (a *Application) DownloadSubtitles(media types.LocalMediaList, lang set.Interface, out io.Writer) (int, error) {
	numsubs := 0

	dry := a.Context().GlobalBool("dry")
	hi := a.Context().GlobalBool("impaired")
	score := a.Context().GlobalInt("score")
	force := a.Context().GlobalBool("force")
	delay, err := parse.Duration(a.Context().GlobalString("delay"))

	if err != nil {
		return 0, errors.New("could not parse delay time format")
	}

	// Iterate all media files found in each path
	for i, item := range media.List() {
		ctx := log.WithFields(log.Fields{
			"media": item,
			"item":  fmt.Sprintf("%v/%v", i+1, media.Len()),
		})

		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return -1, cli.NewExitError(err, 2)
		}

		var missingLangs set.Interface
		if !force {
			missingLangs = set.Difference(lang, cursubs.LanguageSet())

			if missingLangs.Size() == 0 {
				continue
			}
		} else {
			missingLangs = lang
		}

		subs := list.RatedSubtitles(item)

		if !dry {
			search, err := a.SearchSubtitles(item)
			if err != nil {
				return -1, cli.NewExitError(err, 2)
			}
			for _, s := range search {
				subs.Add(s)
			}
		}

		subs = subs.HearingImpaired(hi)

		// Download subtitle for each language
		for _, l := range missingLangs.List() {
			ctx = ctx.WithField("lang", display.English.Languages().Name(l))

			if delay > 0 {
				time.Sleep(delay)
			}

			l, ok := l.(language.Tag)

			if !ok {
				return -1, cli.NewExitError(err, 3)
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !dry {
				ctx.Warn("Subtitle not found")
				continue
			}

			if !dry {
				sub, best := langsubs.Best()
				if best < (float32(score) / 100.0) {
					ctx.Warnf("Score too low %.0f%%\n", best*100.0)
					continue
				}
				onl, ok := sub.(types.OnlineSubtitle)
				if !ok {
					ctx.Fatal("Subtitle could not be cast to online subtitle")
				}
				saved, err := item.SaveSubtitle(onl, onl.Language())
				if err != nil {
					ctx.Error(err.Error())
					continue
				}
				for _, plugin := range a.Plugins() {
					err := plugin.Run(saved)
					if err != nil {
						ctx.Errorf("Plugin failed: %s\n", plugin.Name())
					} else {
						ctx.Infof("Plugin finished: %s\n", plugin.Name())
					}
				}
				numsubs++
			}

			ctx.Info("Subtitle downloaded")
		}
	}
	return numsubs, nil
}
