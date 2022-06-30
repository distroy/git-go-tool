/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/git"
)

type modeBase struct {
	config  *Config
	gitRoot string
	gitSubs []*git.SubModule
	cache   *cache
}

func (m *modeBase) mustInit(c *Config) {
	m.config = c
	m.gitRoot = git.MustGetRootDir()
	m.gitSubs = git.MustGetSubModules()
	m.cache = newCache(m.gitRoot)
}

func (m *modeBase) isFileIgnored(file string) bool {
	if m.isGitSub(file) {
		return true
	}
	if m.config.FileFilter == nil {
		return false
	}
	return !m.config.FileFilter(file)
}

func (m *modeBase) isLineIgnored(line string) bool {
	line = strings.TrimSpace(line)
	// return len(line) == 0
	return len(line) == 0 || line == "}"
}

func (m *modeBase) isGitSub(filename string) bool {
	for _, sub := range m.gitSubs {
		if strings.HasPrefix(filename, sub.Path) {
			return true
		}
	}
	return false
}

func (m *modeBase) mustWalkFile(file *filecore.File, fn WalkFunc) {
	filename := file.Name

	if !file.IsGo() {
		return
	}

	if m.isFileIgnored(file.Name) {
		return
	}

	cache := m.cache
	defer cache.Del(filename)

	lines := file.MustReadLines()

	f := file.MustParse()
	for _, decl := range f.Decls {
		fDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		body := fDecl.Body
		if body == nil {
			// log.Printf(" === func name:%s", fDecl.Name.Name)
			continue
		}

		pos := file.Position(body.Lbrace)
		end := file.Position(body.Rbrace)
		for i := pos.Line - 1; i < end.Line; i++ {
			if m.isLineIgnored(lines[i]) {
				continue
			}

			fn(filename, i+1, i+1)
		}
	}
}
