/*
 * Copyright (C) distroy
 */

package modeservice

import (
	"fmt"

	"github.com/distroy/git-go-tool/core/filecore"
)

var (
	errInvalidRange = fmt.Errorf("invalid file range")
)

type cache struct {
	*filecore.Cache
}

func newCache(rootPath string) *cache {
	return &cache{
		Cache: filecore.NewCache(rootPath),
	}
}

// begin: start form 1
// [begin, end]
func (c *cache) CheckFileRange(filename string, begin, end int, check func(line string) bool) (bool, error) {
	f := c.Get(filename)

	lines, err := f.ReadLines()
	if err != nil {
		return false, err
	}

	begin--
	if begin < 0 || end > len(lines) {
		// log.Printf(" === %s, %s, %#v", f.Path, f.Name, lines)
		return false, errInvalidRange
	}

	for _, line := range lines[begin:end] {
		if check(line) {
			return true, nil
		}
	}

	return false, nil
}
