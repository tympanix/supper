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

	renameCmd.Flags().StringP("action", "a", "hardlink", fmt.Sprintf("renaming action %v", strings.Join(actions, "|")))
	renameCmd.Flags().StringP("dir", "d", "", "output directory for renamed files")

	viper.BindPFlag("action", renameCmd.Flags().Lookup("action"))
	viper.BindPFlag("dir", renameCmd.Flags().Lookup("dir"))
	viper.BindPFlag("conflict", renameCmd.Flags().Lookup("conflict"))

	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:    "rename",
	Short:  "Rename and process media files",
	PreRun: validateRenameFlags,
	Run:    rename,
}

func validateRenameFlags(cmd *cobra.Command, args []string) {
	if _, ok := app.Renamers[viper.GetString("action")]; !ok {
		log.Fatalf("invalid action flag %v", viper.GetString("action"))
	}
}

func rename(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()

	medialist, err := app.FindMedia(args...)

	if err != nil {
		log.WithError(err).Fatal("could not find media in path")
	}

	if err := app.RenameMedia(medialist); err != nil {
		log.WithError(err).Fatal("could not rename media files")
	}

}
