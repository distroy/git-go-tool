/*
 * Copyright (C) distroy
 */

package filelinecache

import (
	"fmt"
	"os"
	"path"

	"github.com/distroy/git-go-tool/core/iocore"
)

var (
	ErrInvalidRange = fmt.Errorf("invalid file range")
)

type File struct {
	Filename string
	Lines    []string
}

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

// [begin, end)
func (c *Cache) CheckFileRange(filename string, begin, end int, check func(line string) bool) (bool, error) {
	f, err := c.GetFile(filename)
	if err != nil {
		return false, err
	}

	if begin < 0 || end > len(f.Lines) {
		return false, ErrInvalidRange
	}

	for _, line := range f.Lines[begin:end] {
		if check(line) {
			return true, nil
		}
	}
	return false, nil
}

func (c *Cache) GetFile(filename string) (*File, error) {
	f := c.files[filename]
	if f != nil {
		return f, nil
	}

	f, err := c.loadFile(filename)
	if err != nil {
		return nil, err
	}

	c.files[filename] = f

	return f, nil
}

func (c *Cache) loadFile(filename string) (*File, error) {
	filePath := path.Join(c.root, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := iocore.NewLineReader(file)

	lines, err := r.ReadAllLineStrings()
	if err != nil {
		return nil, err
	}

	return &File{
		Filename: filename,
		Lines:    lines,
	}, nil
}
