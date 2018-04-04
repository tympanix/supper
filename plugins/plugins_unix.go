// +build !windows

package plugins

import shellquote "github.com/kballard/go-shellquote"

var shell = []string{"sh", "-c"}

func escape(cmd string, args ...string) string {
	cmds, err := shellquote.Split(cmd)
	if err != nil {
		panic(err)
	}
	return shellquote.Join(append(cmds, args...)...)
}
