package main

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/tympanix/supper/app"
	"github.com/tympanix/supper/cmd"

	"github.com/urfave/cli"
)

// set application version with -ldflags -X
var version string

func main() {

	cmd.Execute()

}

func startWebServer(c *cli.Context) error {
	app := application.New(c)
	address := fmt.Sprintf(":%v", c.Int("port"))

	log.Infof("Listening on %v...\n", c.Int("port"))
	log.Error(http.ListenAndServe(address, app).Error())
	return nil
}
