package application

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)
var yellow = color.New(color.FgYellow)

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

		fmt.Fprintf(out, "(%v/%v) - %s\n", i+1, media.Len(), item)

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
			if delay > 0 {
				time.Sleep(delay)
			}

			l, ok := l.(language.Tag)

			if !ok {
				return -1, cli.NewExitError(err, 3)
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !dry {
				red.Fprintln(out, " - no subtitles found")
				continue
			}

			if !dry {
				sub, best := langsubs.Best()
				if best < (float32(score) / 100.0) {
					yellow.Fprintf(out, " - score too low %.0f%%\n", best*100.0)
					continue
				}
				onl, ok := sub.(types.OnlineSubtitle)
				if !ok {
					panic("subtitle could not be cast to online subtitle")
				}
				saved, err := item.SaveSubtitle(onl, onl.Language())
				if err != nil {
					red.Fprintln(out, err.Error())
					continue
				}
				for _, plugin := range a.Plugins() {
					err := plugin.Run(saved)
					if err != nil {
						red.Fprintf(out, " - Plugin failed: %s\n", plugin.Name())
					} else {
						green.Fprintf(out, " - Plugin finished: %s\n", plugin.Name())
					}
				}
				numsubs++
			}

			green.Fprintf(out, " - %v\n", display.English.Languages().Name(l))
		}
	}
	return numsubs, nil
}
