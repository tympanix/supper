package cmd

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/app"
)

func init() {
	actions := make([]string, 0)
	for a := range app.Renamers {
		actions = append(actions, a)
	}

	flags := renameCmd.Flags()

	flags.StringP("action", "a", "hardlink", fmt.Sprintf("renaming action %v", strings.Join(actions, "|")))

	viper.BindPFlag("action", renameCmd.Flags().Lookup("action"))

	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename and process media files",
	Aliases: []string{"ren"},
	PreRun:  validateRenameFlags,
	Run:     renameMedia,
}

func validateRenameFlags(cmd *cobra.Command, args []string) {
	if _, ok := app.Renamers[viper.GetString("action")]; !ok {
		log.Fatalf("invalid action flag %v", viper.GetString("action"))
	}
}

func renameMedia(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()

	medialist, err := app.FindMedia(args...)

	if err != nil {
		log.WithError(err).Fatal("could not find media in path")
	}

	if err := app.RenameMedia(medialist); err != nil {
		log.WithError(err).Fatal("could not rename media files")
	}

}
