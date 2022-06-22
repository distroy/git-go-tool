/*
 * Copyright (C) distroy
 */

package goformat

type Level int

const (
	LevelError Level = iota
	LevelWarning
	LevelInfo
)

func (l Level) String() string {
	switch l {
	case LevelError:
		return "error"
	case LevelWarning:
		return "warning"
	case LevelInfo:
		return "info"
	}
	return "error"
}

type Issue struct {
	Filename    string
	BeginLine   int
	EndLine     int
	Level       Level
	Description string
	// Code        []string
}

type Count struct {
	Error   int
	Warning int
	Info    int
}

type sortedIssues []*Issue

func (s sortedIssues) Len() int      { return len(s) }
func (s sortedIssues) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortedIssues) Less(i, j int) bool {
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
}
