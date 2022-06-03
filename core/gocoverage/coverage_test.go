/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"reflect"
	"testing"
)

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
				prefix: "github.com/distroy/git-go-tool/",
				line:   "mode: set",
			},
			want:  nil,
			want1: true,
		},
		{
			name: "coverage",
			args: args{
				prefix: "github.com/distroy/git-go-tool/",
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
