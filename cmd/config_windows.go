// +build windows

package cmd

import (
	"path/filepath"

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
	configGlobalPath = filepath.Join(home, "AppData", "Roaming", AppName)
	configDefaultsPath = filepath.Join("C:\\ProgramData", AppName)
}
