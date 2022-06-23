/*
 * Copyright (C) distroy
 */

package goformat

import (
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
)

type typeInfo struct {
	String     string
	IsFunc     bool // func
	IsEllipsis bool // ...
	IsPointer  bool // *
	IsSlice    bool // []
	Package    string
	Name       string
}

func getTypeInfo(f *filecore.File, typ ast.Expr) *typeInfo {
	info := &typeInfo{
		String: getTypeName(f, typ),
	}

	fileTypeInfo(info, typ)

	return info
}

func getTypeName(f *filecore.File, typ ast.Expr) string {
	buf := &strings.Builder{}
	f.WriteCode(buf, typ)
	return buf.String()
}

func fileTypeInfo(info *typeInfo, typ ast.Expr) {
	switch tt := typ.(type) {
	case *ast.Ident:
		// log.Printf(" === ident %s", tt.Name)
		info.Name = tt.String()
		break

	case *ast.StarExpr:
		// log.Printf(" === star %T: %v, %s", tt.X, tt.X, tt.X)
		info.IsPointer = true
		fileTypeInfo(info, tt.X)

	case *ast.ArrayType:
		// log.Printf(" === array len: %T: %v, %s", tt.Len, tt.Len, tt.Len)
		// log.Printf(" === array elt: %T: %v, %s", tt.Elt, tt.Elt, tt.Elt)
		info.IsSlice = true
		fileTypeInfo(info, tt.Elt)

	case *ast.Ellipsis:
		// log.Printf(" === ellipsis %T: %v, %s", tt.Elt, tt.Elt, tt.Elt)
		info.IsEllipsis = true
		fileTypeInfo(info, tt.Elt)

	case *ast.SelectorExpr:
		info.Name = tt.Sel.String()
		if xx, _ := tt.X.(*ast.Ident); xx != nil {
			info.Package = xx.Name
		}

		// log.Printf(" === selector %T: %v, %s", tt.X, tt.X, info.Package)
		// log.Printf(" === selector %T: %v, %s", tt.Sel, tt.Sel, info.Name)
	}
}
