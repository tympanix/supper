// +build !windows

package cmd

import (
	"path/filepath"
	"strings"

	"github.com/apex/log"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	configHomePath     string
	configGlobalPath   string
	configDefaultsPath string
)

func init() {
	home, err := homedir.Dir()

	if err != nil {
		log.WithError(err).Fatal("Could not find user home directory")
	}

	configHomePath = home
	configGlobalPath = filepath.Join("/etc", strings.ToLower(AppName))
	configDefaultsPath = filepath.Join("/etc/defaults")
}
