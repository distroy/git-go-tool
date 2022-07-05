/*
 * Copyright (C) distroy
 */

package goformat

import (
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)

func Test_funcParamsChecker_Check(t *testing.T) {
	filename := "test.go"

	type args struct {
		cfg  *FuncParamsConfig
		data string
	}
	tests := []struct {
		name string
		c    importChecker
		args args
		want []*Issue
	}{
		{
			args: args{
				cfg: &FuncParamsConfig{
					InputNum:  3,
					OutputNum: 3,
				},
				data: `
package test
func TestFunc(a *Context, b, c, d int) (o, p, q int, r Error) {
	fn1 := func(a, b, c, d, e int) (o, p, q, r, s int) {
		return a ^ c, b & d, c - e, d % a, e / b
	}
	fn2 := func(a, b, c int) (o, p, q int) {
		return a ^ c, b & d, c - a
	}
	func() {
		go func (a, b, c, d int) (int, int, uint, int) {
			a, b, c, d = fn(a, b, c, d)
			return a + b, b - c, c * d, d / a
		}(a, b, c , d)
	}()
}
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   3,
					EndLine:     3,
					Level:       LevelError,
					Description: "the num of input parameters should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   3,
					EndLine:     3,
					Level:       LevelError,
					Description: "the num of output parameters should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   4,
					EndLine:     4,
					Level:       LevelError,
					Description: "the num of input parameters should not be more than 3, there are 5",
				},
				{
					Filename:    filename,
					BeginLine:   4,
					EndLine:     4,
					Level:       LevelError,
					Description: "the num of output parameters should not be more than 3, there are 5",
				},
				{
					Filename:    filename,
					BeginLine:   11,
					EndLine:     11,
					Level:       LevelError,
					Description: "the num of input parameters should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   11,
					EndLine:     11,
					Level:       LevelError,
					Description: "the num of output parameters should not be more than 3, there are 4",
				},
			},
		},
		{
			args: args{
				cfg: &FuncParamsConfig{
					InputNum:               3,
					OutputNum:              3,
					InputNumWithoutContext: true,
					OutputNumWithoutError:  true,
					ContextFirst:           true,
					ErrorLast:              true,
					ContextErrorMatch:      true,
					NamedOutput:            true,
				},
				data: `
package test

import (
	"context"
	ctx "context"
)

func TestFunc(ctx *Context, a, b, c int) (o, p, q int, err error) {
	fn1 := func(a, b, c int) (o, p, q int) {
		return a ^ c, b & d, c - e, d % a, e / b
	}
	fn2 := func(_ *Context, a, b, c, d int) (o, p, q , r int, err *error) {
		return a ^ c, b & d, c - a, d + a, nil
	}
	fn3 := func(ctx context.Context) *error { return nil }
	fn4 := func(ctx ctx.Context) Error { return nil }
	func() {
		go func (a, b, c, d int) (int, int, uint, int) {
			a, b, c, d = fn1(a, b, c, d)
			return a + b, b - c, c * d, d / a
		}(a, b, c , d)
	}()
}
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   9,
					EndLine:     9,
					Level:       LevelError,
					Description: "context '*Context' is not standard context, but error 'error' is standard error",
				},
				{
					Filename:    filename,
					BeginLine:   13,
					EndLine:     13,
					Level:       LevelError,
					Description: "the num of input parameters without '*Context' should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   13,
					EndLine:     13,
					Level:       LevelError,
					Description: "the num of output parameters without '*error' should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   16,
					EndLine:     16,
					Level:       LevelError,
					Description: "context 'context.Context' is standard context, but error '*error' is not standard error",
				},
				{
					Filename:    filename,
					BeginLine:   17,
					EndLine:     17,
					Level:       LevelError,
					Description: "context 'ctx.Context' is standard context, but error 'Error' is not standard error",
				},
				{
					Filename:    filename,
					BeginLine:   19,
					EndLine:     19,
					Level:       LevelError,
					Description: "output parameter types are similar, please name them",
				},
				{
					Filename:    filename,
					BeginLine:   19,
					EndLine:     19,
					Level:       LevelError,
					Description: "the num of input parameters should not be more than 3, there are 4",
				},
				{
					Filename:    filename,
					BeginLine:   19,
					EndLine:     19,
					Level:       LevelError,
					Description: "the num of output parameters should not be more than 3, there are 4",
				},
			},
		},
		{
			args: args{
				cfg: &FuncParamsConfig{
					InputNum:               3,
					OutputNum:              3,
					InputNumWithoutContext: true,
					OutputNumWithoutError:  true,
					ContextFirst:           true,
					ErrorLast:              true,
					ContextErrorMatch:      true,
					NamedOutput:            true,
				},
				data: `
package test

func TestFunc(f func (x *Context) error) {
}
`,
			},
			want: []*Issue{
				// {
				// 	Filename:    filename,
				// 	BeginLine:   4,
				// 	EndLine:     4,
				// 	Level:       LevelError,
				// 	Description: "context '*Context' is not standard context, but error 'error' is standard error",
				// },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FuncParamsChecker(tt.args.cfg)
			f := filecore.NewTestFile(filename, strcore.StrToBytesUnsafe(tt.args.data))
			x := NewContext(f)

			c.Check(x)
			if got := x.Issues(); !reflect.DeepEqual(got, tt.want) {
				testPrintCheckResult(t, got, tt.want)
			}
		})
	}
}
