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
	InputNum               int
	OutputNum              int
	InputNumWithoutContext bool
	OutputNumWithoutError  bool
	ContextFirst           bool
	ErrorLast              bool
	ContextErrorMatch      bool
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

		if err := c.checkFuncParams(x, fn.Type); err != nil {
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
			err = c.checkFuncParams(x, nn.Type)

		case *ast.FuncDecl:
			err = c.checkFuncParams(x, nn.Type)
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

		if len(param.Names) == 0 {
			res = append(res, &funcParamInfo{
				Name: "",
				Type: typ,
			})
			continue
		}

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

func (c funcParamsChecker) checkFuncParams(x *Context, fn *ast.FuncType) error {
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
	// log.Printf(" === file:%s:%d", x.Name, pos.Line)
	// log.Printf(" === file:%s:%d, ins:%s, outs:%s", x.Name, pos.Line, mustJsonMarshal(ins), mustJsonMarshal(outs))

	ctxIdx, ctx := c.indexParamByTypeName(ins, "context")
	if !c.isContextFirst(ins, ctxIdx) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the input parameter '%s' should be the first", ctx.Type.String),
		})
	}

	errIdx, err := c.indexParamByTypeName(outs, "error")
	if !c.isErrorLast(outs, errIdx) {
		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the output parameter '%s' should not be more last", err.Type.String),
		})
	}

	if !c.isContextErrorMatch(x, ctx, err) {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the context '%s' is not matched the error '%s'",
				ctx.Type.String, err.Type.String),
		})
	}

	if !c.isInNumValidWithoutContext(x, ins, ctx) {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters without '%s' should not be more than %d, there are %d",
				ctx.Type.String, inLimit, inNum-1),
		})

	} else if !c.isInNumValid(x, ins, ctx) {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters should not be more than %d, there are %d",
				inLimit, inNum),
		})
	}

	if !c.isOutNumValidWithoutError(x, outs, err) {
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of output parameters without '%s' should not be more than %d, there are %d",
				err.Type.String, outLimit, outNum-1),
		})

	} else if !c.isOutNumValid(x, outs, err) {
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
		typ := v.Type
		// log.Printf(" === %s", mustJsonMarshal(typ))
		isSpecial := typ.IsEllipsis || typ.IsFunc || typ.IsSlice
		if !isSpecial && strings.EqualFold(typ.Name, typeName) {
			// log.Printf(" === return %d %s", i, mustJsonMarshal(typ))
			return i, v
		}
	}
	return -1, nil
}

func (c funcParamsChecker) isContextFirst(params []*funcParamInfo, idx int) bool {
	if !c.cfg.ContextFirst {
		return true
	}

	if idx < 0 {
		return true
	}

	return idx == 0
}

func (c funcParamsChecker) isErrorLast(params []*funcParamInfo, idx int) bool {
	if !c.cfg.ErrorLast {
		return true
	}

	if idx < 0 {
		return true
	}

	return idx == len(params)-1
}

func (c funcParamsChecker) isInNumValidWithoutContext(x *Context, params []*funcParamInfo, ctx *funcParamInfo) bool {
	limit := c.cfg.InputNum
	num := len(params)

	if limit <= 0 {
		return true
	}

	if !c.cfg.InputNumWithoutContext || ctx == nil {
		return true
	}

	num--
	return num <= limit
}

func (c funcParamsChecker) isInNumValid(x *Context, params []*funcParamInfo, ctx *funcParamInfo) bool {
	limit := c.cfg.InputNum
	num := len(params)

	if limit <= 0 {
		return true
	}

	if c.cfg.InputNumWithoutContext && ctx != nil {
		return true
	}

	return num <= limit
}

func (c funcParamsChecker) isOutNumValidWithoutError(x *Context, params []*funcParamInfo, err *funcParamInfo) bool {
	limit := c.cfg.OutputNum
	num := len(params)

	if limit <= 0 {
		return true
	}

	if !c.cfg.OutputNumWithoutError || err == nil {
		return true
	}

	num--
	return num <= limit
}

func (c funcParamsChecker) isOutNumValid(x *Context, params []*funcParamInfo, err *funcParamInfo) bool {
	limit := c.cfg.OutputNum
	num := len(params)

	if limit <= 0 {
		return true
	}

	if c.cfg.OutputNumWithoutError && err != nil {
		return true
	}

	return num <= limit
}

func (c funcParamsChecker) isContextErrorMatch(x *Context, ctx, err *funcParamInfo) bool {
	// if ctx != nil {
	// 	log.Printf(" === ctx:%s, %v", ctx.Type.String, c.isStdContext(x, ctx))
	// }
	// if err != nil {
	// 	log.Printf(" === err:%s, %v", err.Type.String, c.isStdError(x, err))
	// }

	if ctx == nil || err == nil {
		return true
	}
	if !c.cfg.ContextErrorMatch {
		return true
	}

	return c.isStdContext(x, ctx) == c.isStdError(x, err)
}

func (c funcParamsChecker) isStdContext(x *Context, ctx *funcParamInfo) bool {
	typ := ctx.Type
	if typ.IsPointer {
		return false
	}

	if typ.Package == "" {
		return false
	}

	f := x.MustParse()
	imps := convertImports(x, f.Imports)
	if len(imps) == 0 {
		return true
	}

	for _, imp := range imps {
		if imp.Name == typ.Package || imp.Path == typ.Package {
			return imp.Path == "context"
		}
	}
	return true
}

func (c funcParamsChecker) isStdError(x *Context, err *funcParamInfo) bool {
	typ := err.Type
	return !typ.IsPointer && typ.Package == "" && typ.Name == "error"
}
