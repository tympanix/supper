package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Tympanix/supper/app"
	"github.com/Tympanix/supper/providers"
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
		cli.StringFlag{
			Name:  "lang, l",
			Value: "en",
			Usage: "subtitle language",
		},
		cli.BoolFlag{
			Name:  "impaired, i",
			Usage: "hearing impaired",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 1)
		}

		media, err := sup.FindMedia(c.Args().First())

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		lang, err := language.Parse(c.String("lang"))

		if err != nil {
			log.Printf("Unknown language %s\n", c.String("lang"))
			os.Exit(1)
		}

		log.Printf("Finding subtitles for lang %s\n", c.String("lang"))

		if len(media) == 0 {
			return errors.New("No subtitles found")
		}

		subs, err := sup.SearchSubtitles(media[0])

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if c.Bool("impaired") {
			subs.RemoveNotHI()
		} else {
			subs.RemoveHI()
		}

		subs.FilterLanguage(lang)

		fmt.Println(subs)

		media[0].SaveSubtitle(subs.Best())

		return nil
	}

	app.Run(os.Args)

}
