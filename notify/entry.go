package notify

// Level is the notification level
type Level uint32

const (
	// LevelDebug is for development notifications
	LevelDebug Level = iota

	// LevelInfo is for notifictions which are non-critical
	LevelInfo

	// LevelWarn is for notifications which may be troublesome
	LevelWarn

	// LevelError is for notifications which is crtical
	LevelError

	// LevelFatal is for notifications which are fatal
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// MarshalJSON returns the notification level as a string
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte("\"" + l.String() + "\""), nil
}

// Fields is a collection of key, value pairs
type Fields map[string]interface{}

// Entry is a single instance of a notification
type Entry struct {
	Context
	Message string `json:"message"`
	Level   Level  `json:"level"`
}

// Error returns a string representation of the notification
func (e *Entry) Error() string {
	return e.Message
}

// WithField creates a new context with a single field
func WithField(name string, value interface{}) Context {
	return empty().addField(name, value)
}

// WithFields creates a new context with multiple fields
func WithFields(f Fields) Context {
	c := empty()
	for k, v := range f {
		c.addField(k, v)
	}
	return c
}

// WithError creates a new context with an error
func WithError(err error) Context {
	return empty().addField("error", err)
}

// WithExtra creates a new context with an extra property
func WithExtra(name string, value interface{}) Context {
	return empty().addExtra(name, value)
}
