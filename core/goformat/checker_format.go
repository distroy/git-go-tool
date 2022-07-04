/*
 * Copyright (C) distroy
 */

package goformat

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/distroy/git-go-tool/core/mathcore"
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
	buffer.Grow(c.getBufferSize(x, data))
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

func (c formatChecker) getBufferSize(x *Context, data []byte) int {
	size := len(data)
	size = size + size/10
	size = mathcore.MaxInt(size, 4096)
	return size
}
