/*
 * Copyright (C) distroy
 */

package iocore

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

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

func TestLineReader_ReadAllLineStrings(t *testing.T) {
	type fields struct {
		reader     io.Reader
		bufferSize int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			fields: fields{
				reader: strings.NewReader("1111\r\n\n2222\n3333\r\n"),
			},
			want: []string{
				"1111",
				"",
				"2222",
				"3333",
			},
		},
		{
			fields: fields{
				reader: bytes.NewBuffer([]byte("1111\r\n\n2222\n3333\n")),
			},
			want: []string{
				"1111",
				"",
				"2222",
				"3333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]LineReaderOption, 0)
			if tt.fields.bufferSize > 0 {
				opts = append(opts, LineReaderBufferSize(tt.fields.bufferSize))
			}

			r := NewLineReader(tt.fields.reader, opts...)
			got, err := r.ReadAllLineStrings()
			if (err != nil) != tt.wantErr {
				t.Errorf("LineReader.ReadAllLineStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LineReader.ReadAllLineStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
