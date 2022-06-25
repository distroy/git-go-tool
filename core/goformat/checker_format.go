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
	return c.checkData(x, data)
}

func (c formatChecker) checkData(x *Context, data []byte) Error {
	fmtData, err := format.Source(data)
	if err != nil {
		panic(fmt.Errorf("format file fail. file:%s, err:%v", x.Name, err))
	}

	if !bytes.Equal(data, fmtData) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("source should be formated"),
		})
	}

	return nil
}
