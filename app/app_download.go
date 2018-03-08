package application

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
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

	// Iterate all media files found in each path
	for i, item := range media.List() {
		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return -1, cli.NewExitError(err, 2)
		}

		missingLangs := set.Difference(lang, cursubs.LanguageSet())

		if missingLangs.Size() == 0 {
			continue
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
				err := item.SaveSubtitle(onl, onl.Language())
				if err != nil {
					red.Fprintln(out, err.Error())
					continue
				}
				numsubs++
			}

			green.Fprintf(out, " - %v\n", display.English.Languages().Name(l))
		}
	}
	return numsubs, nil
}
