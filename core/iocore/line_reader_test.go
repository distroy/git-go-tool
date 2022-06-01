/*
 * Copyright (C) distroy
 */

package iocore

import (
	"io"
	"strings"
	"testing"
)

type testLineReaderWant struct {
	want    string
	wantErr bool
}

func testLineReader(t *testing.T, f func() (string, error), tests []testLineReaderWant) {
	for _, tt := range tests {
		got, err := f()
		if (err != nil) != tt.wantErr {
			t.Errorf("%s() error = %v, wantErr %v", t.Name(), err, tt.wantErr)
			return
		}
		if err != nil && got != tt.want {
			t.Errorf("%s() = %v want:%s", t.Name(), got, tt.want)
		}
	}
}

func TestLineReader_Peek(t *testing.T) {
	text := "1111\r\n2222\n3333\r\n"
	r := NewLineReader(strings.NewReader(text))

	tests := []struct{ want string }{
		{want: "1111"},
		{want: "2222"},
		{want: "3333"},
	}

	for _, tt := range tests {
		got, err := r.PeekString()
		if err != nil {
			t.Errorf("LineReader.PeekString() error = %v", err)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.PeekString() = %v, want:%s", got, tt.want)
		}
		r.ReadString()
	}

	if got, err := r.PeekString(); err != io.EOF {
		t.Errorf("LineReader.ReadString() = %v, error = %s", got, err)
	}
}

func TestLineReader_OverSize(t *testing.T) {
	text := "1234"
	r := NewLineReader(strings.NewReader(text))
	r.maxSize = 4

	if got, err := r.PeekString(); err != ErrOverMaxSize {
		t.Errorf("LineReader.PeekString() = %v, error = %s", got, err)
	}

	if got, err := r.ReadString(); err != ErrOverMaxSize {
		t.Errorf("LineReader.ReadString() = %v, error = %s", got, err)
	}
}

func TestLineReader_Read(t *testing.T) {
	text := "1234\nabcd"
	r := NewLineReader(strings.NewReader(text))

	tests := []struct{ want string }{
		{want: "1234"},
		{want: "abcd"},
	}

	for _, tt := range tests {
		got, err := r.ReadString()
		if err != nil {
			t.Errorf("LineReader.ReadString() error:%v, want:%s", err, tt.want)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.ReadString() = %v, want:%s", got, tt.want)
		}
	}

	if got, err := r.PeekString(); err != io.EOF {
		t.Errorf("LineReader.ReadString() = %v, error = %s", got, err)
	}
	if got, err := r.PeekString(); err != io.EOF {
		t.Errorf("LineReader.ReadString() = %v, error = %s", got, err)
	}
}
