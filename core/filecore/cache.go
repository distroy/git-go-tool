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
	f := c.files[filename]
	if f == nil {
		filePath := path.Join(c.root, filename)
		f = &File{
			Path: filePath,
			Name: filename,
		}

		c.files[filename] = f
	}

	return f
}

func (c *Cache) Del(filename string) *File {
	f := c.files[filename]
	delete(c.files, filename)
	return f
}

func (c *Cache) WalkFiles(walkFunc func(file *File) error) error {
	return WalkFiles(c.root, func(file *File) error {
		c.files[file.Name] = file
		return walkFunc(file)
	})
}
