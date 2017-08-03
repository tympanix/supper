package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Tympanix/supper/app"
	"github.com/Tympanix/supper/providers"
)

func main() {

	app := &application.Application{
		Provider: new(provider.Subscene),
	}

	flag.Parse()
	root := flag.Arg(0)

	if len(root) == 0 {
		log.Println("Missing file root")
		os.Exit(1)
	}

	log.Printf("Walking: %s\n", root)

	media, err := app.FindMedia(root)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, m := range media {
		fmt.Println(m)
	}

}
