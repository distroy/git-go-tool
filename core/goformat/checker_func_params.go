/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"strings"
)

type FuncParamsConfig struct {
	InputNum     int
	OutputNum    int
	ContextFirst bool
	ErrorLast    bool
}

func FuncParamsChecker(cfg *FuncParamsConfig) Checker {
	return funcParamsChecker{cfg: cfg}
}

type funcParamInfo struct {
	Name string
	Type *typeInfo
}

type iAstFunc interface {
}

type funcParamsChecker struct {
	cfg *FuncParamsConfig
}

func (c funcParamsChecker) Check(x *Context) error {
	file := x.MustParse()

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if err := c.checkFuncDecl(x, fn.Type); err != nil {
			return err
		}

		if err := c.walkFunc(x, fn); err != nil {
			return err
		}
	}

	return nil
}

func (c funcParamsChecker) walkFunc(x *Context, fn *ast.FuncDecl) error {
	if fn.Body == nil {
		return nil
	}

	var err error
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// log.Printf(" === %T %#v", n, n)

		switch nn := n.(type) {
		case *ast.FuncLit:
			err = c.checkFuncDecl(x, nn.Type)

		case *ast.FuncDecl:
			err = c.checkFuncDecl(x, nn.Type)
		}

		return err == nil
	})
	return err
}

func (c funcParamsChecker) convertParams(x *Context, params *ast.FieldList) []*funcParamInfo {
	if params == nil {
		return nil
	}

	n := params.NumFields()
	res := make([]*funcParamInfo, 0, n)

	for _, param := range params.List {
		typ := getTypeInfo(x.File, param.Type)
		// log.Printf(" === %T: %v, %#v", param.Type, param.Type, typ)

		for _, v := range param.Names {
			name := ""
			if v != nil {
				name = v.Name
			}

			res = append(res, &funcParamInfo{
				Name: name,
				Type: typ,
			})
		}

	}

	return res
}

func (c funcParamsChecker) checkFuncDecl(x *Context, fn *ast.FuncType) error {
	inLimit := c.cfg.InputNum
	outLimit := c.cfg.OutputNum

	ins := c.convertParams(x, fn.Params)
	outs := c.convertParams(x, fn.Results)

	inNum := len(ins)
	outNum := len(outs)
	if inNum == 0 && outNum == 0 {
		return nil
	}

	pos := x.Position(fn.Pos())
	// log.Printf(" === file:%s:%d, func:%s", f.Name, pos.Line, fn.Name.Name)

	ctxIdx, ctx := c.indexParamByTypeName(ins, "context")
	if c.isContextFirst(ins, ctxIdx) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the input parameter '%s' should be the first", ctx.Type.String),
		})
	}

	if inLimit > 0 && ctx != nil && inNum > inLimit+1 {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters without '%s' should not be more than %d, there are %d",
				ctx.Type.String, inLimit, inNum-1),
		})

	} else if inLimit > 0 && ctx == nil && inNum > inLimit {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters should not be more than %d, there are %d",
				inLimit, inNum),
		})
	}

	errIdx, err := c.indexParamByTypeName(outs, "error")
	if c.isErrorLast(outs, errIdx) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the output parameter '%s' should not be more last", err.Type.String),
		})
	}

	if outLimit > 0 && err != nil && outNum > outLimit+1 {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of output parameters without '%s' should not be more than %d, there are %d",
				err.Type.String, outLimit, outNum-1),
		})

	} else if outLimit > 0 && err == nil && outNum > outLimit {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of output parameters should not be more than %d, there are %d",
				outLimit, outNum),
		})
	}

	return nil
}

func (c funcParamsChecker) indexParamByTypeName(params []*funcParamInfo, typeName string) (int, *funcParamInfo) {
	for i, v := range params {
		if strings.EqualFold(v.Type.Name, typeName) {
			return i, v
		}
	}
	return -1, nil
}

func (c funcParamsChecker) isContextFirst(params []*funcParamInfo, idx int) bool {
	if !c.cfg.ContextFirst {
		return false
	}

	return idx > 0
}

func (c funcParamsChecker) isErrorLast(params []*funcParamInfo, idx int) bool {
	if !c.cfg.ErrorLast {
		return false
	}

	if idx < 0 {
		return false
	}

	return idx < len(params)-1
}
