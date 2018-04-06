package cmd

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tympanix/supper/app"
)

func init() {
	webCmd.Flags().IntP("port", "p", 5670, "port used to serve the web application")
	webCmd.Flags().String("movies", "", "path to your movie collection")
	webCmd.Flags().String("shows", "", "path to your tv show collection")
	webCmd.Flags().String("static", "", "path to the web files to serve")

	viper.BindPFlag("port", webCmd.Flags().Lookup("port"))
	viper.BindPFlag("movies", webCmd.Flags().Lookup("movies"))
	viper.BindPFlag("shows", webCmd.Flags().Lookup("shows"))
	viper.BindPFlag("static", webCmd.Flags().Lookup("static"))

	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Listen and server the web application",
	Run:   startWebServer,
}

func startWebServer(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()
	address := fmt.Sprintf(":%v", viper.GetInt("port"))

	log.Infof("Listening on %v...\n", viper.GetInt("port"))
	log.Error(http.ListenAndServe(address, app).Error())
}
