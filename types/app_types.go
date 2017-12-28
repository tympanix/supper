package types

import (
	"net/http"

	"github.com/urfave/cli"
)

type App interface {
	Provider
	http.Handler
	Args() cli.Args
	FindMedia(...string) (LocalMediaList, error)
}
