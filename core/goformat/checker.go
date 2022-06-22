/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

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

func AddChecker(args ...Checker) Checker {
	if len(args) == 0 {
		return checkers(nil)
	}

	size := getCheckersSize(args...)
	if size == len(args) {
		return checkers(args)
	}

	res := make(checkers, 0, size)
	for _, c := range args {
		res = appendChecker(res, c)
	}

	return res
}

func appendChecker(res checkers, c Checker) checkers {
	if c == nil {
		return res
	}
	if cc, ok := c.(checkers); ok {
		res = append(res, cc...)
	} else {
		res = append(res, c)
	}
	return res
}

func getCheckerSize(checker Checker) int {
	if checker == nil {
		return 0
	}
	if v, ok := checker.(checkers); ok {
		return len(v)
	}
	return 1
}

func getCheckersSize(args ...Checker) int {
	size := 0
	for _, c := range args {
		size += getCheckerSize(c)
	}
	return size
}
