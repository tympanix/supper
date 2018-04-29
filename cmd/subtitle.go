package cmd

import (
	"strconv"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/cfg"
)

func init() {
	flags := subtitleCmd.Flags()

	flags.IntP("score", "s", 0, "only download subtitles ranking higher than specified percent")
	flags.String("delay", "", "wait specified duration before downloading next subtitle")
	flags.StringSliceP("lang", "l", []string{}, "download subtitle in specified language")
	flags.BoolP("impaired", "i", false, "hearing impaired subtitles only")
	flags.Int("limit", 12, "limit maximum number of media to process")
	flags.StringP("modified", "m", "", "only process media modified within specified duration")

	viper.BindPFlag("languages", flags.Lookup("lang"))
	viper.BindPFlag("impaired", flags.Lookup("impaired"))
	viper.BindPFlag("limit", flags.Lookup("limit"))
	viper.BindPFlag("modified", flags.Lookup("modified"))
	viper.BindPFlag("score", flags.Lookup("score"))
	viper.BindPFlag("delay", flags.Lookup("delay"))

	rootCmd.AddCommand(subtitleCmd)
}

var subtitleCmd = &cobra.Command{
	Use:     "subtitle",
	Short:   "Download subtitles for media",
	Aliases: []string{"sub"},
	Args:    validateMedia,
	PreRun:  validateSubtitleFlags,
	Run:     downloadSubtitles,
}

func validateSubtitleFlags(cmd *cobra.Command, args []string) {
	if cfg.Default.Languages().Size() == 0 {
		log.Fatal("Missing language flag(s)")
	}

	// Make sure score is between 0 and 100
	if cfg.Default.Score() < 0 || cfg.Default.Score() > 100 {
		log.WithField("score", cfg.Default.Score()).Fatalf("Score must be between 0 and 100")
	}
}

func downloadSubtitles(cmd *cobra.Command, args []string) {
	// Create new application
	app := app.NewFromDefault()
	config := app.Config()

	// Search all argument paths for media
	media, err := app.FindMedia(args...)

	if err != nil {
		log.WithError(err).Fatal("Online search failed")
	}

	if config.Modified() > 0 {
		media = media.FilterModified(config.Modified())
	}

	if media.Len() > config.Limit() && !config.Dry() && config.Limit() != -1 {
		log.WithFields(log.Fields{
			"media": strconv.Itoa(media.Len()),
			"limit": strconv.Itoa(config.Limit()),
		}).Fatal("Media limit exceeded")
	}

	numsubs, err := app.DownloadSubtitles(media, config.Languages())

	if err != nil {
		log.WithError(err).Fatal("Download incomplete")
	}

	if config.Dry() {
		ctx := log.WithField("reason", "dry-run")
		ctx.Warn("Nothing performed")
		ctx.Warnf("Media files: %v", media.Len())
		ctx.Warnf("Missing subtitles: %v", numsubs)
	}
}
