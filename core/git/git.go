/*
 * Copyright (C) distroy
 */

package git

import (
	"fmt"
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

func MustGetBranch() string {
	branch, err := GetBranch()
	if err != nil {
		panic(err)
	}
	return branch
}

func GetBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	_, err := cmd.Output()
	if err == nil {
		return "HEAD", nil
	}

	cmd = exec.Command("git", "hash-object", "-t", "tree", "/dev/null")
	out, err := cmd.Output()
	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			return "", fmt.Errorf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v.Error())
		}

		return "", fmt.Errorf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), err.Error())
	}

	return string(out), nil
}

func MustGetRootDir() string {
	root, err := GetRootDir()
	if err != nil {
		panic(err)
	}
	return root
}

func GetRootDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			return "", fmt.Errorf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v.Error())
		}

		return "", fmt.Errorf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), err.Error())
	}

	return strings.TrimSpace(string(out)), nil
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
