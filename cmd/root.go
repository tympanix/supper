package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/cfg"
	"github.com/tympanix/supper/logutil"
)

var (
	appName    = "Supper"
	appDesc    = "Download subtitles in a breeze"
	appVersion = "master"  // set with ldflags -X
	appCommit  = "HEAD"    // set with ldflags -X
	appDate    = "unknown" // set with ldflags -X
)

// AppName returns the name of the application
func AppName() string {
	return appName
}

// AppVersion returns the version of the application
func AppVersion() string {
	return appVersion
}

// AppDesc return the description of the application
func AppDesc() string {
	return appDesc
}

// AppCommit returns the scm commit hash of the application
func AppCommit() string {
	return appCommit
}

// AppDate returns the build date of the application
func AppDate() string {
	return appDate
}

var rootCmd = &cobra.Command{
	Use:              strings.ToLower(AppName()),
	Short:            AppDesc(),
	PersistentPreRun: validateFlags,
	Args:             validateArgs,
	Run:              downloadSubtitles,
}

// Execute executes the CLI application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Application error")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(readConfigFiles)
	flags := rootCmd.PersistentFlags()

	// Set up cobra command line flags
	flags.StringSliceP("lang", "l", []string{}, "download subtitle in specified language")
	flags.BoolP("impaired", "i", false, "hearing impaired subtitles only")
	flags.Int("limit", 12, "limit maximum number of media to process")
	flags.StringP("modified", "m", "", "only process media modified within specified duration")
	flags.Bool("dry", false, "scan media but do not download any subtitles")
	flags.IntP("score", "s", 0, "only download subtitles ranking higher than specified percent")
	flags.String("delay", "", "wait specified duration before downloading next subtitle")
	flags.Bool("force", false, "overwrite existing subtitles")
	flags.String("config", "", "load config file at specified path")
	flags.String("logfile", "", "store application logs in specified path")
	flags.BoolP("verbose", "v", false, "enable verbose logging")
	flags.Bool("strict", false, "exit the application on any error instead of proceeding to next media item")
	flags.Bool("version", false, "show the application version and exit")

	// Bind flags to viper
	viper.BindPFlag("lang", flags.Lookup("lang"))
	viper.BindPFlag("impaired", flags.Lookup("impaired"))
	viper.BindPFlag("limit", flags.Lookup("limit"))
	viper.BindPFlag("modified", flags.Lookup("modified"))
	viper.BindPFlag("dry", flags.Lookup("dry"))
	viper.BindPFlag("score", flags.Lookup("score"))
	viper.BindPFlag("delay", flags.Lookup("delay"))
	viper.BindPFlag("force", flags.Lookup("force"))
	viper.BindPFlag("config", flags.Lookup("config"))
	viper.BindPFlag("logfile", flags.Lookup("logfile"))
	viper.BindPFlag("verbose", flags.Lookup("verbose"))
	viper.BindPFlag("strict", flags.Lookup("strict"))
	viper.BindPFlag("version", flags.Lookup("version"))

	// Set up aliases
	viper.RegisterAlias("lang", "languages")

	viper.SetDefault("author", "tympanix <tympanix@gmail.com>")
	viper.SetDefault("license", "GNUv3.0")
}

func readConfigFiles() {
	// Parse and set global configuration reference
	cfg.Initialize()

	// Set up logging capabilities from configuration
	logutil.Initialize(cfg.Default)

	config := viper.GetString("config")
	if config != "" {
		// Use config file from the flag
		viper.SetConfigFile(config)

		if err := viper.ReadInConfig(); err != nil {
			log.WithError(err).Fatal("Could not read config file")
			os.Exit(1)
		} else {
			log.WithField("file", viper.ConfigFileUsed()).
				Debug("Loaded configuration file")
		}
	} else {
		// Use default configuration
		viper.SetConfigName(strings.ToLower(AppName()))
		viper.AddConfigPath(cfg.DefaultPath(AppName()))
		if err := viper.ReadInConfig(); err == nil {
			log.WithField("file", viper.ConfigFileUsed()).
				Debug("Loaded default configuration")
		}

		// Merge in local configuration
		viper.SetConfigName(fmt.Sprintf(".%v", strings.ToLower(AppName())))
		viper.AddConfigPath(cfg.HomePath(AppName()))
		viper.AddConfigPath(".")

		if err := viper.MergeInConfig(); err != nil {
			// If no local configuration, use global configuration
			viper.AddConfigPath(cfg.GlobalPath(AppName()))
			viper.MergeInConfig()

			log.WithField("file", viper.ConfigFileUsed()).
				Debug("Loaded global configuration")
		} else {
			log.WithField("file", viper.ConfigFileUsed()).
				Debug("Loaded local configuration")
		}
	}

	// Parse and set global configuration reference
	cfg.Initialize()

	// Set up logging capabilities from configuration
	logutil.Initialize(cfg.Default)
}

func validateFlags(cmd *cobra.Command, args []string) {
	if cfg.Default.Languages().Size() == 0 {
		log.Fatal("Missing language flag(s)")
	}

	// Make sure score is between 0 and 100
	if cfg.Default.Score() < 0 || cfg.Default.Score() > 100 {
		log.WithField("score", cfg.Default.Score()).Fatalf("Score must be between 0 and 100")
	}
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if viper.GetBool("version") {
		showVersion(cmd, args)
		os.Exit(0)
	}

	if len(args) == 0 {
		log.WithField("args", fmt.Sprintf("%v", len(args))).
			Fatal("Missing media arguments")
	}

	// Make sure every arg is a valid file path
	for _, arg := range args {
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			log.WithField("path", arg).Fatal("Invalid file path")
		}
	}

	return nil
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
