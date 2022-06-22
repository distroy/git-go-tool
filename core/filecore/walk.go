/*
 * Copyright (C) distroy
 */

package filecore

import (
	"os"
	"path/filepath"
)

func WalkFiles(dirPath string, walkFunc func(file *File) error) error {
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
