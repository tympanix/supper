package cmd

import (
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

	flags.StringP("action", "a", "hardlink", strings.Join(actions, "|"))
	flags.BoolP("extract", "x", false, "extract media from archives")
	flags.BoolP("movies", "m", false, "rename only movies")
	flags.BoolP("tvshows", "t", false, "rename only tv shows")
	flags.BoolP("subtitles", "s", false, "rename only subtitles")
	flags.BoolP("singular", "i", false, "disallow duplicates of same media")
	flags.BoolP("upgrades", "u", false, "allow duplicate of media when of better quality")


	viper.BindPFlag("action", flags.Lookup("action"))
	viper.BindPFlag("extract", flags.Lookup("extract"))
	viper.BindPFlag("movies", flags.Lookup("movies"))
	viper.BindPFlag("tvshows", flags.Lookup("tvshows"))
	viper.BindPFlag("subtitles", flags.Lookup("subtitles"))
	viper.BindPFlag("singular", flags.Lookup("singular"))
	viper.BindPFlag("upgrades", flags.Lookup("upgrades"))

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

	var selected int
	for _, b := range []bool{
		viper.GetBool("movies"),
		viper.GetBool("tvshows"),
		viper.GetBool("subtitles"),
	} {
		if b {
			selected++
		}
	}

	if selected > 1 {
		log.Fatal("renaming movies, tvshows and subtitles are mutually exclusive")
	}

	if viper.GetBool("singular") && viper.GetBool("upgrades") {
		log.Fatal("flags singular and upgrades are mutually exclusive")
	}
}

func renameMedia(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()

	medialist, err := app.FindMedia(args...)

	if err != nil {
		log.WithError(err).Fatal("Could not find media in path")
	}

	if viper.GetBool("movies") {
		medialist = medialist.FilterMovies()
	}

	if viper.GetBool("tvshows") {
		medialist = medialist.FilterEpisodes()
	}

	if viper.GetBool("subtitles") {
		medialist = medialist.FilterSubtitles()
	}

	if err := app.RenameMedia(medialist); err != nil {
		log.WithError(err).Fatal("Could not rename media files")
	}

	if viper.GetBool("extract") {
		extractMedia(cmd, args)
	}
}

func extractMedia(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()

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
