/*
 * Copyright (C) distroy
 */

package git

import (
	"reflect"
	"testing"
)

func Test_parseFilenameFromNewFileLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "+++ b/core/git/git.go",
			args:    args{line: "+++ b/core/git/git.go"},
			want:    "core/git/git.go",
			wantErr: false,
		},
		{
			name:    "+++ /dev/null",
			args:    args{line: "+++ /dev/null"},
			want:    "/dev/null",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFilenameFromNewFileLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFilenameFromNewFileLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseFilenameFromNewFileLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePositionFromSummaryLine(t *testing.T) {
	type args struct {
		summary string
	}
	tests := []struct {
		name      string
		args      args
		wantBegin int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "@@ -0,0 +10,32 @@",
			args:      args{summary: "@@ -0,0 +10,32 @@"},
			wantBegin: 10,
			wantEnd:   41,
			wantErr:   false,
		},
		{
			name:      "@@ -52 +52 @@",
			args:      args{summary: "@@ -52 +52 @@"},
			wantBegin: 52,
			wantEnd:   52,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBegin, gotEnd, err := parsePositionFromSummaryLine(tt.args.summary)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePositionFromSummaryLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBegin != tt.wantBegin {
				t.Errorf("parsePositionFromSummaryLine() gotBegin = %v, want %v", gotBegin, tt.wantBegin)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("parsePositionFromSummaryLine() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func Test_removeLineNoFromDifferent(t *testing.T) {
	type args struct {
		diff    *Different
		lineNos []int
	}
	tests := []struct {
		name string
		args args
		want []*Different
	}{
		{
			name: "remove head and tail",
			args: args{
				diff:    &Different{BeginLine: 1, EndLine: 100},
				lineNos: []int{1, 2, 4, 5, 6, 50, 51, 99, 100},
			},
			want: []*Different{
				{BeginLine: 3, EndLine: 3},
				{BeginLine: 7, EndLine: 49},
				{BeginLine: 52, EndLine: 98},
			},
		},
		{
			name: "reserve head and tail",
			args: args{
				diff:    &Different{BeginLine: 1, EndLine: 100},
				lineNos: []int{2, 4, 5, 6, 50, 51, 99},
			},
			want: []*Different{
				{BeginLine: 1, EndLine: 1},
				{BeginLine: 3, EndLine: 3},
				{BeginLine: 7, EndLine: 49},
				{BeginLine: 52, EndLine: 98},
				{BeginLine: 100, EndLine: 100},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeLineNoFromDifferent(tt.args.diff, tt.args.lineNos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeLineNoFromDifferent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNewLinesFromStatmentLines(t *testing.T) {
	type args struct {
		filename string
		lines    []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Different
		wantErr bool
	}{
		{
			name: "",
			args: args{
				lines: []string{
					`@@ -107,0 +113,12 @@ func readFileLines(r iocore.LineReader) ([]string, error) {`,
					`+`,
					`+func getCommandString(c *exec.Cmd) string {`,
					`+       // report the exact executable path (plus args)`,
					`+       b := &strings.Builder{}`,
					`+       b.WriteString(c.Path)`,
					`+`,
					`+       for _, a := range c.Args[1:] {`,
					`+               b.WriteByte(' ')`,
					`+               b.WriteString(a)`,
					`+       }`,
					`+       return b.String()`,
					`+}`,
				},
			},
			want: []*Different{
				{BeginLine: 114, EndLine: 117},
				{BeginLine: 119, EndLine: 124},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNewLinesFromStatmentLines(tt.args.filename, tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNewLinesFromStatmentLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNewLinesFromStatmentLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBlankLineNosFromStatmentLines(t *testing.T) {
	type args struct {
		lines []string
		diff  *Different
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "",
			args: args{
				diff: &Different{BeginLine: 100},
				lines: []string{
					`@@ -0,0 +100,7 @@`,
					`-aaa`,
					`-bbb`,
					`+`,
					`+111`,
					`+222`,
					`+`,
					`+`,
					`+555`,
					`+`,
				},
			},
			want: []int{100, 103, 104, 106},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBlankLineNosFromStatmentLines(tt.args.lines, tt.args.diff); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBlankLineNosFromStatmentLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNewLinesFromFileLines(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Different
		wantErr bool
	}{
		{
			name: "",
			args: args{
				lines: []string{
					`diff --git a/.gitignore b/.gitignore`,
					`index 3d37341..2fb5f78 100644`,
					`--- a/.gitignore`,
					`+++ b/.gitignore`,
					`@@ -18,0 +19 @@ unit_test.db`,
					`+/cmd/git-diff-go-cognitive/git-diff-go-cognitive`,
				},
			},
			want: []*Different{
				{Filename: ".gitignore", BeginLine: 19, EndLine: 19},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNewLinesFromFileLines(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNewLinesFromFileLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNewLinesFromFileLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
