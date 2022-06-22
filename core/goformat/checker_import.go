/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
)

type importInfo struct {
	Name   string
	Path   string
	Line   int
	StdLib bool
}

func ImportChecker() Checker {
	return importChecker{}
}

type importChecker struct {
}

func (c importChecker) Check(f *filecore.File) []*Issue {
	file := f.MustParse()
	if len(file.Imports) <= 0 {
		return nil
	}

	stds, others := c.converterImport(f, file.Imports)
	return c.checkImport(f, stds, others)
}

func (c importChecker) checkImport(f *filecore.File, stds, others []*importInfo) []*Issue {
	res := make([]*Issue, 0, 8)

	if len(stds)+len(others) <= 1 {
		return nil
	}

	if imps := stds; c.hasBlankLine(imps) {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   imps[0].Line,
			EndLine:     imps[len(imps)-1].Line,
			Level:       LevelError,
			Description: fmt.Sprintf("must not have blank line in std imports"),
		})
	}

	if imps := others; c.hasBlankLine(imps) {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   imps[0].Line,
			EndLine:     imps[len(imps)-1].Line,
			Level:       LevelError,
			Description: fmt.Sprintf("must not have blank line in other imports"),
		})
	}

	if n := c.getStdLibCount(stds); n < len(stds) {
		imps := stds
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   imps[0].Line,
			EndLine:     imps[len(imps)-1].Line,
			Level:       LevelError,
			Description: fmt.Sprintf("must not have blank line in other imports"),
		})

	} else if n := c.getStdLibCount(others); n > 0 {
	}

	for _, imp := range stds {
		log.Printf(" === line:%d, name:%s, path:%s", imp.Line, imp.Name, imp.Path)
	}
	log.Printf("")
	for _, imp := range others {
		log.Printf(" === line:%d, name:%s, path:%s", imp.Line, imp.Name, imp.Path)
	}

	return res
}

func (c importChecker) getStdLibCount(imps []*importInfo) int {
	count := 0
	for _, imp := range imps {
		if imp.StdLib {
			count++
		}
	}
	return count
}

func (c importChecker) hasBlankLine(imps []*importInfo) bool {
	if len(imps) < 2 {
		return false
	}

	lastImp := imps[0]
	for _, imp := range imps[1:] {
		if imp.Line > lastImp.Line+1 {
			return true
		}
		lastImp = imp
	}

	return false
}

func (c importChecker) isStdLibPath(path string) bool {
	idx0 := strings.Index(path, "/")
	if idx0 < 0 {
		return true
	}

	idx1 := strings.Index(path[:idx0], ".")
	if idx1 < 0 {
		return true
	}

	return false
}

func (c importChecker) converterImport(f *filecore.File, imps []*ast.ImportSpec) (stds, others []*importInfo) {
	buf := make([]*importInfo, 0, len(imps))
	for _, imp := range imps {
		v := &importInfo{
			Line: f.Position(imp.Pos()).Line,
		}

		if imp.Name != nil {
			v.Name = imp.Name.Name
		}

		v.Path, _ = strconv.Unquote(imp.Path.Value)
		v.StdLib = c.isStdLibPath(v.Path)

		buf = append(buf, v)
	}

	// sort.Slice(buf, func(i, j int) bool { return buf[i].Line < buf[j].Line })
	n := filter.FilterSlice(buf, func(v *importInfo) bool { return v.StdLib })
	stds = buf[:n]
	others = buf[n:]

	sort.Slice(stds, func(i, j int) bool { return buf[i].Line < buf[j].Line })
	sort.Slice(others, func(i, j int) bool { return buf[i].Line < buf[j].Line })

	return stds, others
}
