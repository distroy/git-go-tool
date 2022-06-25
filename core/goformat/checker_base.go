/*
 * Copyright (C) distroy
 */

package goformat

type checkerNil struct{}

func (c checkerNil) Check(x *Context) Error {
	return nil
}
