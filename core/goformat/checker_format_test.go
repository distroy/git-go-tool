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

func Test_formatChecker_Check(t *testing.T) {
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
package test
func TestFunc() {}
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					Level:       LevelError,
					Description: "source should be formated",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := formatChecker{}
			f := filecore.NewTestFile(tt.args.name, strcore.StrToBytesUnsafe(tt.args.data))

			if got := c.Check(f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatChecker.Check() = %v, want %v", testIssuesString(got), testIssuesString(tt.want))
			}
		})
	}
}
