/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"log"
	"os"
	"path/filepath"
)

type modeAll struct {
	modeBase
}

func (m *modeAll) mustInit(c *Config) {
	m.modeBase.mustInit(c)
}

func (m *modeAll) IsIn(file string, begin, end int) bool {
	cache := m.cache

	ok, err := cache.CheckFileRange(file, begin, end, func(line string) bool {
		return !m.isLineIgnored(line)
	})
	if err != nil {
		log.Fatalf("check file range fail. file:%s, begin:%d, end:%d, err:%v",
			file, begin, end, err)
	}
	return ok
}

func (m *modeAll) Walk(fn WalkFunc) {
	rootDir := m.rootDir

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			m.mustWalkFile(path, fn)
		}
		return err
	})

	if err != nil {
		log.Fatalf("walk dir fail. dir:%s, err:%v", rootDir, err)
	}
}
