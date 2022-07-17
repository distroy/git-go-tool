/*
 * Copyright (C) distroy
 */

package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/distroy/git-go-tool/core/execcore"
)

type Different struct {
	Filename  string
	BeginLine int
	EndLine   int
}

func (d *Different) String() string {
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
	_, err := execcore.GetOutput("git", "rev-parse", "--verify", "HEAD")
	if err == nil {
		return "HEAD", nil
	}

	out, err := execcore.GetOutput("git", "hash-object", "-t", "tree", "/dev/null")
	if err != nil {
		return "", err
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
	out, err := execcore.GetOutput("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
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
