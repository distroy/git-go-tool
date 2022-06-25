/*
 * Copyright (C) distroy
 */

package git

import "sort"

type Differents []Different

func (s Differents) Len() int      { return len(s) }
func (s Differents) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Differents) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}

func (s Differents) IsIn(begin, end int) bool {
	idx := sort.Search(len(s), func(i int) bool {
		d := s[i]
		return d.EndLine >= begin
	})
	if idx >= len(s) {
		return false
	}

	d := s[idx]
	if begin > d.EndLine || end < d.BeginLine {
		return false
	}

	return true
}

type Files map[string]Differents

func NewFileDifferents(slice Differents) Files {
	sort.Sort(slice)

	m := make(map[string]Differents)
	lastIdx := 0
	for i, v1 := range slice {
		if v1.Filename == slice[lastIdx].Filename {
			continue
		}

		v0 := slice[lastIdx]
		m[v0.Filename] = slice[lastIdx:i]
		lastIdx = i
	}

	if lastIdx < len(slice) {
		v0 := slice[lastIdx]
		m[v0.Filename] = slice[lastIdx:]
	}

	return m
}

// if begin == 0 && end == 0, check the whole file
func (m Files) IsIn(file string, begin, end int) bool {
	s := m[file]

	if len(s) == 0 {
		return false
	}

	if begin == 0 && end == 0 {
		return true
	}

	return s.IsIn(begin, end)
}
