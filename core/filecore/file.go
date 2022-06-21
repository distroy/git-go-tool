/*
 * Copyright (C) distroy
 */

package filecore

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/distroy/git-go-tool/core/iocore"
)

type File struct {
	Path  string
	Name  string
	fset  *token.FileSet
	file  *ast.File
	data  []byte
	lines []string
}

func (f *File) Read() ([]byte, error) {
	if f.data != nil {
		return f.data, nil
	}

	data, err := iocore.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []byte{}
	}

	f.data = data
	return data, nil
}

func (f *File) ReadLines() ([]string, error) {
	if f.lines != nil {
		return f.lines, nil
	}

	data, err := f.Read()
	if err != nil {
		return nil, err
	}

	// log.Printf(" === %s", data)

	r := iocore.NewLineReader(bytes.NewBuffer(data))
	lines, err := r.ReadAllLineStrings()
	if err != nil {
		return nil, err
	}
	// log.Printf(" === %v", lines)

	if lines == nil {
		lines = []string{}
	}

	f.lines = lines
	return lines, nil
}

func (f *File) Parse() (*ast.File, error) {
	if f.file != nil {
		return f.file, nil
	}

	data, err := f.Read()
	if err != nil {
		return nil, err
	}

	fset := f.fset
	if fset == nil {
		fset = token.NewFileSet()
	}

	file, err := parser.ParseFile(fset, f.Path, data, 0)
	if err != nil {
		return nil, err
	}

	f.fset = fset
	f.file = file
	return f.file, nil
}

func (f *File) Position(p token.Pos) token.Position {
	return f.fset.Position(p)
}
