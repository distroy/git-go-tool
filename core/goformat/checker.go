/*
 * Copyright (C) distroy
 */

package goformat

import (
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
)

type Checker interface {
	Check(f *filecore.File) []*Issue
}

type checkers []Checker

func (c checkers) Check(f *filecore.File) []*Issue {
	res := make([]*Issue, 0, 16)
	for _, checker := range c {
		r := checker.Check(f)
		res = append(res, r...)
	}
	return res
}

func Checkers(args ...Checker) Checker {
	n := filter.FilterSlice(args, func(v Checker) bool {
		if v == nil {
			return false
		}

		switch vv := v.(type) {
		case checkerNil, *checkerNil:
			return false

		case checkers:
			return len(vv) > 0
		}

		return true
	})

	args = args[:n]
	return checkers(args)
}
