// +build windows

package plugin

import (
	"fmt"
	"strings"
)

var shell = []string{"cmd.exe", "/C"}

func escape(cmd string, args ...string) string {
	allargs := []string{cmd}

	for _, a := range args {
		allargs = append(allargs, fmt.Sprintf("\"%s\"", a))
	}
	return strings.Join(allargs, " ")
}
