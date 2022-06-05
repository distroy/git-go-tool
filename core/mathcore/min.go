/*
 * Copyright (C) distroy
 */

package mathcore

func MinInt(n int, args ...int) int {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinInt8(n int8, args ...int8) int8 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinInt16(n int16, args ...int16) int16 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinInt32(n int32, args ...int32) int32 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinInt64(n int64, args ...int64) int64 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinUint(n uint, args ...uint) uint {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinUint8(n uint8, args ...uint8) uint8 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinUint16(n uint16, args ...uint16) uint16 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinUint32(n uint32, args ...uint32) uint32 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func MinUint64(n uint64, args ...uint64) uint64 {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}
