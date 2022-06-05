/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"reflect"
	"testing"
)

func TestCoverages_Add(t *testing.T) {
	type args struct {
		c Coverage
	}
	s := Coverages{
		{BeginLine: 100, EndLine: 110},
		{BeginLine: 120, EndLine: 130},
		{BeginLine: 140, EndLine: 140},
		{BeginLine: 160, EndLine: 170},
		{BeginLine: 180, EndLine: 180},
	}
	tests := []struct {
		name string
		s    Coverages
		args args
		want Coverages
	}{
		{
			name: "add before head",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 90, EndLine: 95}},
			want: Coverages{
				{BeginLine: 90, EndLine: 95},
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add after tail",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 200, EndLine: 210}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
				{BeginLine: 200, EndLine: 210},
			},
		},
		{
			name: "add in middle without overlap 1",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 112, EndLine: 117}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 112, EndLine: 117},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle without overlap 2",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 111, EndLine: 119}},
			want: Coverages{
				{BeginLine: 100, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with overlap without merge 1",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 123, EndLine: 127}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with overlap without merge 2",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 135, EndLine: 145}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 135, EndLine: 145},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge left 1",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 125, EndLine: 145}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 145},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge left 2",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 111, EndLine: 145}},
			want: Coverages{
				{BeginLine: 100, EndLine: 145},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge right 1",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 135, EndLine: 160}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 135, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge right 2",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 135, EndLine: 179}},
			want: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 135, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge left and right 1",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 111, EndLine: 179}},
			want: Coverages{
				{BeginLine: 100, EndLine: 180},
			},
		},
		{
			name: "add in middle with merge left and right 2",
			s:    append(Coverages(nil), s...),
			args: args{c: Coverage{BeginLine: 90, EndLine: 200}},
			want: Coverages{
				{BeginLine: 90, EndLine: 200},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Add(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Coverages.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoverages_GetCount(t *testing.T) {
	tests := []struct {
		name string
		s    Coverages
		want int
	}{
		{
			s: Coverages{
				{BeginLine: 100, EndLine: 110},
				{BeginLine: 120, EndLine: 130},
				{BeginLine: 140, EndLine: 140},
				{BeginLine: 160, EndLine: 170},
				{BeginLine: 180, EndLine: 180},
			},
			want: 35,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetCount(); got != tt.want {
				t.Errorf("Coverages.GetCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoverages_Index(t *testing.T) {
	type args struct {
		begin int
		end   int
	}
	s := Coverages{
		{BeginLine: 100, EndLine: 110},
		{BeginLine: 120, EndLine: 130},
		{BeginLine: 140, EndLine: 140},
		{BeginLine: 160, EndLine: 170},
		{BeginLine: 180, EndLine: 180},
	}
	tests := []struct {
		name string
		s    Coverages
		args args
		want int
	}{
		{
			s:    s,
			args: args{begin: 99, end: 99},
			want: -1,
		},
		{
			s:    s,
			args: args{begin: 99, end: 100},
			want: 0,
		},
		{
			s:    s,
			args: args{begin: 99, end: 120},
			want: 0,
		},
		{
			s:    s,
			args: args{begin: 115, end: 120},
			want: 1,
		},
		{
			s:    s,
			args: args{begin: 115, end: 119},
			want: -1,
		},
		{
			s:    s,
			args: args{begin: 200, end: 200},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Index(tt.args.begin, tt.args.end); got != tt.want {
				t.Errorf("Coverages.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}
