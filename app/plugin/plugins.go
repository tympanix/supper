package plugin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tympanix/supper/types"

	"github.com/apex/log"
	"github.com/mitchellh/mapstructure"
)

// NewFromMap construct a plugin from an map of attributes
func NewFromMap(i interface{}) (*Plugin, error) {
	plugin := &Plugin{}

	if err := mapstructure.Decode(i, plugin); err != nil {
		return nil, err
	}

	if err := plugin.valid(); err != nil {
		return nil, err
	}

	return plugin, nil
}

// Plugin is a struct enabling external functionality
type Plugin struct {
	PluginName string `mapstructure:"name"`
	Exec       string `mapstructure:"exec"`
}

// Run executes the plugin
func (p *Plugin) Run(s types.LocalSubtitle) error {
	cmd := exec.Command(shell[0], shell[1], p.Exec)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SUBTITLE=%s", s.Path()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).WithField("plugin", p.Name()).
			Debugf("Plugin debug\n%s", string(out))
	}
	return err
}

func (p *Plugin) valid() error {
	if p.Name() == "" {
		return fmt.Errorf("Missing plugin name")
	}
	if p.Exec == "" {
		return fmt.Errorf("Missing plugin exec for %v", p.Name())
	}
	return nil
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return p.PluginName
}
