/*
 * Copyright (C) distroy
 */

package goformat

import (
	"encoding/json"

	"github.com/distroy/git-go-tool/core/filtercore"
	"github.com/distroy/git-go-tool/core/strcore"
)

type checkers []Checker

func (c checkers) Check(ctx *Context) Error {
	for _, checker := range c {
		err := checker.Check(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func Checkers(args ...Checker) Checker {
	n := filtercore.FilterSlice(args, func(v Checker) bool {
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

func mustJsonMarshal(v interface{}) string {
	d, _ := json.Marshal(v)
	return strcore.BytesToStrUnsafe(d)
}
