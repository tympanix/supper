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

// set application version with -ldflags -X
var version string

func main() {

	app := cli.NewApp()
	app.Name = "supper"
	app.Version = version
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
			Name:   "lang, l",
			Usage:  "only download subtitles in language `LANG`",
			EnvVar: "SUPPER_LANGS",
		},
		cli.BoolFlag{
			Name:   "impaired, i",
			Usage:  "hearing impaired subtitles only",
			EnvVar: "SUPPER_IMPAIRED",
		},
		cli.IntFlag{
			Name:   "limit",
			Value:  12,
			Usage:  "limit maximum number of media to process to `NUM`",
			EnvVar: "SUPPER_LIMIT",
		},
		cli.StringFlag{
			Name:   "modified, m",
			Usage:  "only process media modified within `TIME` duration",
			EnvVar: "SUPPER_MODIFIED",
		},
		cli.BoolFlag{
			Name:  "dry, dry-run",
			Usage: "scan media but do not download any subtitles",
		},
		cli.IntFlag{
			Name:   "score, s",
			Value:  0,
			Usage:  "only download subtitles ranking higher than `SCORE` percent",
			EnvVar: "SUPPER_SCORE",
		},
		cli.StringFlag{
			Name:   "delay",
			Usage:  "wait `TIME` duration before downloading next subtitle",
			EnvVar: "SUPPER_DELAY",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "overwrite existing subtitles",
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

		// Make sure score is between 0 and 100
		if c.Int("score") < 0 || c.Int("score") > 100 {
			return cli.NewExitError("score must be between 0 and 100", 1)
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

		if err != nil {
			return cli.NewExitError(err, 3)
		}

		if media.Len() > c.Int("limit") && !c.Bool("dry") {
			err := fmt.Errorf("number of media files exceeded: %v", media.Len())
			return cli.NewExitError(err, 3)
		}

		numsubs, err := sup.DownloadSubtitles(media, lang, os.Stdout)

		if err != nil {
			return cli.NewExitError(err, 5)
		}

		if c.Bool("dry") {
			fmt.Println()
			color.Blue("dry run, nothing performed")
			color.Blue("total media files: %v", media.Len())
			color.Blue("total missing subtitles: %v", numsubs)
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
