/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"

	"github.com/distroy/git-go-tool/core/filecore"
)

func FileLineChecker(fileLine int) Checker {
	if fileLine <= 0 {
		return checkerNil{}
	}

	return fileLineChecker{
		fileLine: fileLine,
	}
}

type fileLineChecker struct {
	fileLine int
}

func (c fileLineChecker) Check(f *filecore.File) []*Issue {
	if c.fileLine <= 0 {
		return nil
	}

	lines := f.MustReadLines()
	return c.checkLines(f, lines)
}

func (c fileLineChecker) checkLines(f *filecore.File, lines []string) []*Issue {
	if f.IsGoTest() {
		return nil
	}

	res := make([]*Issue, 0, 1)

	limit := c.fileLine
	if len(lines) > limit {
		res = append(res, &Issue{
			Filename:    f.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("file lines(%d) is more than %d, should split the file", len(lines), limit),
		})
	}

	return res
}
