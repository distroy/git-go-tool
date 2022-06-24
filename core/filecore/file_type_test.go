/*
 * Copyright (C) distroy
 */

package filecore

import (
	"os"
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/strcore"
)

func TestFile_ReadLines(t *testing.T) {
	tests := []struct {
		name    string
		file    *File
		want    []string
		wantErr bool
	}{
		{
			file: NewTestFile("abc", []byte("aaa\nbbb\r\nccc\r\n")),
			want: []string{"aaa", "bbb", "ccc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			got, err := f.ReadLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.ReadLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File.ReadLines() = %v, want %v", got, tt.want)
			}
			got, err = f.ReadLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.ReadLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File.ReadLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Parse(t *testing.T) {
	text := `
package test
func Test() {
}
`
	tests := []struct {
		name    string
		file    *File
		want    bool
		wantErr bool
	}{
		{
			file: NewTestFile("test_file.go", strcore.StrToBytesUnsafe(text)),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			got, err := f.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != (got != nil) {
				t.Errorf("File.Parse() = %v, want %v", got, tt.want)
			}

			got, err = f.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != (got != nil) {
				t.Errorf("File.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Read(t *testing.T) {
	filename := MustTempFile()
	defer os.Remove(filename)

	text := strcore.StrToBytesUnsafe(`
package test
func Test() {
}
`)
	if err := WriteFile(filename, text, os.ModePerm); err != nil {
		t.Errorf("WriteFile() error = %v", err)
		return
	}

	tests := []struct {
		name    string
		file    *File
		want    []byte
		wantErr bool
	}{
		{
			file: NewTestFile(filename, text),
			want: text,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			got, err := f.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
