/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"sort"
	"strconv"
	"strings"

	"github.com/distroy/git-go-tool/core/filtercore"
	"github.com/distroy/git-go-tool/core/mathcore"
)

type importInfo struct {
	Name   string
	Path   string
	Line   int
	StdLib bool
}

func ImportChecker(enable bool) Checker {
	if !enable {
		return checkerNil{}
	}

	return importChecker{}
}

type importChecker struct{}

func (c importChecker) Check(x *Context) Error {
	f := x.MustParse()
	if len(f.Imports) <= 0 {
		return nil
	}

	imps := convertImports(x, f.Imports)
	return c.checkImport(x, imps)
}

func (c importChecker) checkImport(x *Context, imps []*importInfo) Error {
	for _, imp := range imps {
		// log.Printf(" === line:%d, name:%s, path:%s", imp.Line, imp.Name, imp.Path)
		if imp.Name == "." {
			x.AddIssue(&Issue{
				Filename:    x.Name,
				BeginLine:   imp.Line,
				EndLine:     imp.Line,
				Level:       LevelError,
				Description: fmt.Sprintf("do not use the dot import"),
			})
		}
	}

	if len(imps) <= 1 {
		return nil
	}

	begin, end := c.getImportRange(imps)

	n := filtercore.FilterSlice(imps, func(v *importInfo) bool { return v.StdLib })
	stds := imps[:n]
	others := imps[n:]

	sort.Slice(stds, func(i, j int) bool { return stds[i].Line < stds[j].Line })
	sort.Slice(others, func(i, j int) bool { return others[i].Line < others[j].Line })

	if !c.isGroupedAndOrdered(stds, others) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   begin,
			EndLine:     end,
			Level:       LevelError,
			Description: fmt.Sprintf("imports should be grouped and ordered by standards and others"),
		})
	}

	return nil
}

func (c importChecker) isGroupedAndOrdered(stds, others []*importInfo) bool {
	if c.hasBlankLine(stds) {
		return false
	}
	if c.hasBlankLine(others) {
		return false
	}
	if c.getStdLibCount(stds) < len(stds) {
		return false
	}
	if c.getStdLibCount(others) > 0 {
		return false
	}

	if len(stds) > 0 && len(others) > 0 {
		if stds[len(stds)-1].Line > others[0].Line {
			return false
		}
		if others[0].Line-stds[len(stds)-1].Line != 2 {
			return false
		}
	}

	return true
}

func (c importChecker) getImportRange(imps []*importInfo) (begin, end int) {
	begin = imps[0].Line
	end = imps[0].Line
	for _, imp := range imps {
		begin = mathcore.MinInt(begin, imp.Line)
		end = mathcore.MaxInt(end, imp.Line)
	}
	return
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

func isStdLibPath(path string) bool {
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

func convertImports(ctx *Context, imps []*ast.ImportSpec) []*importInfo {
	buf := make([]*importInfo, 0, len(imps))
	for _, imp := range imps {
		v := &importInfo{
			Line: ctx.Position(imp.Pos()).Line,
		}

		if imp.Name != nil {
			v.Name = imp.Name.Name
		}

		v.Path, _ = strconv.Unquote(imp.Path.Value)
		v.StdLib = isStdLibPath(v.Path)

		buf = append(buf, v)
	}

	return buf
}
