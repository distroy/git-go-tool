/*
 * Copyright (C) distroy
 */

package goformat

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)

func testIssuesString(issues []*Issue) string {
	data, _ := json.Marshal(issues)
	return strcore.BytesToStrUnsafe(data)
}

func Test_importChecker_Check(t *testing.T) {
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
				name: filename,
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
				name: filename,
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
			c := importChecker{}
			f := filecore.NewTestFile(tt.args.name, strcore.StrToBytesUnsafe(tt.args.data))

			if got := c.Check(f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importChecker.Check() = %v, want %v", testIssuesString(got), testIssuesString(tt.want))
			}
		})
	}
}
