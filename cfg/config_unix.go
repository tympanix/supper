// +build !windows

package cfg

import (
	"path/filepath"
	"strings"
)

// HomePath returns the configuration path for the users home
func HomePath(app string) string {
	return homePath
}

// GlobalPath returns the global configuration path
func GlobalPath(app string) string {
	return filepath.Join("/etc", strings.ToLower(app))
}

// DefaultPath returns the default configuration path
func DefaultPath(app string) string {
	return filepath.Join("/etc/defaults")
}
