/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

func ImportChecker() Checker {
	return importChecker{}
}

type importChecker struct {
}

func (c importChecker) Check(f *filecore.File) []*Issue {
	file := f.MustParse()
	for _, imp := range file.Imports {

	}
}
