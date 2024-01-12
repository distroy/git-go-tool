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
	Issues() []*Issue
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
	issues []*Issue
	count  Count
}

func (w *writerWrapper) Count() Count     { return w.count }
func (w *writerWrapper) Issues() []*Issue { return w.issues }

func (w *writerWrapper) Write(issues *Issue) {
	w.writer.Write(issues)
	w.issues = append(w.issues, issues)

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
	sort.Slice(issues, func(i, j int) bool {
		s := issues
		a, b := s[i], s[j]
		if a.Filename != b.Filename {
			return a.Filename < b.Filename
		}
		if a.BeginLine != b.BeginLine {
			return a.BeginLine < b.BeginLine
		}
		if a.EndLine != b.EndLine {
			return a.EndLine < b.EndLine
		}
		return a.Level > b.Level
	})

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
	fmt.Fprintf(w.writer, "%s [%s] %s\n", w.fileAndLine(issue), issue.Level.String(), issue.Description)
}

func (w *writer) fileAndLine(issue *Issue) string {
	if issue.BeginLine <= 0 || issue.EndLine <= 0 {
		return issue.Filename
	}
	if issue.BeginLine == issue.EndLine {
		return fmt.Sprintf("%s:%d", issue.Filename, issue.BeginLine)
	}
	return fmt.Sprintf("%s:%d,%d", issue.Filename, issue.BeginLine, issue.EndLine)
}
