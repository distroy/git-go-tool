/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type orderType int

const (
	orderType_Asc orderType = iota + 1
	orderType_Desc
)

type FuncParamsConfig struct {
	InputNum               int
	OutputNum              int
	NamedOutput            bool
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

func (c funcParamsChecker) Check(x *Context) Error {
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

func (c funcParamsChecker) walkFunc(x *Context, fn *ast.FuncDecl) Error {
	if fn.Body == nil {
		return nil
	}

	if x.IsGoTest() && fn.Name != nil && strings.HasPrefix(fn.Name.Name, "Test") {
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

func (c funcParamsChecker) checkFuncParams(x *Context, fn *ast.FuncType) Error {
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

	ctxIdx, ctx := c.indexParamByTypeName(ins, "context", orderType_Asc)
	if e := c.checkContextFirst(x, pos, ins, ctxIdx); e != nil {
		return e
	}

	errIdx, err := c.indexParamByTypeName(outs, "error", orderType_Desc)
	if e := c.checkErrorLast(x, pos, outs, errIdx); e != nil {
		return e
	}

	if e := c.checkContextErrorMatch(x, pos, ctx, err); e != nil {
		return e
	}

	if e := c.checkNamedOutput(x, pos, outs); e != nil {
		return e
	}

	if e := c.checkInNumValidWithoutContext(x, pos, ins, ctx); e != nil {
		return e
	}

	if e := c.checkInNumValid(x, pos, ins, ctx); e != nil {
		return e
	}

	if e := c.checkOutNumValidWithoutError(x, pos, outs, err); e != nil {
		return e
	}

	if e := c.checkOutNumValid(x, pos, outs, err); e != nil {
		return e
	}

	return nil
}

func (c funcParamsChecker) indexParamByTypeName(params []*funcParamInfo, typeName string, order orderType) (int, *funcParamInfo) {
	for i, l := 0, len(params); i < l; i++ {
		idx := i
		if order == orderType_Desc {
			idx = l - i - 1
		}

		param := params[idx]
		typ := param.Type
		// log.Printf(" === %s", mustJsonMarshal(typ))
		isSpecial := typ.IsEllipsis || typ.IsFunc || typ.IsSlice
		if !isSpecial && strings.Contains(strings.ToLower(typ.Name), typeName) {
			// log.Printf(" === return %d %s", idx, mustJsonMarshal(typ))
			return idx, param
		}
	}
	return -1, nil
}

func (c funcParamsChecker) checkContextFirst(x *Context, pos token.Position, params []*funcParamInfo, idx int) Error {
	if !c.cfg.ContextFirst {
		return nil
	}

	if idx < 0 {
		return nil
	}

	if idx == 0 {
		return nil
	}

	ctx := params[idx]

	x.AddIssue(&Issue{
		Filename:    x.Name,
		BeginLine:   pos.Line,
		EndLine:     pos.Line,
		Level:       LevelError,
		Description: fmt.Sprintf("the input parameter '%s' should be the first", ctx.Type.String),
	})
	return nil
}

func (c funcParamsChecker) checkErrorLast(x *Context, pos token.Position, params []*funcParamInfo, idx int) Error {
	if !c.cfg.ErrorLast {
		return nil
	}

	if idx < 0 {
		return nil
	}

	if idx == len(params)-1 {
		return nil
	}

	err := params[idx]

	x.AddIssue(&Issue{
		Filename:    x.Name,
		BeginLine:   pos.Line,
		EndLine:     pos.Line,
		Level:       LevelError,
		Description: fmt.Sprintf("the output parameter '%s' should not be more last", err.Type.String),
	})
	return nil
}

func (c funcParamsChecker) checkInNumValidWithoutContext(x *Context, pos token.Position, params []*funcParamInfo, ctx *funcParamInfo) Error {
	limit := c.cfg.InputNum
	num := len(params)

	if limit <= 0 {
		return nil
	}

	if !c.cfg.InputNumWithoutContext || ctx == nil {
		return nil
	}

	num--
	if num <= limit {
		return nil
	}

	x.AddIssue(&Issue{
		Filename:  x.Name,
		BeginLine: pos.Line,
		EndLine:   pos.Line,
		Level:     LevelError,
		Description: fmt.Sprintf("the num of input parameters without '%s' should not be more than %d, there are %d",
			ctx.Type.String, limit, num),
	})
	return nil
}

func (c funcParamsChecker) checkInNumValid(x *Context, pos token.Position, params []*funcParamInfo, ctx *funcParamInfo) Error {
	limit := c.cfg.InputNum
	num := len(params)

	if limit <= 0 {
		return nil
	}

	if c.cfg.InputNumWithoutContext && ctx != nil {
		return nil
	}

	if num <= limit {
		return nil
	}

	x.AddIssue(&Issue{
		Filename:  x.Name,
		BeginLine: pos.Line,
		EndLine:   pos.Line,
		Level:     LevelError,
		Description: fmt.Sprintf("the num of input parameters should not be more than %d, there are %d",
			limit, num),
	})
	return nil
}

func (c funcParamsChecker) checkOutNumValidWithoutError(x *Context, pos token.Position, params []*funcParamInfo, err *funcParamInfo) Error {
	limit := c.cfg.OutputNum
	num := len(params)

	if limit <= 0 {
		return nil
	}

	if !c.cfg.OutputNumWithoutError || err == nil {
		return nil
	}

	num--
	if num <= limit {
		return nil
	}

	x.AddIssue(&Issue{
		Filename:  x.Name,
		BeginLine: pos.Line,
		EndLine:   pos.Line,
		Level:     LevelError,
		Description: fmt.Sprintf("the num of output parameters without '%s' should not be more than %d, there are %d",
			err.Type.String, limit, num),
	})
	return nil
}

func (c funcParamsChecker) checkOutNumValid(x *Context, pos token.Position, params []*funcParamInfo, err *funcParamInfo) Error {
	limit := c.cfg.OutputNum
	num := len(params)

	if limit <= 0 {
		return nil
	}

	if c.cfg.OutputNumWithoutError && err != nil {
		return nil
	}

	if num <= limit {
		return nil
	}

	x.AddIssue(&Issue{
		Filename:  x.Name,
		BeginLine: pos.Line,
		EndLine:   pos.Line,
		Level:     LevelError,
		Description: fmt.Sprintf("the num of output parameters should not be more than %d, there are %d",
			limit, num),
	})
	return nil
}

func (c funcParamsChecker) checkContextErrorMatch(x *Context, pos token.Position, ctx, err *funcParamInfo) Error {
	// if ctx != nil {
	// 	log.Printf(" === ctx:%s, %v", ctx.Type.String, c.isStdContext(x, ctx))
	// }
	// if err != nil {
	// 	log.Printf(" === err:%s, %v", err.Type.String, c.isStdError(x, err))
	// }

	if ctx == nil || err == nil {
		return nil
	}
	if !c.cfg.ContextErrorMatch {
		return nil
	}

	isStdCtx := c.isStdContext(x, ctx)
	isStdErr := c.isStdError(x, err)
	if isStdCtx == isStdErr {
		return nil
	}

	desc := ""
	if isStdCtx {
		desc = fmt.Sprintf("context '%s' is standard context, but error '%s' is not standard error",
			ctx.Type.String, err.Type.String)
	} else {
		desc = fmt.Sprintf("context '%s' is not standard context, but error '%s' is standard error",
			ctx.Type.String, err.Type.String)
	}

	// log.Printf(" === ctx:%#v(%v) err:%#v(%v) desc:%s", ctx.Type, isStdCtx, err.Type, isStdErr, desc)

	x.AddIssue(&Issue{
		Filename:    x.Name,
		BeginLine:   pos.Line,
		EndLine:     pos.Line,
		Level:       LevelError,
		Description: desc,
	})

	return nil
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
		if imp.Path == "context" && (typ.Package == imp.Path || typ.Package == imp.Name) {
			return true
		}
	}

	return false
}

func (c funcParamsChecker) isStdError(x *Context, err *funcParamInfo) bool {
	typ := err.Type
	return !typ.IsPointer && typ.Package == "" && typ.Name == "error"
}

func (c funcParamsChecker) getOutputTypeName(x *Context, out *funcParamInfo) string {
	typeName := out.Type.Name
	typeName = strings.ToLower(typeName)

	if strings.HasPrefix(typeName, "int") {
		return "-number"
	}

	if strings.HasPrefix(typeName, "uint") {
		return "-number"
	}

	if strings.HasPrefix(typeName, "float") {
		return "-number"
	}

	if strings.HasPrefix(typeName, "byte") {
		return "-number"
	}

	return typeName
}

func (c funcParamsChecker) checkNamedOutput(x *Context, pos token.Position, outs []*funcParamInfo) Error {
	if !c.cfg.NamedOutput {
		return nil
	}

	if len(outs) <= 1 {
		return nil
	}

	if outs[0].Name != "" {
		return nil
	}

	m := make(map[string]struct{}, len(outs))
	for _, out := range outs {
		typeName := c.getOutputTypeName(x, out)
		if _, ok := m[typeName]; !ok {
			m[typeName] = struct{}{}
			continue
		}

		x.AddIssue(&Issue{
			Filename:    x.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("output parameter types are similar, please name them"),
		})
		break
	}

	return nil
}
