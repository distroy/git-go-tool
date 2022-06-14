/*
 * Copyright (C) distroy
 */

package main

import (
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/gocoverage"
)

func Test_getTopNonCoverageFiles(t *testing.T) {
	type args struct {
		files gocoverage.Files
		top   int
	}
	tests := []struct {
		name string
		args args
		want Files
	}{
		{
			name: "top 2 without all coverage file",
			args: args{
				top: 2,
				files: gocoverage.Files{
					"a": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "a", BeginLine: 1, EndLine: 10},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "a", BeginLine: 101, EndLine: 120},
						},
					},
					"b": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 1, EndLine: 40},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 101, EndLine: 119},
						},
					},
					"c": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "c", BeginLine: 1, EndLine: 9},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "c", BeginLine: 101, EndLine: 120},
						},
					},
				},
			},
			want: Files{
				{
					Count: gocoverage.Count{Coverages: 9, NonCoverages: 20},
					File: &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "c", BeginLine: 1, EndLine: 9},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "c", BeginLine: 101, EndLine: 120},
						},
					},
				},
				{
					Count: gocoverage.Count{Coverages: 10, NonCoverages: 20},
					File: &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "a", BeginLine: 1, EndLine: 10},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "a", BeginLine: 101, EndLine: 120},
						},
					},
				},
			},
		},
		{
			name: "top 2 with all coverage file",
			args: args{
				top: 2,
				files: gocoverage.Files{
					"a": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "a", BeginLine: 1, EndLine: 10},
						},
						NonCoverages: gocoverage.Coverages{},
					},
					"b": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 1, EndLine: 40},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 101, EndLine: 119},
						},
					},
					"c": &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "c", BeginLine: 1, EndLine: 9},
						},
						NonCoverages: gocoverage.Coverages{},
					},
				},
			},
			want: Files{
				{
					Count: gocoverage.Count{Coverages: 40, NonCoverages: 19},
					File: &gocoverage.FileCoverages{
						Coverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 1, EndLine: 40},
						},
						NonCoverages: gocoverage.Coverages{
							{Filename: "b", BeginLine: 101, EndLine: 119},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopNonCoverageFiles(tt.args.files, tt.args.top); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTopNonCoverageFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_GetNonCoverageLineNos(t *testing.T) {
	type fields struct {
		Count gocoverage.Count
		File  *gocoverage.FileCoverages
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				// Count: gocoverage.Count{Coverages: 40, NonCoverages: 19},
				File: &gocoverage.FileCoverages{
					// Coverages: gocoverage.Coverages{
					// 	{Filename: "b", BeginLine: 1, EndLine: 40},
					// },
					NonCoverages: gocoverage.Coverages{
						{Filename: "b", BeginLine: 10, EndLine: 10},
						{Filename: "b", BeginLine: 20, EndLine: 25},
					},
				},
			},
			want: "10, [20,25]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File{
				Count: tt.fields.Count,
				File:  tt.fields.File,
			}
			if got := f.GetNonCoverageLineNos(); got != tt.want {
				t.Errorf("File.GetNonCoverageLineNos() = %v, want %v", got, tt.want)
			}
		})
	}
}
