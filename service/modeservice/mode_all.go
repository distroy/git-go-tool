/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/distroy/git-go-tool/core/git"
)

type modeAll struct {
	modeBase

	rootDir string
	cache   *cache
}

func (m *modeAll) mustInit(c *Config) {
	rootDir := git.GetRootDir()
	cache := newCache(m.rootDir)

	m.rootDir = rootDir
	m.cache = cache
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
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			m.mustWalkFile(path, fn)
		}
		return err
	})

	if err != nil {
		log.Fatalf("walk dir fail. dir:%s, err:%v", rootDir, err)
	}
}

func (m *modeAll) mustWalkFile(path string, fn WalkFunc) {
	rootDir := m.rootDir
	cache := m.cache

	filename, _ := filepath.Rel(rootDir, path)
	file := cache.MustGetFile(filename)
	cache.DelFile(filename)

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
			if m.isLineIgnored(file.Lines[i]) {
				continue
			}

			fn(filename, i+1, i+1)
		}
	}
}
