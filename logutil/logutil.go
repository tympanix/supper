package logutil

import (
	"os"

	"github.com/apex/log"
	clilog "github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
	"github.com/urfave/cli"
)

var handlers []log.Handler

// Context initialises the logger from a cli context
func Context(ctx *cli.Context) {
	// Use temporary logger during initialisation
	log.SetHandler(text.Default)

	logfile := ctx.GlobalString("logfile")

	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			log.WithField("logfile", logfile).Error("could not open logfile for reading")
			os.Exit(1)
		}

		handlers = append(handlers, logfmt.New(file))
	}

	if verbose := ctx.GlobalBool("verbose"); verbose {
		log.SetLevel(log.DebugLevel)
	}

	handlers = append(handlers, clilog.Default)

	multilog := multi.New(handlers...)
	log.SetHandler(multilog)
}
