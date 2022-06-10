/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/distroy/git-go-tool/core/iocore"
)

var (
	errInvalidRange = fmt.Errorf("invalid file range")
)

type file struct {
	Filename string
	Lines    []string
}

type cache struct {
	root  string
	files map[string]*file
}

func newCache(rootPath string) *cache {
	return &cache{
		root:  rootPath,
		files: make(map[string]*file),
	}
}

// begin: start form 1
// [begin, end]
func (c *cache) CheckFileRange(filename string, begin, end int, check func(line string) bool) (bool, error) {
	f, err := c.GetFile(filename)
	if err != nil {
		return false, err
	}

	begin--
	if begin < 0 || end > len(f.Lines) {
		return false, errInvalidRange
	}

	for _, line := range f.Lines[begin:end] {
		if check(line) {
			return true, nil
		}
	}

	return false, nil
}

func (c *cache) MustGetFile(filename string) *file {
	f, err := c.GetFile(filename)
	if err != nil {
		log.Fatalf("read file fail. file:%s, err:%v", filename, err)
	}
	return f
}

func (c *cache) GetFile(filename string) (*file, error) {
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

func (c *cache) DelFile(filename string) *file {
	f := c.files[filename]
	delete(c.files, filename)
	return f
}

func (c *cache) loadFile(filename string) (*file, error) {
	filePath := path.Join(c.root, filename)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := iocore.NewLineReader(f)

	lines, err := r.ReadAllLineStrings()
	if err != nil {
		return nil, err
	}

	return &file{
		Filename: filename,
		Lines:    lines,
	}, nil
}
