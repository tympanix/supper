package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	lang   []string
	config string
}

var cfg = config{}

var rootCmd = &cobra.Command{
	Use:   "supper",
	Short: "Supper downloads subtitles in a breeze",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetStringSlice("lang"))
	},
}

// Execute executes the CLI application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
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
	viper.BindPFlag("logfile", flags.Lookup("logfile"))
	viper.BindPFlag("verbose", flags.Lookup("verbose"))
	viper.BindPFlag("strict", flags.Lookup("strict"))

	viper.SetDefault("author", "tympanix <tympanix@gmail.com>")
	viper.SetDefault("license", "GNUv3.0")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfg.config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg.config)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("supper")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
