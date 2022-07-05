/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/git"
)

var (
	errInvalidRange = fmt.Errorf("invalid file range")
)

type modeBase struct {
	gitRoot string
	gitSubs []*git.SubModule
	cache   *filecore.Cache
}

func (m *modeBase) mustInit(c *Config) {
	m.gitRoot = git.MustGetRootDir()
	m.gitSubs = git.MustGetSubModules()
	m.cache = filecore.NewCache(m.gitRoot)
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

func (m *modeBase) isIn(filename string, begin, end int) (bool, error) {
	f := m.cache.Get(filename)

	lines, err := f.ReadLines()
	if err != nil {
		return false, err
	}

	if begin == 0 && end == 0 {
		return true, nil
	}

	begin--
	if begin < 0 || end > len(lines) {
		// log.Printf(" === %s, %s, %#v", f.Path, f.Name, lines)
		return false, errInvalidRange
	}

	for _, line := range lines[begin:end] {
		if m.isLineIgnored(line) {
			return true, nil
		}
	}

	return false, nil
}

func (m *modeBase) mustWalkFile(file *filecore.File, fn WalkFunc) {
	filename := file.Name

	if !file.IsGo() {
		return
	}

	if m.isGitSub(file.Name) {
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
