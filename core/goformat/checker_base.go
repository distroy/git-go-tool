/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

type checkerBase struct{}

func (c checkerBase) getLinesByRange(f *filecore.File, beginLine, endLine int) []string {
	lines := f.MustReadLines()
	if beginLine <= 0 || endLine <= 0 || endLine > len(lines) {
		return lines
	}

	beginLine--
	return lines[beginLine:endLine]
}
