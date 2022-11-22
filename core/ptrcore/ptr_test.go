/*
 * Copyright (C) distroy
 */

package ptrcore

import (
	"testing"
	"time"

	"github.com/distroy/git-go-tool/3rd/convey"
)

func TestGetBool(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetBool(nil), convey.ShouldEqual, false)
		convey.So(GetBool(nil, true), convey.ShouldEqual, true)
		convey.So(GetBool(NewBool(true), false), convey.ShouldEqual, true)
	})
}

func TestGetByte(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetByte(nil), convey.ShouldEqual, 0)
		convey.So(GetByte(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetByte(NewByte(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetRune(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetRune(nil), convey.ShouldEqual, 0)
		convey.So(GetRune(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetRune(NewRune(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetInt(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetInt(nil), convey.ShouldEqual, 0)
		convey.So(GetInt(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetInt(NewInt(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetInt8(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetInt8(nil), convey.ShouldEqual, 0)
		convey.So(GetInt8(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetInt8(NewInt8(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetInt16(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetInt16(nil), convey.ShouldEqual, 0)
		convey.So(GetInt16(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetInt16(NewInt16(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetInt32(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetInt32(nil), convey.ShouldEqual, 0)
		convey.So(GetInt32(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetInt32(NewInt32(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetInt64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetInt64(nil), convey.ShouldEqual, 0)
		convey.So(GetInt64(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetInt64(NewInt64(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUint(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUint(nil), convey.ShouldEqual, 0)
		convey.So(GetUint(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUint(NewUint(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUint8(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUint8(nil), convey.ShouldEqual, 0)
		convey.So(GetUint8(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUint8(NewUint8(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUint16(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUint16(nil), convey.ShouldEqual, 0)
		convey.So(GetUint16(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUint16(NewUint16(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUint32(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUint32(nil), convey.ShouldEqual, 0)
		convey.So(GetUint32(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUint32(NewUint32(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUint64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUint64(nil), convey.ShouldEqual, 0)
		convey.So(GetUint64(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUint64(NewUint64(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetUintptr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetUintptr(nil), convey.ShouldEqual, 0)
		convey.So(GetUintptr(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetUintptr(NewUintptr(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetFloat32(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetFloat32(nil), convey.ShouldEqual, 0)
		convey.So(GetFloat32(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetFloat32(NewFloat32(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetFloat64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetFloat64(nil), convey.ShouldEqual, 0)
		convey.So(GetFloat64(nil, 1), convey.ShouldEqual, 1)
		convey.So(GetFloat64(NewFloat64(100), 0), convey.ShouldEqual, 100)
	})
}

func TestGetString(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetString(nil), convey.ShouldEqual, "")
		convey.So(GetString(nil, "1"), convey.ShouldEqual, "1")
		convey.So(GetString(NewString("100"), "0"), convey.ShouldEqual, "100")
	})
}

func TestGetComplex64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetComplex64(nil), convey.ShouldEqual, complex(0, 0))
		convey.So(GetComplex64(nil, complex(1, 1)), convey.ShouldEqual, complex(1, 1))
		convey.So(GetComplex64(NewComplex64(complex(100, 1)), complex(0, 1)), convey.ShouldEqual, complex(100, 1))
	})
}

func TestGetComplex128(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetComplex128(nil), convey.ShouldEqual, complex(0, 0))
		convey.So(GetComplex128(nil, complex(1, 1)), convey.ShouldEqual, complex(1, 1))
		convey.So(GetComplex128(NewComplex128(complex(100, 1)), complex(0, 1)), convey.ShouldEqual, complex(100, 1))
	})
}

func TestGetTime(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetTime(nil), convey.ShouldEqual, time.Time{})
		convey.So(GetTime(nil, time.Unix(10, 0)), convey.ShouldEqual, time.Unix(10, 0))
		convey.So(GetTime(NewTime(time.Unix(10, 0)), time.Unix(0, 0)), convey.ShouldEqual, time.Unix(10, 0))
	})
}

func TestGetDuration(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(GetDuration(nil), convey.ShouldEqual, time.Duration(0))
		convey.So(GetDuration(nil, 1), convey.ShouldEqual, time.Duration(1))
		convey.So(GetDuration(NewDuration(100), time.Duration(0)), convey.ShouldEqual, time.Duration(100))
	})
}
