package notify

func empty() Context {
	return Context(make(map[string]interface{}))
}

// Debug creates a debuggin notification
func Debug(s string, v ...interface{}) *Entry {
	return empty().Debug(s, v...)
}

// Info creates a info notification
func Info(s string, v ...interface{}) *Entry {
	return empty().Info(s, v...)
}

// Warn creates a warn notification
func Warn(s string, v ...interface{}) *Entry {
	return empty().Warn(s, v...)
}

// Error creates a warn notification
func Error(s string, v ...interface{}) *Entry {
	return empty().Error(s, v...)
}

// Fatal creates a warn notification
func Fatal(s string, v ...interface{}) *Entry {
	return empty().Fatal(s, v...)
}
