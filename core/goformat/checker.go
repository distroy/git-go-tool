/*
 * Copyright (C) distroy
 */

package goformat

import (
	"github.com/distroy/git-go-tool/core/filter"
)

type Checker interface {
	Check(x *Context) error
}

type checkers []Checker

func (c checkers) Check(ctx *Context) error {
	for _, checker := range c {
		err := checker.Check(ctx)
		if err != nil {
			return err
		}
	}
	return nil
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
