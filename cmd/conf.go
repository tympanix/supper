package cmd

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/cfg"
)

func init() {
	rootCmd.AddCommand(confCmd)
}

var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Show application configuration and exit",
	Run:   showApplicationConfiguration,
}

func showApplicationConfiguration(cmd *cobra.Command, args []string) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "help" {
			return
		}
		log.WithField(flag.Name, fmt.Sprintf("%v", viper.Get(flag.Name))).Info(flag.Usage)
	})

	for _, p := range cfg.Default.Plugins() {
		log.Info(p.Name())
	}
}
