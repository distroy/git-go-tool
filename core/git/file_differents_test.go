/*
 * Copyright (C) distroy
 */

package git

import (
	"reflect"
	"testing"
)

func TestDifferents_IsIn(t *testing.T) {
	type args struct {
		begin int
		end   int
	}
	s := Differents{
		{BeginLine: 100, EndLine: 110},
		{BeginLine: 120, EndLine: 130},
		{BeginLine: 140, EndLine: 140},
	}
	tests := []struct {
		name string
		s    Differents
		args args
		want bool
	}{
		{
			name: "< head",
			s:    s,
			args: args{begin: 99, end: 99},
			want: false,
		},
		{
			name: "in the head 1",
			s:    s,
			args: args{begin: 99, end: 100},
			want: true,
		},
		{
			name: "in the head 2",
			s:    s,
			args: args{begin: 108, end: 112},
			want: true,
		},
		{
			name: "> head && < middle",
			s:    s,
			args: args{begin: 113, end: 117},
			want: false,
		},
		{
			name: "in the middle 1",
			s:    s,
			args: args{begin: 118, end: 132},
			want: true,
		},
		{
			name: "in the middle 2",
			s:    s,
			args: args{begin: 124, end: 126},
			want: true,
		},
		{
			name: "in the tail 1",
			s:    s,
			args: args{begin: 140, end: 140},
			want: true,
		},
		{
			name: "> tail",
			s:    s,
			args: args{begin: 141, end: 141},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.IsIn(tt.args.begin, tt.args.end)
			if got != tt.want {
				t.Errorf("Differents.IsIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFileDifferents(t *testing.T) {
	type args struct {
		slice []Different
	}
	tests := []struct {
		name string
		args args
		want Files
	}{
		{
			name: "",
			args: args{
				slice: Differents{
					{Filename: "a", BeginLine: 101, EndLine: 102},
					{Filename: "c", BeginLine: 104, EndLine: 110},
					{Filename: "b", BeginLine: 110, EndLine: 114},
					{Filename: "a", BeginLine: 104, EndLine: 110},
					{Filename: "b", BeginLine: 101, EndLine: 102},
				},
			},
			want: map[string]Differents{
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
			if got := NewFileDifferents(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileDifferents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileDifferents_IsIn(t *testing.T) {
	type args struct {
		file  string
		begin int
		end   int
	}
	tests := []struct {
		name string
		m    Files
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsIn(tt.args.file, tt.args.begin, tt.args.end); got != tt.want {
				t.Errorf("FileDifferents.IsIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
