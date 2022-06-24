/*
 * Copyright (C) distroy
 */

package goformat

import "github.com/distroy/git-go-tool/core/filecore"

type Context struct {
	*filecore.File
	issues []*Issue
}

func NewContext(f *filecore.File) *Context {
	return &Context{
		File:   f,
		issues: make([]*Issue, 0, 16),
	}
}

func (c *Context) Issues() []*Issue {
	return c.issues
}

func (c *Context) AddIssue(issue *Issue) {
	c.issues = append(c.issues, issue)
}
