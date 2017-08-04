package main_test

import (
	"fmt"
	"testing"

	"github.com/Tympanix/supper/app"
	"github.com/Tympanix/supper/providers"
)

func TestSupper(t *testing.T) {
	app := &application.Application{
		Provider: new(provider.Subscene),
	}

	media, _ := app.FindMedia(`E:\Media\Movies\Inception (2010)\`)

	subs, _ := app.SearchSubtitles(media[0])

	fmt.Println(subs)

}
