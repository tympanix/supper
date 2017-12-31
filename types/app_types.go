package types

import (
	"net/http"

	"github.com/fatih/set"
	"github.com/urfave/cli"
)

type App interface {
	Provider
	http.Handler
	Context() *cli.Context
	FindMedia(...string) (LocalMediaList, error)
	Languages() set.Interface
}
