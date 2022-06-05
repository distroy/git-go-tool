/*
 * Copyright (C) distroy
 */

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/distroy/git-go-tool/core/gocoverage"
)

type File struct {
	Count gocoverage.Count
	File  *gocoverage.FileCoverages
}

func (f File) GetNonCoverageLineNos() string {
	buffer := &strings.Builder{}

	for i, c := range f.File.NonCoverages {
		if i != 0 {
			fmt.Fprintf(buffer, ", ")
		}

		if c.BeginLine == c.EndLine {
			fmt.Fprintf(buffer, "%d", c.BeginLine)
		} else {
			fmt.Fprintf(buffer, "[%d,%d]", c.BeginLine, c.EndLine)
		}
	}

	return buffer.String()
}

type Files []File

func (s Files) Len() int      { return len(s) }
func (s Files) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Files) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Count.NonCoverages != b.Count.NonCoverages {
		return a.Count.NonCoverages > b.Count.NonCoverages
	}
	return a.Count.Coverages <= b.Count.Coverages
}

func getTopNonCoverageFiles(files gocoverage.Files, top int) Files {
	res := make(Files, 0, len(files))
	for _, f := range files {
		count := f.GetCount()
		if count.NonCoverages == 0 || count.IsZero() {
			continue
		}

		res = append(res, File{
			Count: count,
			File:  f,
		})
	}

	sort.Sort(res)

	if top <= 0 || top >= len(res) {
		return res
	}

	return res[:top]
}
