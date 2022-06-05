/*
 * Copyright (C) distroy
 */

package mathcore

func MaxInt(n int, args ...int) int {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxInt8(n int8, args ...int8) int8 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxInt16(n int16, args ...int16) int16 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxInt32(n int32, args ...int32) int32 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxInt64(n int64, args ...int64) int64 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxUint(n uint, args ...uint) uint {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxUint8(n uint8, args ...uint8) uint8 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxUint16(n uint16, args ...uint16) uint16 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxUint32(n uint32, args ...uint32) uint32 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func MaxUint64(n uint64, args ...uint64) uint64 {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}
