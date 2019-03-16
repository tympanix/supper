package logutil

import (
	"errors"
	"fmt"
	"os"

	"github.com/tympanix/supper/types"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/memory"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

// Initialize initialises the logger from a configuration
func Initialize(config types.Config) {
	// Use temporary logger during initialisation
	log.SetHandler(text.Default)

	var handlers []log.Handler

	if config.Logfile() != "" {
		file, err := os.OpenFile(config.Logfile(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			log.WithField("logfile", config.Logfile()).Fatal("could not open logfile for reading")
		}

		handlers = append(handlers, logfmt.New(file))
	}

	if config.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	handlers = append(handlers, cli.Default)

	multilog := multi.New(handlers...)
	log.SetHandler(multilog)
}

// Mock substitues the handler with an in-memory handler which can be used for
// testing purposes
func Mock() *memory.Handler {
	h := memory.New()
	log.SetHandler(h)
	return h
}

// Error logs the error message and returns an error with the same message
func Error(msg string) error {
	log.Error(msg)
	return errors.New(msg)
}

// Errorf formats and logs the error message and returns an error with the same
// message
func Errorf(msg string, v ...interface{}) error {
	m := fmt.Sprintf(msg, v...)
	log.Error(m)
	return errors.New(m)
}
