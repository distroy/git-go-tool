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
				data: `
package test
func TestFunc(a, b, c, d int) (o, p, q, r int) {
	fn1 := func(a, b, c, d, e int) (o, p, q, r, s int) {
		return a ^ c, b & d, c - e, d % a, e / b
	}
	fn2 := func(a, b, c int) (o, p, q int) {
		return a ^ c, b & d, c - a
	}
	func() {
		go func (a, b, c, d int) (o, p, q, r int) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FuncParamsChecker(&FuncParamsConfig{
				InputNum:  3,
				OutputNum: 3,
			})
			f := filecore.NewTestFile(filename, strcore.StrToBytesUnsafe(tt.args.data))
			x := NewContext(f)

			c.Check(x)
			if got := x.Issues(); !reflect.DeepEqual(got, tt.want) {
				testPrintCheckResult(t, got, tt.want)
			}
		})
	}
}
