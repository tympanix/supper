package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/fatih/set"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/providers"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/urfave/cli"
)

func main() {

	sup := &application.Application{
		Provider: new(provider.Subscene),
	}

	app := cli.NewApp()
	app.Name = "supper"
	app.Version = "0.1.0"
	app.Usage = "An automatic subtitle downloader"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "lang, l",
			Usage: "subtitle language",
		},
		cli.BoolFlag{
			Name:  "impaired, i",
			Usage: "hearing impaired",
		},
		cli.IntFlag{
			Name:  "limit",
			Value: 12,
			Usage: "limit maximum number of media to process",
		},
		cli.StringFlag{
			Name:  "modified, m",
			Usage: "filter media modified within duration",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 1)
		}

		// Parse all language flags into slice of tags
		lang := set.New()
		for _, tag := range c.StringSlice("lang") {
			_lang, err := language.Parse(tag)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			lang.Add(_lang)
		}

		if lang.Size() == 0 {
			return cli.NewExitError("missing language flags(s)", 1)
		}

		// Make sure every arg is a valid file path
		for _, arg := range c.Args() {
			if _, err := os.Stat(arg); err == os.ErrNotExist {
				return cli.NewExitError(err, 1)
			}
		}

		// Search all argument paths for media
		media, err := sup.FindMedia(c.Args()...)

		if err != nil {
			return cli.NewExitError(err, 2)
		}

		modified, err := parse.Duration(c.String("modified"))

		if modified > 0 {
			media = media.FilterModified(modified)
		}

		media, err = media.FilterMissingSubs(lang)

		if err != nil {
			return cli.NewExitError(err, 2)
		}

		if err != nil {
			return cli.NewExitError(err, 3)
		}

		if media.Len() > c.Int("limit") {
			err := fmt.Errorf("number of media files exceeded: %v", media.Len())
			return cli.NewExitError(err, 3)
		}

		// Iterate all media files found in each path
		for i, item := range media.List() {
			cursubs, err := item.ExistingSubtitles()

			if err != nil {
				return cli.NewExitError(err, 2)
			}

			missingLangs := set.Difference(lang, cursubs.LanguageSet())

			if missingLangs.Size() == 0 {
				continue
			}

			fmt.Printf("(%v/%v) - %s\n", i+1, media.Len(), item)

			subs, err := sup.SearchSubtitles(item)

			if err != nil {
				return cli.NewExitError(err, 2)
			}

			subs = subs.HearingImpaired(c.Bool("impaired"))

			// Download subtitle for each language
			for _, l := range missingLangs.List() {
				l, ok := l.(language.Tag)

				if !ok {
					return cli.NewExitError(err, 3)
				}

				langsubs := subs.FilterLanguage(l)

				if langsubs.Len() == 0 {
					color.Red(" - no subtitles found")
					continue
				}

				err := item.SaveSubtitle(langsubs.Best())

				if err != nil {
					color.Red(err.Error())
					continue
				}
				color.Green(" - %v\n", display.English.Languages().Name(l))
			}
		}

		return nil
	}

	app.Run(os.Args)

}
