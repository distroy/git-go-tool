/*
 * Copyright (C) distroy
 */

package iocore

import "io"

var (
	_ io.Writer = discard{}
)

func Discard() discard { return discard{} }

type discard struct{}

func (discard) Write(p []byte) (int, error)       { return len(p), nil }
func (discard) WriteString(s string) (int, error) { return len(s), nil }
