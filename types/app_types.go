package types

import (
	"net/http"

	"github.com/urfave/cli"
)

type App interface {
	Provider
	http.Handler
	Context() *cli.Context
	FindMedia(...string) (LocalMediaList, error)
}
