package parse_test

import (
	"log"
	"testing"

	"github.com/Tympanix/supper/parse"
)

func TestCleanName(t *testing.T) {

	names := []string{
		"The Office (US)",
		"Marvel's Agents of S.H.I.E.L.D.",
		"The O.C.",
		"The Flash (2014)",
	}

	for _, name := range names {
		log.Println(parse.CleanName(name))
	}
}
