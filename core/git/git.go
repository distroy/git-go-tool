/*
 * Copyright (C) distroy
 */

package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Different struct {
	Filename  string
	BeginLine int
	EndLine   int
}

func (d Different) String() string {
	return fmt.Sprintf("%s:%d,%d", d.Filename, d.BeginLine, d.EndLine)
}

func GetBranch() string {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	_, err := cmd.Output()
	if err == nil {
		return "HEAD"
	}

	cmd = exec.Command("git", "hash-object", "-t", "tree", "/dev/null")
	out, err := cmd.Output()
	if err != nil {
		switch v := err.(type) {
		default:
			log.Fatalf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), v.Error())

		case *exec.ExitError:
			log.Fatalf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v.Error())
		}

		return ""
	}

	return string(out)
}

func GetRootDir() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		switch v := err.(type) {
		default:
			log.Fatalf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), v.Error())

		case *exec.ExitError:
			log.Fatalf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v.Error())
		}

		return ""
	}

	return strings.TrimSpace(string(out))
}

func getCommandString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := &strings.Builder{}
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
