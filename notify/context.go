package notify

import "fmt"

// Map holds stateful information of key, value pairs
type Map map[string]interface{}

func (m Map) copy() Map {
	ctx := make(map[string]interface{})
	for k, v := range m {
		ctx[k] = v
	}
	return Map(ctx)
}

// Context holds stateful information for notifications
type Context struct {
	Fields Map `json:"data"`
	Extra  Map `json:"extra"`
}

func empty() Context {
	return Context{
		Fields: make(map[string]interface{}),
		Extra:  make(map[string]interface{}),
	}
}

func (c Context) addField(name string, value interface{}) Context {
	c.Fields[name] = value
	return c
}

func (c Context) addExtra(name string, value interface{}) Context {
	c.Extra[name] = value
	return c
}

func (c Context) copy() Context {
	return Context{
		Fields: c.Fields.copy(),
		Extra:  c.Extra.copy(),
	}
}

// WithField extends the current context with anohter field
func (c Context) WithField(name string, value interface{}) Context {
	return c.copy().addField(name, value)
}

// WithFields extends the current context with multiple fields
func (c Context) WithFields(f Fields) Context {
	copy := c.copy()
	for k, v := range f {
		copy.addField(k, v)
	}
	return copy
}

// WithError extedns the current context with an error
func (c Context) WithError(err error) Context {
	copy := c.copy()
	copy.addField("error", err)
	return copy
}

// WithExtra adds additional information which is only used internally
func (c Context) WithExtra(name string, value interface{}) Context {
	return c.copy().addExtra(name, value)
}

// Notify creates a new notification entry
func (c Context) Notify(l Level, s string, v ...interface{}) *Entry {
	return &Entry{
		Context: c,
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
