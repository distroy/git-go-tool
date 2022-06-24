/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

type checkerNil struct{}

func (c checkerNil) Check(f *filecore.File) []*Issue { return nil }
