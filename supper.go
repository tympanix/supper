package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Tympanix/supper/providers"
	"github.com/Tympanix/supper/types"
)

// Application is an configuration instance of the application
type Application struct {
	types.Provider
}

func main() {

	file, _ := os.Open("./supper.go")

	fmt.Println(path.Base(file.Name()))

	return

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
