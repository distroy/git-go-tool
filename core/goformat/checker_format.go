/*
 * Copyright (C) distroy
 */

package goformat

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/distroy/git-go-tool/core/filecore"
)

type formatChecker struct {
}

func (c formatChecker) Check(f *filecore.File) []*Issue {
	data := f.MustRead()
	return c.checkData(f, data)
}

func (c formatChecker) checkData(f *filecore.File, data []byte) []*Issue {
	res := make([]*Issue, 0, 1)

	fmtData, err := format.Source(data)
	if err != nil {
		panic(fmt.Errorf("format file fail. file:%s, err:%v", f.Name, err))
	}

	if !bytes.Equal(data, fmtData) {
		res = append(res, &Issue{
			Filename:    f.Name,
			Level:       LevelError,
			Description: fmt.Sprintf("source must be formated"),
		})
	}

	return res
}
