/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
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

func (c fileLineChecker) Check(x *Context) error {
	if c.fileLine <= 0 {
		return nil
	}

	lines := x.MustReadLines()
	return c.checkLines(x, lines)
}

func (c fileLineChecker) checkLines(x *Context, lines []string) error {
	if x.IsGoTest() {
		return nil
	}

	limit := c.fileLine
	if len(lines) > limit {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("file lines(%d) is more than %d, should split the file", len(lines), limit),
		})
	}

	return nil
}
