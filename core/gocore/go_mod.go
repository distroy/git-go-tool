/*
 * Copyright (C) distroy
 */

package gocore

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/distroy/git-go-tool/core/git"
)

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

func mustGetCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func GetModPrefix() string {
	gitRootDir := git.GetRootDir()

	goModFile := path.Join(gitRootDir, "go.mod")
	if _, err := os.Stat(goModFile); err != nil {
		log.Fatalf("cannot find the go.mod file. work path:%s, git root:%s, err:%v",
			mustGetCwd(), gitRootDir, err)
	}

	cmd := exec.Command("grep", "module", goModFile)
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

	prefix := strings.Split(string(out), " ")[1]
	prefix = strings.TrimSpace(prefix)

	return prefix
}
