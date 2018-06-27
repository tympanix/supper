package notify

import "fmt"

// Context holds stateful information for notifications
type Context map[string]interface{}

func (c Context) copy() Context {
	ctx := make(map[string]interface{})
	for k, v := range c {
		ctx[k] = v
	}
	return Context(ctx)
}

// WithField extends the current context with anohter field
func (c Context) WithField(name string, value interface{}) Context {
	copy := c.copy()
	copy[name] = value
	return copy
}

// WithFields extends the current context with multiple fields
func (c Context) WithFields(f Fields) Context {
	copy := c.copy()
	for k, v := range f {
		copy[k] = v
	}
	return copy
}

// WithError extedns the current context with an error
func (c Context) WithError(err error) Context {
	copy := c.copy()
	copy["error"] = err
	return copy
}

// Notify creates a new notification entry
func (c Context) Notify(l Level, s string, v ...interface{}) *Entry {
	return &Entry{
		Context: c.copy(),
		Level:   l,
		Message: fmt.Sprintf(s, v...),
	}
}

// Debug creates a new debugging notification
func (c Context) Debug(s string, v ...interface{}) *Entry {
	return c.Notify(LevelDebug, s, v...)
}

// Info creates a new info notification
func (c Context) Info(s string, v ...interface{}) *Entry {
	return c.Notify(LevelInfo, s, v...)
}

// Warn creates a new info notification
func (c Context) Warn(s string, v ...interface{}) *Entry {
	return c.Notify(LevelWarn, s, v...)
}

// Error creates a new info notification
func (c Context) Error(s string, v ...interface{}) *Entry {
	return c.Notify(LevelError, s, v...)
}

// Fatal creates a new info notification
func (c Context) Fatal(s string, v ...interface{}) *Entry {
	return c.Notify(LevelFatal, s, v...)
}
