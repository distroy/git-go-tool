/*
 * Copyright (C) distroy
 */

package strcore

import (
	"reflect"
	"testing"
)

func TestBytesToStr(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{b: nil},
			want: "",
		},
		{
			name: "abc",
			args: args{b: []byte("abc")},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToStr(tt.args.b); got != tt.want {
				t.Errorf("BytesToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestStrToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "",
			args: args{s: ""},
			want: []byte(""),
		},
		{
			name: "aaa",
			args: args{s: "aaa"},
			want: []byte("aaa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToStrUnsafe(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{b: nil},
			want: "",
		},
		{
			name: "abc",
			args: args{b: []byte("abc")},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToStrUnsafe(tt.args.b); got != tt.want {
				t.Errorf("BytesToStrUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToBytesUnsafe(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil",
			args: args{s: ""},
			want: []byte{},
		},
		{
			name: "aaa",
			args: args{s: "aaa"},
			want: []byte("aaa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrToBytesUnsafe(tt.args.s)
			if got == nil {
				got = []byte{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToBytesUnsafe() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
