/*
 * Copyright (C) distroy
 */

package iocore

import (
	"bytes"
	"fmt"
	"io"
)

var (
	ErrOverMaxSize = fmt.Errorf("over max size")
)

type LineReader interface {
	Read() ([]byte, error)
	ReadString() (string, error)

	Peek() ([]byte, error)
	PeekString() (string, error)
}

type lineReader struct {
	reader      io.Reader
	token       []byte
	buffer      []byte
	maxSize     int
	tokenPos    int
	bufferPos   int
	lastLineEnd byte
	err         error
}

func NewLineReader(r io.Reader) *lineReader {
	maxSize := 4096
	return &lineReader{
		reader:    r,
		token:     make([]byte, maxSize),
		buffer:    make([]byte, maxSize),
		maxSize:   maxSize,
		tokenPos:  -1,
		bufferPos: 0,
		err:       nil,
	}
}

func (r *lineReader) PeekString() (string, error) {
	b, err := r.Peek()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *lineReader) ReadString() (string, error) {
	b, err := r.Read()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *lineReader) Peek() ([]byte, error) {
	err := r.read()
	if err != nil {
		return nil, err
	}

	if r.tokenPos >= 0 {
		return r.token[:r.tokenPos], nil
	}

	return nil, r.err
}

func (r *lineReader) Read() ([]byte, error) {
	err := r.read()
	if err != nil {
		return nil, err
	}

	if pos := r.tokenPos; pos >= 0 {
		r.tokenPos = -1
		return r.token[:pos], nil
	}

	return nil, r.err
}

func (r *lineReader) read() error {
	if r.tokenPos >= 0 {
		return nil
	}

	if r.err != nil {
		if r.err == io.EOF && r.bufferPos > 0 {
			r.tokenPos = r.bufferPos
			return nil
		}
		return r.err
	}

	idx := r.indexToken(0)
	if idx >= 0 {
		r.copyToken(idx)
		return nil
	}

	return r.readLineLoop()
}

func (r *lineReader) readLineLoop() error {
	for {
		pos := r.bufferPos
		buf := r.buffer[pos:]
		n, err := r.reader.Read(buf)
		if err != nil {
			r.err = err
			break
		}

		r.bufferPos += n
		idx := r.indexToken(pos)
		if idx >= 0 {
			r.copyToken(idx)
			return nil
		}

		if r.bufferPos >= r.maxSize {
			r.err = ErrOverMaxSize
			return r.err
		}
	}

	if r.err == io.EOF {
		if r.bufferPos > 0 {
			r.copyToken(r.bufferPos)
		}
		return nil
	}
	return r.err
}

func (r *lineReader) indexToken(pos int) int {
	end := r.bufferPos
	return bytes.IndexByte(r.buffer[pos:end], '\n')
}

func (r *lineReader) copyToken(pos int) {
	copy(r.token[:pos], r.buffer[:pos])
	r.tokenPos = pos
	if r.token[pos-1] == '\r' {
		r.tokenPos--
	}

	r.lastLineEnd = r.buffer[pos]

	pos++
	if pos >= r.bufferPos {
		r.bufferPos = 0
		return
	}

	copy(r.buffer, r.buffer[pos:r.bufferPos])
	r.bufferPos -= pos
}
