/*
 * Copyright (C) distroy
 */

package goformat

import (
	"go/format"

	"github.com/distroy/git-go-tool/core/filecore"
)

type formatChecker struct {
}

func (c formatChecker) Check(f *filecore.File) []*Issue {
	data := f.MustRead()

	fmtData, err := format.Source(data)
}
