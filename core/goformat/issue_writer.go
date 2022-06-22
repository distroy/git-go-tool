/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"io"
	"sort"
)

type IssueOnlyWriter interface {
	Write(issue *Issue)
}

type IssueWriter interface {
	IssueOnlyWriter

	WriteIssues(issues []*Issue)
	Count() Count
}

func WrapWriter(w IssueOnlyWriter) IssueWriter {
	if ww, ok := w.(IssueWriter); ok && ww != nil {
		return ww
	}

	return &writerWrapper{
		writer: w,
		count:  Count{},
	}
}

type writerWrapper struct {
	writer IssueOnlyWriter
	count  Count
}

func (w *writerWrapper) Count() Count { return w.count }

func (w *writerWrapper) Write(issues *Issue) {
	w.writer.Write(issues)

	switch issues.Level {
	default:
		fallthrough
	case LevelError:
		w.count.Error++

	case LevelWarning:
		w.count.Warning++

	case LevelInfo:
		w.count.Info++
	}
}

func (w *writerWrapper) WriteIssues(issues []*Issue) {
	sort.Sort(sortedIssues(issues))
	for _, v := range issues {
		w.Write(v)
	}
}

func NewIssueWriter(w io.Writer) IssueWriter {
	return WrapWriter(&writer{
		writer: w,
	})
}

type writer struct {
	writer io.Writer
}

func (w *writer) Write(issue *Issue) {
	if issue.BeginLine == issue.EndLine {
		fmt.Fprintf(w.writer, "%s %s:%d %s\n", issue.Level.String(), issue.Filename, issue.BeginLine, issue.Description)
	} else {
		fmt.Fprintf(w.writer, "%s %s:%d,%d %s\n", issue.Level.String(), issue.Filename, issue.BeginLine, issue.EndLine, issue.Description)
	}
}
