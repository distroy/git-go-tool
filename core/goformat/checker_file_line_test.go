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

func Test_fileLineChecker_Check(t *testing.T) {
	filename := "test.go"

	type args struct {
		fileLine int
		filename string
		data     string
	}
	tests := []struct {
		name string
		c    importChecker
		args args
		want []*Issue
	}{
		{
			args: args{
				fileLine: 10,
				data: `
package test
import (
	"reflect"
	"testing"
)
`,
			},
			want: []*Issue{},
		},
		{
			args: args{
				fileLine: 5,
				filename: "for_test.go",
				data: `
package test
import (
	"reflect"
	"testing"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)
`,
			},
			want: []*Issue{},
		},
		{
			args: args{
				fileLine: 5,
				filename: "",
				data: `
package test
import (
	"reflect"
	"testing"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					Level:       LevelError,
					Description: "file lines(8) is more than 5, should split the file",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.filename == "" {
				tt.args.filename = filename
			}
			c := FileLineChecker(tt.args.fileLine)
			f := filecore.NewTestFile(tt.args.filename, strcore.StrToBytesUnsafe(tt.args.data))
			x := NewContext(f)

			c.Check(x)
			if got := x.Issues(); !reflect.DeepEqual(got, tt.want) {
				testPrintCheckResult(t, got, tt.want)
			}
		})
	}
}
