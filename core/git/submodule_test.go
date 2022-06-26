/*
 * Copyright (C) distroy
 */

package git

import (
	"reflect"
	"testing"
)

func Test_parseSubModules(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []*SubModule
	}{
		{
			args: args{
				s: `
-283a207a9ea00874fb5261df490d017305b6e648 git-go-tool
 358cd80c59873f54a1790e3628cc392cd3e72019 zsh/on-my-zsh (v3.0-2-g283a207)
`,
			},
			want: []*SubModule{
				{CommitId: "283a207a9ea00874fb5261df490d017305b6e648", Path: `git-go-tool`},
				{CommitId: "358cd80c59873f54a1790e3628cc392cd3e72019", Path: `zsh/on-my-zsh`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSubModules(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSubModules() = %v, want %v", got, tt.want)
			}
		})
	}
}
