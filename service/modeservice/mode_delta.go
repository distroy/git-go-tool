/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"log"

	"github.com/distroy/git-go-tool/core/git"
)

type modeDelta struct {
	modeBase

	files git.Files
}

func (m *modeDelta) mustInit(c *Config) {
	branch := c.Branch
	if branch == "" {
		branch = git.GetBranch()
	}

	s, err := git.ParseNewLines(branch)
	if err != nil {
		log.Fatalf("parse the git different relative to the branch:%s. err:%v", branch, err)
	}

	m.files = git.NewFileDifferents(s)
}

func (m *modeDelta) IsIn(file string, begin, end int) bool {
	return m.files.IsIn(file, begin, end)
}

func (m *modeDelta) Walk(fn WalkFunc) {
	for _, file := range m.files {
		for _, diff := range file {
			fn(diff.Filename, diff.BeginLine, diff.EndLine)
		}
	}
}
