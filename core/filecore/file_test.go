/*
 * Copyright (C) distroy
 */

package filecore

import (
	"os"
	"reflect"
	"testing"
)

func Test_TempFile_WriteFile_ReadFile(t *testing.T) {
	filename := MustTempFile()
	defer os.Remove(filename)
	t.Logf("temp file: %s", filename)

	type args struct {
		name string
	}
	tests := []struct {
		name string
		want []byte
	}{
		{
			want: []byte("111\n2222\n333\naaaa\nbbbb"),
		},
		{
			want: []byte("111"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(filename, tt.want, os.ModePerm); err != nil {
				t.Errorf("WriteFile() error = %v", err)
				return
			}

			got, err := ReadFile(filename)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
