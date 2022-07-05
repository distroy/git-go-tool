/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

type Context struct {
	*filecore.File

	cache  *Cache
	issues []*Issue
}

func NewContext(f *filecore.File) *Context {
	return &Context{
		File:   f,
		issues: make([]*Issue, 0, 16),
	}
}

func (x *Context) Issues() []*Issue {
	return x.issues
}

func (x *Context) AddIssue(issue *Issue) {
	x.issues = append(x.issues, issue)
}
