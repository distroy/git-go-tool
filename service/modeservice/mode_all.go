/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"log"

	"github.com/distroy/git-go-tool/core/filecore"
)

type modeAll struct {
	modeBase
}

func (m *modeAll) mustInit(c *Config) {
	m.modeBase.mustInit(c)
}

func (m *modeAll) IsIn(file string, begin, end int) bool {
	if m.isFileIgnored(file) {
		return false
	}

	cache := m.cache
	ok, err := cache.CheckFileRange(&cacheCheckFileRangeReq{
		Filename:  file,
		BeginLine: begin,
		EndLine:   end,
		Check:     func(line string) bool { return !m.isLineIgnored(line) },
	})
	if err != nil {
		log.Fatalf("check file range fail. file:%s, begin:%d, end:%d, err:%v",
			file, begin, end, err)
	}
	return ok
}

func (m *modeAll) Walk(fn WalkFunc) {
	cache := m.cache
	rootDir := m.gitRoot

	err := cache.WalkFiles(func(file *filecore.File) error {
		m.mustWalkFile(file, fn)
		return nil
	})

	if err != nil {
		log.Fatalf("walk dir fail. dir:%s, err:%v", rootDir, err)
	}
}
