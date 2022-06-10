/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"log"
	"path"

	"github.com/distroy/git-go-tool/core/git"
)

type modeDelta struct {
	modeBase

	files git.Files
}

func (m *modeDelta) mustInit(c *Config) {
	m.modeBase.mustInit(c)

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
	for _, differents := range m.files {
		if len(differents) == 0 {
			continue
		}

		filename := differents[0].Filename
		// if filename == "/dev/null" {
		// 	continue
		// }

		filePath := path.Join(m.rootDir, filename)
		m.mustWalkFile(filePath, func(file string, begin, end int) {
			for i := begin; i <= end; i++ {
				if !m.IsIn(file, begin, end) {
					continue
				}
				fn(file, i, i)
			}
		})
	}
}
