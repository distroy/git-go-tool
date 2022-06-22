/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"

	"github.com/distroy/git-go-tool/core/filecore"
)

func FileLineChecker(fileLine int) Checker {
	return fileLineChecker{fileLine: fileLine}
}

type fileLineChecker struct {
	fileLine int
}

func (c fileLineChecker) Check(f *filecore.File) []*Issue {
	lines := f.MustReadLines()
	return c.checkLines(f, lines)
}

func (c fileLineChecker) checkLines(f *filecore.File, lines []string) []*Issue {
	res := make([]*Issue, 0, 1)
	limit := c.fileLine

	if len(lines) > limit {
		res = append(res, &Issue{
			Filename:    f.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("file lines(%d) is more than %d, must split the file", len(lines), limit),
		})
	}

	return res
}
