package main

import (
	"fmt"
	"os"

	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/providers"
	"golang.org/x/text/language"

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
		lang := make([]language.Tag, 0)
		for _, tag := range c.StringSlice("lang") {
			_lang, err := language.Parse(tag)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			lang = append(lang, _lang)
		}

		if len(lang) == 0 {
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

		if err != nil {
			return cli.NewExitError(err, 3)
		}

		if media.Len() > c.Int("limit") {
			err := fmt.Errorf("number of media files exceeded: %v", media.Len())
			return cli.NewExitError(err, 3)
		}

		// Iterate all media files found in each path
		for _, item := range media.List() {
			subs, err := sup.SearchSubtitles(item)

			if err != nil {
				return cli.NewExitError(err, 2)
			}

			subs = subs.HearingImpaired(c.Bool("impaired"))

			// Download subtitle for each language
			for _, l := range lang {
				langsubs := subs.FilterLanguage(l)

				fmt.Println(langsubs)

				item.SaveSubtitle(langsubs.Best())
			}
		}

		return nil
	}

	app.Run(os.Args)

}
