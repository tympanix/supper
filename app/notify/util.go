package notify

import (
	"github.com/apex/log"
)

// AsyncLogger receives notification entries from a channel and logs them
func AsyncLogger() (chan<- *Entry, chan bool) {
	c := make(chan *Entry)
	d := make(chan bool, 1)
	go func(c <-chan *Entry, d chan<- bool) {
		defer func() {
			d <- true
		}()
		for e := range c {
			ctx := log.WithFields(log.Fields(e.Fields))
			switch e.Level {
			case LevelDebug:
				ctx.Debug(e.Message)
			case LevelInfo:
				ctx.Info(e.Message)
			case LevelWarn:
				ctx.Warn(e.Message)
			case LevelError:
				ctx.Error(e.Message)
			case LevelFatal:
				ctx.Fatal(e.Message)
			default:
				ctx.Error(e.Message)
			}
		}
	}(c, d)
	return c, d
}

// AsyncDiscard discards all notifications
func AsyncDiscard() chan<- *Entry {
	c := make(chan *Entry)
	go func() {
		for _ = range c {
		}
	}()
	return c
}
