package notify

import (
	"github.com/apex/log"
)

// AsyncLogger receives notification entries from a channel and logs them
func AsyncLogger() chan<- *Entry {
	c := make(chan *Entry)
	go func() {
		for e := range c {
			ctx := log.WithFields(log.Fields(e.Context))
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
	}()
	return c
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
