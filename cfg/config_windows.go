// +build windows

package cfg

import (
	"path/filepath"
)

// HomePath returns the configuration path for the users home
func HomePath(app string) string {
	return homePath
}

// GlobalPath returns the global configuration path
func GlobalPath(app string) string {
	return filepath.Join(homePath, "AppData", "Roaming", app)
}

// DefaultPath returns the default configuration path
func DefaultPath(app string) string {
	return filepath.Join("C:\\ProgramData", app)
}
