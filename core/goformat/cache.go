/*
 * Copyright (C) distroy
 */

package goformat

import (
	"go/token"
	"os"
	"path/filepath"

	"github.com/distroy/git-go-tool/core/filecore"
)

type Cache struct {
	cache map[string]*Context
	fset  *token.FileSet
}

func NewCache() *Cache {
	w := &Cache{
		cache: make(map[string]*Context),
		fset:  token.NewFileSet(),
	}

	return w
}

func (w *Cache) MustWalkDir(dirPath string, walkFunc func(x *Context) Error) {
	err := w.WalkDir(dirPath, walkFunc)
	if err != nil {
		panic(err)
	}
}

func (w *Cache) WalkDir(dirPath string, walkFunc func(x *Context) Error) error {
	return filecore.WalkFiles(dirPath, func(f *filecore.File) error {
		return w.walkOneFile(f, walkFunc)
	})
}

func (w *Cache) MustWalkPathes(pathes []string, walkFunc func(x *Context) Error) {
	err := w.WalkPathes(pathes, walkFunc)
	if err != nil {
		panic(err)
	}
}

func (w *Cache) WalkPathes(pathes []string, walkFunc func(x *Context) Error) error {
	for _, path := range pathes {
		if !filecore.IsDir(path) {
			f := &filecore.File{
				Path: path,
				Name: path,
			}

			err := w.walkOneFile(f, walkFunc)
			if err != nil {
				return err
			}
		}

		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			f := &filecore.File{
				Path: path,
				Name: path,
			}

			return w.walkOneFile(f, walkFunc)
		})
	}
	return nil
}

func (w *Cache) walkOneFile(f *filecore.File, walkFunc func(x *Context) error) error {
	if !f.IsGo() {
		return nil
	}

	f.FileSet = w.fset

	x := w.cache[f.Path]
	if x == nil {
		x = NewContext(f)
		w.cache[f.Path] = x
	}

	err := walkFunc(x)
	x.issues = []*Issue{}
	return err
}
