/*
 * Copyright (C) distroy
 */

package goformat

import "go/token"

type File struct {
	Filename string
	fset     *token.FileSet
	file     *token.File
}
