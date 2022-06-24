/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/strcore"
)

func PackageChecker(enable bool) Checker {
	if !enable {
		return checkerNil{}
	}
	return packageChecker{}
}

type packageChecker struct{}

func (c packageChecker) Check(x *Context) error {

	file := x.MustParse()

	return c.checkPackageName(x, file)
}

func (c packageChecker) checkPackageName(x *Context, file *ast.File) error {
	pkg := file.Name
	pkgPos := x.Position(pkg.Pos())

	name := pkg.Name
	if strings.Contains(name, "_") {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pkgPos.Line,
			EndLine:     pkgPos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("do not use the underscore in package name '%s'", name),
		})

	} else if !strcore.IsLower(name) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pkgPos.Line,
			EndLine:     pkgPos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("do not use capital letters in package name '%s'", name),
		})
	}

	return nil
}
