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

type LineReader struct {
	reader    io.Reader
	buffer    []byte
	maxSize   int
	tokenPos  int
	tokenEnd  int
	tokenNext int
	bufferEnd int
	err       error
}

func NewLineReader(r io.Reader) *LineReader {
	maxSize := 4096
	return newLineReader(r, maxSize)
}

func newLineReader(r io.Reader, maxSize int) *LineReader {
	return &LineReader{
		reader:    r,
		buffer:    make([]byte, maxSize),
		maxSize:   maxSize,
		tokenPos:  0,
		tokenEnd:  -1,
		tokenNext: 0,
		bufferEnd: 0,
		err:       nil,
	}
}

func (r *LineReader) ReadAllLineStrings() ([]string, error) {
	res := make([]string, 0, 10)

	for {
		l, err := r.ReadLineString()
		if err != nil {
			if err == io.EOF {
				break
			}

			return res, err
		}

		res = append(res, l)
	}

	return res, nil
}

func (r *LineReader) PeekLineString() (string, error) {
	b, err := r.PeekLine()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *LineReader) ReadLineString() (string, error) {
	b, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *LineReader) PeekLine() ([]byte, error) {
	err := r.read()
	if err != nil {
		return nil, err
	}

	if r.tokenEnd >= 0 {
		return r.buffer[r.tokenPos:r.tokenEnd], nil
	}

	return nil, r.err
}

func (r *LineReader) ReadLine() ([]byte, error) {
	err := r.read()
	if err != nil {
		return nil, err
	}

	if end := r.tokenEnd; end >= 0 {
		pos := r.tokenPos
		r.tokenEnd = -1
		return r.buffer[pos:end], nil
	}

	return nil, r.err
}

func (r *LineReader) read() error {
	if r.tokenEnd >= 0 {
		return nil
	}

	idx := r.indexToken(r.tokenNext)
	if idx >= 0 {
		r.setToken(idx)
		return nil
	}

	if r.err != nil {
		if r.err == io.EOF && r.bufferEnd > r.tokenNext {
			r.setToken(r.bufferEnd)
			return nil
		}
		return r.err
	}

	return r.readLineLoop()
}

func (r *LineReader) readLineLoop() error {
	r.moveFront()

	for {
		end := r.bufferEnd
		buf := r.buffer[end:]
		n, err := r.reader.Read(buf)
		if err != nil {
			r.err = err
			break
		}

		r.bufferEnd += n
		idx := r.indexToken(end)
		if idx >= 0 {
			r.setToken(idx)
			return nil
		}

		if r.bufferEnd >= r.maxSize {
			r.err = ErrOverMaxSize
			return r.err
		}
	}

	if r.err == io.EOF {
		if r.bufferEnd > r.tokenNext {
			r.setToken(r.bufferEnd)
			return nil
		}
	}
	return r.err
}

func (r *LineReader) indexToken(pos int) int {
	end := r.bufferEnd
	if pos >= end {
		return -1
		// log.Printf(" *** pos:%d, token pos:%d, token end:%d, token next:%d, buffer end:%d, err:%v",
		// 	pos, r.tokenPos, r.tokenEnd, r.tokenNext, r.bufferEnd, r.err)
	}
	idx := bytes.IndexByte(r.buffer[pos:end], '\n')
	if idx < 0 {
		return -1
	}
	return idx + pos
}

func (r *LineReader) setToken(idx int) {
	r.tokenPos = r.tokenNext
	r.tokenEnd = idx
	r.tokenNext = idx + 1

	if r.tokenEnd > r.tokenPos && r.buffer[r.tokenEnd-1] == '\r' {
		r.tokenEnd--
	}
}

func (r *LineReader) moveFront() {
	if r.tokenNext < r.bufferEnd {
		copy(r.buffer, r.buffer[r.tokenNext:r.bufferEnd])
	}

	r.bufferEnd -= r.tokenNext
	r.tokenNext = 0
}
