/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"sort"

	"github.com/distroy/git-go-tool/core/mathcore"
)

type Coverages []*Coverage

func (s Coverages) Len() int      { return len(s) }
func (s Coverages) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Coverages) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}

func (s Coverages) IsIn(begin, end int) bool {
	idx := s.Index(begin, end)
	// log.Printf(" === array:%v, begin:%d, end:%d, idx:%d", s, begin, end, idx)
	return idx >= 0
}

func (s Coverages) Index(begin, end int) int {
	idx := sort.Search(len(s), func(i int) bool {
		d := s[i]
		return d.EndLine >= begin
	})
	if idx >= len(s) {
		return -1
	}

	d := s[idx]
	if begin > d.EndLine || end < d.BeginLine {
		return -1
	}

	return idx
}

func (s Coverages) GetCount() int {
	count := 0
	for _, c := range s {
		count += c.EndLine - c.BeginLine + 1
	}
	return count
}

func (s Coverages) Add(c *Coverage) Coverages {
	begin := c.BeginLine
	end := c.EndLine

	idx := sort.Search(len(s), func(i int) bool {
		d := s[i]
		return d.EndLine >= begin
	})
	if idx >= len(s) {
		s = append(s, c)
		return s.merge(idx)
	}

	c0 := s[idx].clone()
	if begin > c0.EndLine+1 || end < c0.BeginLine-1 {
		s = append(s, c)
		copy(s[idx+1:], s[idx:])
		s[idx] = c
		return s.merge(idx)
	}

	c0.BeginLine = mathcore.MinInt(c0.BeginLine, begin)
	c0.EndLine = mathcore.MaxInt(c0.EndLine, end)
	s[idx] = c0

	return s.merge(idx)
}

func (s Coverages) merge(idx int) Coverages {
	c := s[idx].clone()

	idxLeft := idx - 1
	for ; idxLeft >= 0; idxLeft-- {
		tmp := s[idxLeft]
		if c.BeginLine > tmp.EndLine+1 || c.EndLine < tmp.BeginLine-1 {
			break
		}
	}
	idxLeft++

	idxRight := idx + 1
	for ; idxRight < len(s); idxRight++ {
		tmp := s[idxRight]
		if c.BeginLine > tmp.EndLine+1 || c.EndLine < tmp.BeginLine-1 {
			break
		}
	}
	idxRight--

	if idxLeft != idx {
		c.BeginLine = mathcore.MinInt(c.BeginLine, s[idxLeft].BeginLine)
		c.EndLine = mathcore.MaxInt(c.EndLine, s[idxLeft].EndLine)
	}

	if idxRight != idx {
		c.BeginLine = mathcore.MinInt(c.BeginLine, s[idxRight].BeginLine)
		c.EndLine = mathcore.MaxInt(c.EndLine, s[idxRight].EndLine)
	}

	copy(s[idxLeft:], s[idxRight:])
	s[idxLeft] = c
	return s[:len(s)-idxRight+idxLeft]
}
