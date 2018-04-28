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
	webCmd.Flags().String("static", "", "path to the web files to serve")

	viper.BindPFlag("port", webCmd.Flags().Lookup("port"))
	viper.BindPFlag("static", webCmd.Flags().Lookup("static"))

	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Listen and serve the web application",
	Args:  cobra.NoArgs,
	Run:   startWebServer,
}

func startWebServer(cmd *cobra.Command, args []string) {
	app := app.NewFromDefault()
	address := fmt.Sprintf(":%v", viper.GetInt("port"))

	log.Infof("Listening on %v...\n", viper.GetInt("port"))
	log.WithError(http.ListenAndServe(address, app)).
		Fatal("Web application exited abnormally")
}
