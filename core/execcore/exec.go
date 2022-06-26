/*
 * Copyright (C) distroy
 */

package execcore

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/distroy/git-go-tool/core/strcore"
)

func GetOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	out, err := cmd.Output()
	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			return "", fmt.Errorf("exec command fail. cmd:%s, code:%d, err:%v",
				CmdString(cmd), v.ExitCode(), v.Error())
		}

		return "", fmt.Errorf("exec command fail. cmd:%s, err:%v", CmdString(cmd), err.Error())
	}

	return strcore.BytesToStrUnsafe(out), nil
}

func CmdString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := &strings.Builder{}
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
