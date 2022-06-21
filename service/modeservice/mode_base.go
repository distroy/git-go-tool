/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/git"
)

type modeBase struct {
	rootDir string
	cache   *cache
}

func (m *modeBase) mustInit(c *Config) {
	m.rootDir = git.GetRootDir()
	m.cache = newCache(m.rootDir)
}

func (m *modeBase) isLineIgnored(line string) bool {
	line = strings.TrimSpace(line)
	// return len(line) == 0
	return len(line) == 0 || line == "}"
}

func (m *modeBase) mustWalkFile(file *filecore.File, fn WalkFunc) {
	filename := file.Name
	path := file.Path

	if !strings.HasSuffix(filename, ".go") {
		return
	}

	cache := m.cache
	defer cache.Del(filename)

	lines, err := file.ReadLines()
	if err != nil {
		log.Fatalf("read file fail. file:%s, err:%v", filename, err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		log.Fatalf("parse file fail. file:%s, err:%v", filename, err)
	}

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

		pos := fset.Position(body.Lbrace)
		end := fset.Position(body.Rbrace)
		for i := pos.Line - 1; i < end.Line; i++ {
			if m.isLineIgnored(lines[i]) {
				continue
			}

			fn(filename, i+1, i+1)
		}
	}
}
