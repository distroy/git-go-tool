/*
 * Copyright (C) distroy
 */

package filecore

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

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

func (f *File) IsGo() bool     { return strings.HasSuffix(f.Name, ".go") }
func (f *File) IsGoTest() bool { return strings.HasSuffix(f.Name, "_test.go") }

func (f *File) MustRead() []byte {
	data, err := f.Read()
	if err != nil {
		panic(fmt.Sprintf("read file fail. file:%s, err:%v", f.Path, err))
	}
	return data
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

func (f *File) MustReadLines() []string {
	lines, err := f.ReadLines()
	if err != nil {
		panic(fmt.Sprintf("read file lines fail. file:%s, err:%v", f.Path, err))
	}
	return lines
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

func (f *File) MustParse() *ast.File {
	file, err := f.Parse()
	if err != nil {
		panic(fmt.Sprintf("parse file fail. file:%s, err:%v", f.Path, err))
	}
	return file
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

	mode := parser.ParseComments
	file, err := parser.ParseFile(fset, f.Path, data, mode)
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
