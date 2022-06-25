/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewFlagSet(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFlagSet(); tt.want != (got != nil) {
				t.Errorf("NewFlagSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSet_printUsageHeader(t *testing.T) {
	type fields struct {
		name string
		args *Flag
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: `name == "" && no args`,
			fields: fields{
				name: "",
				args: nil,
			},
			wantW: "Usage of <command>:\nFlags:\n",
		},
		{
			name: `name == "abc" && args.meta == ""`,
			fields: fields{
				name: "abc",
				args: &Flag{Meta: ""},
			},
			wantW: "Usage: abc [<flags>] [<arg>...]\n\nFlags:\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FlagSet{
				name: tt.fields.name,
				args: tt.fields.args,
			}
			w := &bytes.Buffer{}
			s.printUsageHeader(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FlagSet.printUsageHeader() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func testClearFlagsUnexportedField(s []*Flag) {
	for _, v := range s {
		v.tags = nil
		v.val = reflect.Value{}
		v.Value = nil
	}
}

func TestFlagSet_Model(t *testing.T) {
	type Flags struct {
		Top      int      `flag:"name:top; meta:N; usage:show the top <N>"`
		Avg      bool     `flag:"usage:show the average complexity"`
		DebugLog bool     `flag:"usage:print debug log; bool"`
		Rate     float64  `flag:"default:0.65; usage:"`
		Branch   string   `flag:"meta:branch; usage:git branch name"`
		Pathes   []string `flag:"args; meta:path; default:."`
	}

	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want []*Flag
	}{
		{
			args: args{v: &Flags{}},
			want: []*Flag{
				{
					lvl:     0,
					Name:    "top",
					Meta:    "N",
					Default: "0",
					Usage:   "show the top <N>",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "avg",
					Meta:    "",
					Default: "false",
					Usage:   "show the average complexity",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "debug-log",
					Meta:    "",
					Default: "false",
					Usage:   "print debug log",
					IsArgs:  false,
					Bool:    true,
				},
				{
					lvl:     0,
					Name:    "rate",
					Meta:    "",
					Default: "0.65",
					Usage:   "",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "pathes",
					Meta:    "path",
					Default: ".",
					Usage:   "",
					IsArgs:  false,
					Bool:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFlagSet()
			s.Model(tt.args.v)

			if got := s.flagSlice; reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagSet.Model() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
