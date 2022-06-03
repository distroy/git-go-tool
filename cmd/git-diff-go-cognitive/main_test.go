/*
 * Copyright (C) distroy
 */

package main

import (
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/git"
)

func Test_differents_toMap(t *testing.T) {
	tests := []struct {
		name  string
		array differents
		want  map[string][]git.Different
	}{
		{
			name: "",
			array: []git.Different{
				{Filename: "a", BeginLine: 101, EndLine: 102},
				{Filename: "c", BeginLine: 104, EndLine: 110},
				{Filename: "b", BeginLine: 110, EndLine: 114},
				{Filename: "a", BeginLine: 104, EndLine: 110},
				{Filename: "b", BeginLine: 101, EndLine: 102},
			},
			want: map[string][]git.Different{
				"a": {
					{Filename: "a", BeginLine: 101, EndLine: 102},
					{Filename: "a", BeginLine: 104, EndLine: 110},
				},
				"b": {
					{Filename: "b", BeginLine: 101, EndLine: 102},
					{Filename: "b", BeginLine: 110, EndLine: 114},
				},
				"c": {
					{Filename: "c", BeginLine: 104, EndLine: 110},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.array.toMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("differents.toMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
