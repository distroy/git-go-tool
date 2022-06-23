/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
)

type FuncParamsConfig struct {
	InputNum  int
	OutputNum int
}

func FuncParamsChecker(cfg *FuncParamsConfig) Checker {
	return funcParamsChecker{cfg: cfg}
}

type funcParamInfo struct {
	Name string
	Type *typeInfo
}

type funcParamsChecker struct {
	cfg *FuncParamsConfig
}

func (c funcParamsChecker) Check(f *filecore.File) []*Issue {
	res := make([]*Issue, 0, 8)

	file := f.MustParse()

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		res = c.checkFuncDecl(res, f, fn)

		res = c.walkFunc(res, f, fn)
	}

	return res
}

func (c funcParamsChecker) walkFunc(res []*Issue, f *filecore.File, fn *ast.FuncDecl) []*Issue {
	return res
}

func (c funcParamsChecker) convertParams(f *filecore.File, params *ast.FieldList) []*funcParamInfo {
	if params == nil {
		return nil
	}

	n := params.NumFields()
	res := make([]*funcParamInfo, 0, n)

	for _, param := range params.List {
		typ := getTypeInfo(f, param.Type)
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

func (c funcParamsChecker) checkFuncDecl(res []*Issue, f *filecore.File, fn *ast.FuncDecl) []*Issue {
	inLimit := c.cfg.InputNum
	outLimit := c.cfg.OutputNum

	ins := c.convertParams(f, fn.Type.Params)
	outs := c.convertParams(f, fn.Type.Results)

	inNum := len(ins)
	outNum := len(outs)
	if inNum == 0 && outNum == 0 {
		return res
	}

	pos := f.Position(fn.Pos())
	// log.Printf(" === file:%s:%d, func:%s", f.Name, pos.Line, fn.Name.Name)

	ctxIdx, ctx := c.indexParamByTypeName(ins, "context")
	if ctxIdx > 0 {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the input parameter '%s' should be the first", ctx.Type.String),
		})
	}

	if inLimit > 0 && ctx != nil && inNum > inLimit+1 {
		res = append(res, &Issue{
			Filename:  f.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters without '%s' should be less than %d, there are %d",
				ctx.Type.String, inLimit, inNum-1),
		})

	} else if inLimit > 0 && ctx == nil && inNum > inLimit {
		res = append(res, &Issue{
			Filename:  f.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of input parameters should be less than %d, there are %d",
				inLimit, inNum),
		})
	}

	errIdx, err := c.indexParamByTypeName(outs, "error")
	if errIdx >= 0 && errIdx != outNum-1 {
		res = append(res, &Issue{
			Filename:    f.Name,
			BeginLine:   pos.Line,
			EndLine:     pos.Line,
			Level:       LevelError,
			Description: fmt.Sprintf("the output parameter '%s' should be the last", err.Type.String),
		})
	}

	if outLimit > 0 && err != nil && outNum > outLimit+1 {
		res = append(res, &Issue{
			Filename:  f.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of output parameters without '%s' should be less than %d, there are %d",
				err.Type.String, outLimit, outNum-1),
		})

	} else if outLimit > 0 && err == nil && outNum > outLimit {
		res = append(res, &Issue{
			Filename:  f.Name,
			BeginLine: pos.Line,
			EndLine:   pos.Line,
			Level:     LevelError,
			Description: fmt.Sprintf("the num of output parameters should be less than %d, there are %d",
				outLimit, outNum),
		})
	}

	return res
}

func (c funcParamsChecker) indexParamByTypeName(params []*funcParamInfo, typeName string) (int, *funcParamInfo) {
	for i, v := range params {
		if strings.EqualFold(v.Type.Name, typeName) {
			return i, v
		}
	}
	return -1, nil
}
