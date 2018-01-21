package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/parse"
	"golang.org/x/text/language"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "supper"
	app.Version = "0.1.0"
	app.Usage = "An automatic subtitle downloader"

	app.Commands = []cli.Command{
		{
			Name:   "web",
			Usage:  "listens and serves the web application",
			Action: startWebServer,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port, p",
					Value: 5670,
					Usage: "port used to serve the web application",
				},
				cli.StringFlag{
					Name:  "movies",
					Usage: "path to your movie collection",
				},
				cli.StringFlag{
					Name:  "shows",
					Usage: "path to your tv show collection",
				},
				cli.StringFlag{
					Name:  "static",
					Value: "./web",
					Usage: "path to the web files to serve",
				},
			},
		},
	}

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
		cli.BoolFlag{
			Name:  "dry",
			Usage: "scan media but do not download any subtitles",
		},
	}

	app.Before = func(c *cli.Context) error {
		// Make sure all language flags are defined
		for _, tag := range c.StringSlice("lang") {
			_, err := language.Parse(tag)
			if err != nil {
				err := fmt.Errorf("unknown language tag: %v", tag)
				return cli.NewExitError(err, 1)
			}
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		sup := application.New(c)

		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 1)
		}

		lang := sup.Languages()

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

		if err != nil {
			return cli.NewExitError(err, 2)
		}

		if modified > 0 {
			media = media.FilterModified(modified)
		}

		media, err = media.FilterMissingSubs(lang)

		if err != nil {
			return cli.NewExitError(err, 3)
		}

		if media.Len() > c.Int("limit") && !c.Bool("dry") {
			err := fmt.Errorf("number of media files exceeded: %v", media.Len())
			return cli.NewExitError(err, 3)
		}

		if err := sup.DownloadSubtitles(media, lang, os.Stdout); err != nil {
			return cli.NewExitError(err, 5)
		}

		if c.Bool("dry") {
			fmt.Println()
			color.Blue("dry run, nothing performed")
			color.Blue("total media files: %v", media.Len())
			//color.Blue("total missing subtitles: %v", numsubs)
		}

		return nil
	}

	app.Run(os.Args)

}

func startWebServer(c *cli.Context) error {
	app := application.New(c)
	address := fmt.Sprintf(":%v", c.Int("port"))

	log.Printf("Listening on %v...", c.Int("port"))
	log.Println(http.ListenAndServe(address, app))
	return nil
}
