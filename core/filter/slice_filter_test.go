/*
 * Copyright (C) distroy
 */

package filter

import (
	"reflect"
	"testing"
)

func TestFilterSlice(t *testing.T) {
	type args struct {
		slice  interface{}
		filter interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{
				slice:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				filter: func(n int) bool { return n%2 == 0 },
			},
			want: []int{0, 2, 8, 4, 6},
		},
		{
			args: args{
				slice:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				filter: func(n int) bool { return true },
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := FilterSlice(tt.args.slice, tt.args.filter)
			got := reflect.ValueOf(tt.args.slice).Slice(0, n).Interface()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
