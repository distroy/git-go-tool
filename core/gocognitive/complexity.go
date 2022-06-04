/*
 * Copyright (C) distroy
 */

package gocognitive

import "fmt"

// Complexity is statistic of the complexity.
type Complexity struct {
	PkgName    string
	FuncName   string
	Filename   string
	Complexity int
	BeginLine  int
	EndLine    int
}

func (s Complexity) String() string {
	filePos := fmt.Sprintf("%s:%d,%d", s.Filename, s.BeginLine, s.EndLine)
	return fmt.Sprintf("%d %s %s %s", s.Complexity, s.PkgName, s.FuncName, filePos)
}

type Complexites []Complexity

func (s Complexites) Len() int      { return len(s) }
func (s Complexites) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Complexites) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Complexity != b.Complexity {
		return a.Complexity > b.Complexity
	}
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}
