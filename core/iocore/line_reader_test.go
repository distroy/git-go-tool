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

func TestLineReader_PeekLineString(t *testing.T) {
	text := "1111\r\n\n2222\n3333\r\n"
	r := NewLineReader(strings.NewReader(text))

	tests := []struct{ want string }{
		{want: "1111"},
		{want: ""},
		{want: "2222"},
		{want: "3333"},
	}

	for _, tt := range tests {
		got, err := r.PeekLineString()
		if err != nil {
			t.Errorf("LineReader.PeekLineString() error = %v", err)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.PeekLineString() = %v, want:%s", got, tt.want)
		}
		got, err = r.ReadLineString()
		if err != nil {
			t.Errorf("LineReader.ReadLineString() error = %v", err)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.ReadLineString() = %v, want:%s", got, tt.want)
		}
	}

	if got, err := r.PeekLineString(); err != io.EOF {
		t.Errorf("LineReader.ReadLineString() = %v, error = %s", got, err)
	}
}

func TestLineReader_OverSize(t *testing.T) {
	text := "a\nb\nc\r\n1234\n"
	r := NewLineReader(strings.NewReader(text), LineReaderBufferSize(4))

	tests := []struct{ want string }{
		{want: "a"},
		{want: "b"},
		{want: "c"},
	}
	for _, tt := range tests {
		got, err := r.PeekLineString()
		if err != nil {
			t.Errorf("LineReader.PeekLineString() error = %v", err)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.PeekLineString() = %v, want:%s", got, tt.want)
		}
		got, err = r.ReadLineString()
		if err != nil {
			t.Errorf("LineReader.ReadLineString() error = %v", err)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.ReadLineString() = %v, want:%s", got, tt.want)
		}
	}

	if got, err := r.PeekLineString(); err != ErrOverMaxSize {
		t.Errorf("LineReader.PeekLineString() = %v, error = %s", got, err)
	}

	if got, err := r.ReadLineString(); err != ErrOverMaxSize {
		t.Errorf("LineReader.ReadLineString() = %v, error = %s", got, err)
	}
}

func TestLineReader_ReadLineString(t *testing.T) {
	text := "1234\nabcd"
	r := NewLineReader(strings.NewReader(text))

	tests := []struct{ want string }{
		{want: "1234"},
		{want: "abcd"},
	}

	for _, tt := range tests {
		got, err := r.ReadLineString()
		if err != nil {
			t.Errorf("LineReader.ReadLineString() error:%v, want:%s", err, tt.want)
			return
		}
		if got != tt.want {
			t.Errorf("LineReader.ReadLineString() = %v, want:%s", got, tt.want)
		}
	}

	if got, err := r.PeekLineString(); err != io.EOF {
		t.Errorf("LineReader.ReadLineString() = %v, error = %s", got, err)
	}
	if got, err := r.PeekLineString(); err != io.EOF {
		t.Errorf("LineReader.ReadLineString() = %v, error = %s", got, err)
	}
}
