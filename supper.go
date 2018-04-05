package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/logutil"
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

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the application version and exit",
	}

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
		cli.StringFlag{
			Name:   "config",
			Usage:  "load config file at `PATH` with additional configuration",
			EnvVar: "SUPPER_CONFIG",
		},
		cli.StringFlag{
			Name:   "logfile",
			Usage:  "file at `PATH` in which to store applications logs",
			EnvVar: "SUPPER_LOGFILE",
		},
		cli.BoolFlag{
			Name:   "verbose, v",
			Usage:  "enable verbose logging",
			EnvVar: "SUPPER_VERBOSE",
		},
	}

	app.Before = func(c *cli.Context) error {
		// Initialise logging
		logutil.Context(c)

		// Make sure all language flags are defined
		for _, tag := range c.StringSlice("lang") {
			_, err := language.Parse(tag)
			if err != nil {
				log.Fatalf("Unknown language tag: %v", tag)
			}
		}

		// Make sure score is between 0 and 100
		if c.Int("score") < 0 || c.Int("score") > 100 {
			log.Fatalf("Score must be between 0 and 100")
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
			log.Fatal("Missing language flag(s)")
		}

		// Make sure every arg is a valid file path
		for _, arg := range c.Args() {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				log.WithField("path", arg).Fatal("Invalid file path")
			}
		}

		// Search all argument paths for media
		media, err := sup.FindMedia(c.Args()...)

		if err != nil {
			log.WithError(err).Fatal("Online search failed")
		}

		modified, err := parse.Duration(c.String("modified"))

		if err != nil {
			log.WithError(err).WithField("modified", c.String("modified")).
				Fatal("Invalid duration")
		}

		if modified > 0 {
			media = media.FilterModified(modified)
		}

		if media.Len() > c.Int("limit") && !c.Bool("dry") && c.Int("limit") != -1 {
			log.WithFields(log.Fields{
				"media": strconv.Itoa(media.Len()),
				"limit": strconv.Itoa(c.Int("limit")),
			}).Fatal("Media limit exceeded")
		}

		numsubs, err := sup.DownloadSubtitles(media, lang, os.Stdout)

		if err != nil {
			log.WithError(err).Fatal("Download incomplete")
		}

		if c.Bool("dry") {
			ctx := log.WithField("reason", "dry-run")
			ctx.Warn("Nothing performed")
			ctx.Warnf("Media files: %v", media.Len())
			ctx.Warnf("Missing subtitles: %v", numsubs)
		}

		return nil
	}

	app.Run(os.Args)

}

func startWebServer(c *cli.Context) error {
	app := application.New(c)
	address := fmt.Sprintf(":%v", c.Int("port"))

	log.Infof("Listening on %v...\n", c.Int("port"))
	log.Error(http.ListenAndServe(address, app).Error())
	return nil
}
