/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

type Checker interface {
	Check(f *filecore.File)
}
