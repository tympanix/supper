//go:generate statik -src=web/build

package main

import "github.com/tympanix/supper/app/cli"

func main() {
	cli.Execute()
}
