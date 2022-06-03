/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseReader(t *testing.T) {
	type args struct {
		prefix string
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Coverage
		wantErr bool
	}{
		{
			args: args{
				prefix: "github.com/distroy/git-go-tool",
				reader: bytes.NewBufferString(strings.Join([]string{
					`mode: set`,
					`github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive/main.go:29.26,46.20 8 0`,
					`github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive/main.go:49.2,49.16 1 1`,
					`github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive/main.go:53.2,53.10 1 0`,
					`github.com/distroy/git-go-tool/core/iocore/line_reader.go:84.2,84.33 1 1`,
					`github.com/distroy/git-go-tool/core/iocore/line_reader.go:89.2,89.19 1 0`,
					`github.com/distroy/git-go-tool/core/iocore/line_reader.go:80.16,82.3 1 1`,
				}, "\n")),
			},
			want: []Coverage{
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 29, EndLine: 46, Count: 0},
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 49, EndLine: 49, Count: 1},
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 53, EndLine: 53, Count: 0},
				{Filename: "core/iocore/line_reader.go", BeginLine: 84, EndLine: 84, Count: 1},
				{Filename: "core/iocore/line_reader.go", BeginLine: 89, EndLine: 89, Count: 0},
				{Filename: "core/iocore/line_reader.go", BeginLine: 80, EndLine: 82, Count: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReader(tt.args.prefix, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLine(t *testing.T) {
	type args struct {
		prefix string
		line   string
	}
	tests := []struct {
		name  string
		args  args
		want  *Coverage
		want1 bool
	}{
		{
			name: "mode line",
			args: args{
				prefix: "github.com/distroy/git-go-tool",
				line:   "mode: set",
			},
			want:  nil,
			want1: true,
		},
		{
			name: "coverage",
			args: args{
				prefix: "github.com/distroy/git-go-tool",
				line:   "github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive/main.go:83.22,86.3 2 1",
			},
			want:  &Coverage{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 83, EndLine: 86, Count: 1},
			want1: true,
		},
		{
			name: "non coverage",
			args: args{
				prefix: "github.com/distroy/git-go-tool/",
				line:   "github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive/main.go:49.16,51.3 1 0",
			},
			want:  &Coverage{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 49, EndLine: 51, Count: 0},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseLine(tt.args.prefix, tt.args.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
