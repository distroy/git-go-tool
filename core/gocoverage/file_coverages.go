/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"sort"

	"github.com/distroy/git-go-tool/core/mathcore"
)

type Coverages []Coverage

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
	return s.Index(begin, end) >= 0
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

func (s Coverages) GetLineCount() int {
	count := 0
	for _, c := range s {
		count += c.EndLine - c.BeginLine + 1
	}
	return count
}

func (s Coverages) Add(c Coverage) Coverages {
	begin := c.BeginLine
	end := c.EndLine

	idx := sort.Search(len(s), func(i int) bool {
		d := s[i]
		return d.EndLine >= begin
	})
	if idx >= len(s) {
		return append(s, c)
	}

	c0 := s[idx]
	if begin > c0.EndLine || end < c0.BeginLine {
		s = append(s, c)
		copy(s[idx+1:], s[idx:])
		s[idx] = c
		return s
	}

	c0.BeginLine = mathcore.MinInt(c0.BeginLine, begin)
	c0.EndLine = mathcore.MaxInt(c0.EndLine, end)
	s[idx] = c0

	return s.merge(idx)
}

func (s Coverages) merge(idx int) Coverages {
	c := s[idx]

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

type FileCoverages struct {
	Filename     string
	Coverages    Coverages
	NonCoverages Coverages
}

func (f *FileCoverages) GetCount() Count {
	return Count{
		Coverages:    f.Coverages.GetLineCount(),
		NonCoverages: f.NonCoverages.GetLineCount(),
	}
}

func (f *FileCoverages) Add(c Coverage, filters ...filter) {
	if c.Count > 0 {
		f.Coverages = f.addToCoverages(f.Coverages, c, filters)
		return
	}

	filters = append(filters, func(file string, lineNo int) bool {
		return f.Coverages.IsIn(lineNo, lineNo)
	})
	f.NonCoverages = f.addToCoverages(f.NonCoverages, c, filters)
}

func (f *FileCoverages) addToCoverages(s Coverages, c Coverage, filters []filter) Coverages {
	for i := c.BeginLine; i <= c.EndLine; i++ {
		if !doFilters(c.Filename, i, filters) {
			continue
		}

		s = s.Add(Coverage{
			Filename:  c.Filename,
			BeginLine: i,
			EndLine:   i,
			Count:     0,
		})
	}
	return s
}

type Files map[string]*FileCoverages

func NewFileCoverages(coverages []Coverage, filters ...filter) Files {
	res := make(Files)
	for _, c := range coverages {
		if c.Count <= 0 {
			continue
		}

		res.Add(c, filters...)
	}

	for _, c := range coverages {
		if c.Count > 0 {
			continue
		}

		res.Add(c, filters...)
	}

	return res
}

func (f Files) Add(c Coverage, filters ...filter) *FileCoverages {
	v := f[c.Filename]
	if v == nil {
		v = &FileCoverages{
			Filename: c.Filename,
		}
		f[c.Filename] = v
	}
	v.Add(c, filters...)
	return v
}

func (f Files) GetCount() Count {
	count := Count{}
	for _, f := range f {
		res := f.GetCount()
		count.Coverages += res.Coverages
		count.NonCoverages += res.NonCoverages
	}
	return count
}
