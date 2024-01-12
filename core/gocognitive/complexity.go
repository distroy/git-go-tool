/*
 * Copyright (C) distroy
 */

package gocognitive

import "fmt"

// Complexity is statistic of the complexity.
type Complexity struct {
	PkgName    string `json:"pkg_name"`
	FuncName   string `json:"func_name"`
	Filename   string `json:"file_name"`
	Complexity int    `json:"complexity"`
	BeginLine  int    `json:"begin_line"`
	EndLine    int    `json:"end_line"`
}

func (s *Complexity) String() string {
	filePos := fmt.Sprintf("%s:%d,%d", s.Filename, s.BeginLine, s.EndLine)
	return fmt.Sprintf("%d %s %s %s", s.Complexity, s.PkgName, s.FuncName, filePos)
}

type Complexities []*Complexity

func (s Complexities) Len() int      { return len(s) }
func (s Complexities) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Complexities) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Complexity != b.Complexity {
		return a.Complexity > b.Complexity
	}
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}
