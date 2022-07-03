/*
 * Copyright (C) distroy
 */

package goformat

import (
	"bytes"
	"fmt"
	"go/format"
)

func FormatChecker(enable bool) Checker {
	if !enable {
		return checkerNil{}
	}

	return formatChecker{}
}

type formatChecker struct {
}

func (c formatChecker) Check(x *Context) Error {
	data := x.MustRead()
	file := x.MustParse()
	fset := x.FileSet()

	buffer := &bytes.Buffer{}
	buffer.Grow(len(data))
	if err := format.Node(buffer, fset, file); err != nil {
		panic(fmt.Errorf("format file fail. file:%s, err:%v", x.Name, err))
	}

	fmtData := buffer.Bytes()
	if !bytes.Equal(data, fmtData) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("source should be formated"),
		})
	}

	return nil
}
