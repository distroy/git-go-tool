/*
 * Copyright (C) distroy
 */

package git

import (
	"reflect"
	"testing"
)

func Test_removeLineNoFromDifferent(t *testing.T) {
	type args struct {
		diff    Different
		lineNos []int
	}
	tests := []struct {
		name string
		args args
		want []Different
	}{
		{
			name: "remove head and tail",
			args: args{
				diff:    Different{BeginLine: 1, EndLine: 100},
				lineNos: []int{1, 2, 4, 5, 6, 50, 51, 99, 100},
			},
			want: []Different{
				{BeginLine: 3, EndLine: 3},
				{BeginLine: 7, EndLine: 49},
				{BeginLine: 52, EndLine: 98},
			},
		},
		{
			name: "reserve head and tail",
			args: args{
				diff:    Different{BeginLine: 1, EndLine: 100},
				lineNos: []int{2, 4, 5, 6, 50, 51, 99},
			},
			want: []Different{
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
