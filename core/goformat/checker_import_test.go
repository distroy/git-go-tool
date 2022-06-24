/*
 * Copyright (C) distroy
 */

package goformat

import (
	"reflect"
	"strings"
	"testing"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)

func testIssuesString(issues []*Issue) string {
	buf := &strings.Builder{}
	NewIssueWriter(buf).WriteIssues(issues)
	return buf.String()
}

func testPrintCheckResult(t testing.TB, got, want []*Issue) {
	t.Errorf("Check() = \n%v\n    want:\n%v", testIssuesString(got), testIssuesString(want))
}

func Test_importChecker_Check(t *testing.T) {
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
					BeginLine:   4,
					EndLine:     7,
					Level:       LevelError,
					Description: "imports should be grouped and ordered by standards and others",
				},
			},
		},
		{
			args: args{
				data: `
package test
import (
	"reflect"
	"github.com/distroy/git-go-tool/core/filecore"
	"testing"
	. "io"
	"github.com/distroy/git-go-tool/core/strcore"
)
`,
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   7,
					EndLine:     7,
					Level:       LevelError,
					Description: "do not use the dot import",
				},
				{
					Filename:    filename,
					BeginLine:   4,
					EndLine:     8,
					Level:       LevelError,
					Description: "imports should be grouped and ordered by standards and others",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ImportChecker(true)
			f := filecore.NewTestFile(filename, strcore.StrToBytesUnsafe(tt.args.data))
			x := NewContext(f)

			c.Check(x)
			if got := x.Issues(); !reflect.DeepEqual(got, tt.want) {
				testPrintCheckResult(t, got, tt.want)
			}
		})
	}
}
