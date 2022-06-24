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

func Test_packageChecker_Check(t *testing.T) {
	filename := "test.go"

	type args struct {
		name string
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
				name: filename,
				data: `
package test_package
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   2,
					EndLine:     2,
					Level:       LevelError,
					Description: "do not use the underscore in package name 'test_package'",
				},
			},
		},
		{
			args: args{
				name: filename,
				data: `
package testPackage
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   2,
					EndLine:     2,
					Level:       LevelError,
					Description: "do not use capital letters in package name 'testPackage'",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := packageChecker{}
			f := filecore.NewTestFile(tt.args.name, strcore.StrToBytesUnsafe(tt.args.data))

			if got := c.Check(f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("packageChecker.Check() = %v, want %v", testIssuesString(got), testIssuesString(tt.want))
			}
		})
	}
}
