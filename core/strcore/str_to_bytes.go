/*
 * Copyright (C) distroy
 */

package strcore

import (
	"reflect"
	"unsafe"
)

func BytesToStr(b []byte) string {
	return string(b)
}

func StrToBytes(s string) []byte {
	return []byte(s)
}

// BytesToStrUnsafe the source bytes cannot free
func BytesToStrUnsafe(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StrToBytesUnsafe the result bytes cannot be modified and the source string cannot free
func StrToBytesUnsafe(s string) []byte {
	// return *(*[]byte)(unsafe.Pointer(&s))

	type bytesFromString reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := &bytesFromString{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bh))
}
