/*
 * Copyright (C) distroy
 */

package git

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func Test_parseNewLinesFromCommand(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
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
				exec.Command("echo", strings.Join([]string{
					`diff --git a/.gitignore b/.gitignore`,
					`index 3d37341..2fb5f78 100644`,
					`--- a/.gitignore`,
					`+++ b/.gitignore`,
					`@@ -18,0 +19 @@ unit_test.db`,
					`+/cmd/git-diff-go-cognitive/git-diff-go-cognitive`,
					`diff --git a/cmd/git-diff-go-cognitive/main.go b/cmd/git-diff-go-cognitive/main.go`,
					`index 1d86b6c..96742a4 100644`,
					`--- a/cmd/git-diff-go-cognitive/main.go`,
					`+++ b/cmd/git-diff-go-cognitive/main.go`,
					`@@ -8,0 +9 @@ import (`,
					`+       "log"`,
					`@@ -16,0 +19 @@ type Flags struct {`,
					`+       Branch string`,
					`@@ -20 +23 @@ type Flags struct {`,
					`-func paresFlags() *Flags {`,
					`+`,
					`@@ -35,0 +40,4 @@ func paresFlags() *Flags {`,
					`+       if f.Branch == "" {`,
					`+               f.Branch = git.GetBranch()`,
					`+       }`,
					`+`,
					``,
				}, "\n")),
			},
			want: []*Different{
				{Filename: ".gitignore", BeginLine: 19, EndLine: 19},
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 9, EndLine: 9},
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 19, EndLine: 19},
				{Filename: "cmd/git-diff-go-cognitive/main.go", BeginLine: 40, EndLine: 42},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNewLinesFromCommand(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNewLinesFromCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNewLinesFromCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
