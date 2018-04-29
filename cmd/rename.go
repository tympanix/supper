package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/media"
)

func init() {
	actions := make([]string, 0)
	for a := range app.Renamers {
		actions = append(actions, a)
	}

	flags := renameCmd.Flags()

	flags.StringP("action", "a", "hardlink", fmt.Sprintf("renaming action %v", strings.Join(actions, "|")))
	flags.BoolP("extract", "x", false, "extract media from archives")

	viper.BindPFlag("action", flags.Lookup("action"))
	viper.BindPFlag("extract", flags.Lookup("extract"))

	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename and process media files",
	Aliases: []string{"ren"},
	Args:    validateMedia,
	PreRun:  validateRenameFlags,
	Run:     renameMedia,
}

func validateRenameFlags(cmd *cobra.Command, args []string) {
	if _, ok := app.Renamers[viper.GetString("action")]; !ok {
		log.Fatalf("Invalid action flag %v", viper.GetString("action"))
	}
}

func renameMedia(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()

	medialist, err := app.FindMedia(args...)

	if err != nil {
		log.WithError(err).Fatal("Could not find media in path")
	}

	if err := app.RenameMedia(medialist); err != nil {
		log.WithError(err).Fatal("Could not rename media files")
	}

	if viper.GetBool("extract") {
		archives, err := app.FindArchives(args...)
		if err != nil {
			log.WithError(err).Fatal("Could not open archives")
		}
		for _, a := range archives {
			defer a.Close()

			m, err := a.Next()
			for err == nil {
				if err = app.ExtractMedia(m); err != nil {
					if !media.IsExistsErr(err) {
						if app.Config().Strict() {
							log.WithError(err).Fatal("Extraction failed")
						} else {
							log.WithError(err).Error("Extraction failed")
						}
					}
				}
				defer m.Close()
				m, err = a.Next()
			}

			if err != io.EOF {
				log.WithError(err).Fatal("Extraction failed")
			}
		}
	}

}
