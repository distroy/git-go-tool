/*
 * Copyright (C) distroy
 */

package iocore

type LineReaderOption = func(r *LineReader)

func LineReaderBufferSize(size int) LineReaderOption {
	return func(r *LineReader) {
		r.bufferSize = size
	}
}
