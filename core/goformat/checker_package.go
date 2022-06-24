/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)

func PackageChecker(enable bool) Checker {
	if !enable {
		return checkerNil{}
	}
	return packageChecker{}
}

type packageChecker struct{}

func (c packageChecker) Check(f *filecore.File) []*Issue {
	res := make([]*Issue, 0, 8)

	file := f.MustParse()

	res = c.checkPackageName(res, f, file)

	return res
}

func (c packageChecker) checkPackageName(res []*Issue, f *filecore.File, file *ast.File) []*Issue {
	pkg := file.Name
	pkgPos := f.Position(pkg.Pos())

	name := pkg.Name
	if strings.Contains(name, "_") {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   pkgPos.Line,
			EndLine:     pkgPos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("do not use the underscore in package name '%s'", name),
		})

	} else if !strcore.IsLower(name) {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   pkgPos.Line,
			EndLine:     pkgPos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("do not use capital letters in package name '%s'", name),
		})
	}

	return res
}
