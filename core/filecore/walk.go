/*
 * Copyright (C) distroy
 */

package filecore

import (
	"fmt"
	"os"
	"path/filepath"
)

func MustWalkFiles(dirPath string, walkFunc func(f *File) error) {
	err := WalkFiles(dirPath, walkFunc)
	if err != nil {
		panic(fmt.Errorf("walk dir fail. dir:%s, err%v", dirPath, err))
	}
}

func WalkFiles(dirPath string, walkFunc func(f *File) error) error {
	// fset := token.NewFileSet()
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		filename, _ := filepath.Rel(dirPath, path)
		file := &File{
			Path: path,
			Name: filename,
			// fset: fset,
		}
		return walkFunc(file)
	})
}

func MustWalkPathes(pathes []string, walkFunc func(f *File) error) {
	for _, path := range pathes {
		if err := walkOnePath(path, walkFunc); err != nil {
			panic(fmt.Errorf("walk path fail. path:%s, err%v", path, err))
		}
	}
}

func WalkPathes(pathes []string, walkFunc func(f *File) error) error {
	for _, path := range pathes {
		if err := walkOnePath(path, walkFunc); err != nil {
			return err
		}
	}

	return nil
}

func walkOnePath(path string, walkFunc func(f *File) error) error {
	if !IsDir(path) {
		f := &File{
			Path: path,
			Name: path,
		}

		return walkFunc(f)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		file := &File{
			Path: path,
			Name: path,
		}
		return walkFunc(file)
	})
}
