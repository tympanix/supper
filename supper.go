package main

import (
	"fmt"
	"log"

	"github.com/Tympanix/supper/providers"
	"github.com/Tympanix/supper/types"
)

// Application is an configuration instance of the application
type Application struct {
	types.Provider
}

func main() {

	app := &Application{
		new(provider.Subscene),
	}

	subs, err := app.Search("Guardians of the galaxy")

	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range subs {
		fmt.Println(s)
	}
}
