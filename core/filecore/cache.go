/*
 * Copyright (C) distroy
 */

package filecore

import (
	"fmt"
	"path"
)

var (
	ErrInvalidRange = fmt.Errorf("invalid file range")
)

type Cache struct {
	root  string
	files map[string]*File
}

func NewCache(rootPath string) *Cache {
	return &Cache{
		root:  rootPath,
		files: make(map[string]*File),
	}
}

func (c *Cache) Get(filename string) *File {
	filePath := path.Join(c.root, filename)

	f := c.files[filePath]
	if f == nil {
		f = &File{
			Path: filePath,
			Name: filename,
		}

		c.files[filePath] = f
	}

	return f
}

func (c *Cache) Del(filename string) *File {
	f := c.files[filename]
	delete(c.files, filename)
	return f
}

func (c *Cache) WalkFiles(walkFunc func(file *File) error) error {
	rootDir := c.root
	return WalkFiles(rootDir, func(file *File) error {
		c.files[file.Name] = file
		return walkFunc(file)
	})
}
