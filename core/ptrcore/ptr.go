/*
 * Copyright (C) distroy
 */

package ptrcore

import "time"

func GetBool(p *bool, def ...bool) bool {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func GetByte(p *byte, def ...byte) byte {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetRune(p *rune, def ...rune) rune {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetInt(p *int, def ...int) int {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetInt8(p *int8, def ...int8) int8 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetInt16(p *int16, def ...int16) int16 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetInt32(p *int32, def ...int32) int32 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetInt64(p *int64, def ...int64) int64 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUint(p *uint, def ...uint) uint {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUint8(p *uint8, def ...uint8) uint8 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUint16(p *uint16, def ...uint16) uint16 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUint32(p *uint32, def ...uint32) uint32 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUint64(p *uint64, def ...uint64) uint64 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetUintptr(p *uintptr, def ...uintptr) uintptr {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetFloat32(p *float32, def ...float32) float32 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetFloat64(p *float64, def ...float64) float64 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetString(p *string, def ...string) string {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func GetComplex64(p *complex64, def ...complex64) complex64 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return complex(0, 0)
}

func GetComplex128(p *complex128, def ...complex128) complex128 {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return complex(0, 0)
}

func GetTime(p *time.Time, def ...time.Time) time.Time {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func GetDuration(p *time.Duration, def ...time.Duration) time.Duration {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
