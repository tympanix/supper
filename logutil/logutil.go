package logutil

import (
	"os"

	"github.com/tympanix/supper/types"

	"github.com/apex/log"
	clilog "github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

var handlers []log.Handler

// Initialize initialises the logger from a configuration
func Initialize(config types.Config) {
	// Use temporary logger during initialisation
	log.SetHandler(text.Default)

	logfile := config.Logfile()

	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			log.WithField("logfile", logfile).Error("could not open logfile for reading")
			os.Exit(1)
		}

		handlers = append(handlers, logfmt.New(file))
	}

	if config.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	handlers = append(handlers, clilog.Default)

	multilog := multi.New(handlers...)
	log.SetHandler(multilog)
}
