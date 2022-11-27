/*
 * Copyright (C) distroy
 */

package ptrcore

import (
	"testing"
	"time"

	"github.com/distroy/git-go-tool/3rd/convey"
)

func TestNewBoolByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewBoolByPtr(nil), convey.ShouldBeNil)
		convey.So(NewBoolByPtr(nil, false), convey.ShouldResemble, NewBool(false))
		convey.So(NewBoolByPtr(nil, true), convey.ShouldResemble, NewBool(true))
		convey.So(NewBoolByPtr(NewBool(false)), convey.ShouldResemble, NewBool(false))
		convey.So(NewBoolByPtr(NewBool(true)), convey.ShouldResemble, NewBool(true))
	})
}

func TestNewByteByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewByteByPtr(nil), convey.ShouldBeNil)
		convey.So(NewByteByPtr(nil, 0), convey.ShouldResemble, NewByte(0))
		convey.So(NewByteByPtr(nil, 1), convey.ShouldResemble, NewByte(1))
		convey.So(NewByteByPtr(NewByte(1)), convey.ShouldResemble, NewByte(1))
		convey.So(NewByteByPtr(NewByte(100)), convey.ShouldResemble, NewByte(100))
	})
}

func TestNewRuneByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewRuneByPtr(nil), convey.ShouldBeNil)
		convey.So(NewRuneByPtr(nil, 0), convey.ShouldResemble, NewRune(0))
		convey.So(NewRuneByPtr(nil, 1), convey.ShouldResemble, NewRune(1))
		convey.So(NewRuneByPtr(NewRune(1)), convey.ShouldResemble, NewRune(1))
		convey.So(NewRuneByPtr(NewRune(100)), convey.ShouldResemble, NewRune(100))
	})
}

func TestNewIntByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewIntByPtr(nil), convey.ShouldBeNil)
		convey.So(NewIntByPtr(nil, 0), convey.ShouldResemble, NewInt(0))
		convey.So(NewIntByPtr(nil, 1), convey.ShouldResemble, NewInt(1))
		convey.So(NewIntByPtr(NewInt(1)), convey.ShouldResemble, NewInt(1))
		convey.So(NewIntByPtr(NewInt(100)), convey.ShouldResemble, NewInt(100))
		convey.So(NewIntByPtr(NewInt(-100)), convey.ShouldResemble, NewInt(-100))
	})
}

func TestNewInt8ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewInt8ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewInt8ByPtr(nil, 0), convey.ShouldResemble, NewInt8(0))
		convey.So(NewInt8ByPtr(nil, 1), convey.ShouldResemble, NewInt8(1))
		convey.So(NewInt8ByPtr(NewInt8(1)), convey.ShouldResemble, NewInt8(1))
		convey.So(NewInt8ByPtr(NewInt8(100)), convey.ShouldResemble, NewInt8(100))
		convey.So(NewInt8ByPtr(NewInt8(-100)), convey.ShouldResemble, NewInt8(-100))
	})
}

func TestNewInt16ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewInt16ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewInt16ByPtr(nil, 0), convey.ShouldResemble, NewInt16(0))
		convey.So(NewInt16ByPtr(nil, 1), convey.ShouldResemble, NewInt16(1))
		convey.So(NewInt16ByPtr(NewInt16(1)), convey.ShouldResemble, NewInt16(1))
		convey.So(NewInt16ByPtr(NewInt16(100)), convey.ShouldResemble, NewInt16(100))
		convey.So(NewInt16ByPtr(NewInt16(-100)), convey.ShouldResemble, NewInt16(-100))
	})
}

func TestNewInt32ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewInt32ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewInt32ByPtr(nil, 0), convey.ShouldResemble, NewInt32(0))
		convey.So(NewInt32ByPtr(nil, 1), convey.ShouldResemble, NewInt32(1))
		convey.So(NewInt32ByPtr(NewInt32(1)), convey.ShouldResemble, NewInt32(1))
		convey.So(NewInt32ByPtr(NewInt32(100)), convey.ShouldResemble, NewInt32(100))
		convey.So(NewInt32ByPtr(NewInt32(-100)), convey.ShouldResemble, NewInt32(-100))
	})
}

func TestNewInt64ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewInt64ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewInt64ByPtr(nil, 0), convey.ShouldResemble, NewInt64(0))
		convey.So(NewInt64ByPtr(nil, 1), convey.ShouldResemble, NewInt64(1))
		convey.So(NewInt64ByPtr(NewInt64(1)), convey.ShouldResemble, NewInt64(1))
		convey.So(NewInt64ByPtr(NewInt64(100)), convey.ShouldResemble, NewInt64(100))
		convey.So(NewInt64ByPtr(NewInt64(-100)), convey.ShouldResemble, NewInt64(-100))
	})
}

func TestNewUintByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUintByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUintByPtr(nil, 0), convey.ShouldResemble, NewUint(0))
		convey.So(NewUintByPtr(nil, 1), convey.ShouldResemble, NewUint(1))
		convey.So(NewUintByPtr(NewUint(1)), convey.ShouldResemble, NewUint(1))
		convey.So(NewUintByPtr(NewUint(100)), convey.ShouldResemble, NewUint(100))
	})
}

func TestNewUint8ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUint8ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUint8ByPtr(nil, 0), convey.ShouldResemble, NewUint8(0))
		convey.So(NewUint8ByPtr(nil, 1), convey.ShouldResemble, NewUint8(1))
		convey.So(NewUint8ByPtr(NewUint8(1)), convey.ShouldResemble, NewUint8(1))
		convey.So(NewUint8ByPtr(NewUint8(100)), convey.ShouldResemble, NewUint8(100))
	})
}

func TestNewUint16ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUint16ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUint16ByPtr(nil, 0), convey.ShouldResemble, NewUint16(0))
		convey.So(NewUint16ByPtr(nil, 1), convey.ShouldResemble, NewUint16(1))
		convey.So(NewUint16ByPtr(NewUint16(1)), convey.ShouldResemble, NewUint16(1))
		convey.So(NewUint16ByPtr(NewUint16(100)), convey.ShouldResemble, NewUint16(100))
	})
}

func TestNewUint32ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUint32ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUint32ByPtr(nil, 0), convey.ShouldResemble, NewUint32(0))
		convey.So(NewUint32ByPtr(nil, 1), convey.ShouldResemble, NewUint32(1))
		convey.So(NewUint32ByPtr(NewUint32(1)), convey.ShouldResemble, NewUint32(1))
		convey.So(NewUint32ByPtr(NewUint32(100)), convey.ShouldResemble, NewUint32(100))
	})
}

func TestNewUint64ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUint64ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUint64ByPtr(nil, 0), convey.ShouldResemble, NewUint64(0))
		convey.So(NewUint64ByPtr(nil, 1), convey.ShouldResemble, NewUint64(1))
		convey.So(NewUint64ByPtr(NewUint64(1)), convey.ShouldResemble, NewUint64(1))
		convey.So(NewUint64ByPtr(NewUint64(100)), convey.ShouldResemble, NewUint64(100))
	})
}

func TestNewUintptrByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewUintptrByPtr(nil), convey.ShouldBeNil)
		convey.So(NewUintptrByPtr(nil, 0), convey.ShouldResemble, NewUintptr(0))
		convey.So(NewUintptrByPtr(nil, 1), convey.ShouldResemble, NewUintptr(1))
		convey.So(NewUintptrByPtr(NewUintptr(1)), convey.ShouldResemble, NewUintptr(1))
		convey.So(NewUintptrByPtr(NewUintptr(100)), convey.ShouldResemble, NewUintptr(100))
	})
}

func TestNewFloat32ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewFloat32ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewFloat32ByPtr(nil, 0), convey.ShouldResemble, NewFloat32(0))
		convey.So(NewFloat32ByPtr(nil, 1), convey.ShouldResemble, NewFloat32(1))
		convey.So(NewFloat32ByPtr(NewFloat32(1)), convey.ShouldResemble, NewFloat32(1))
		convey.So(NewFloat32ByPtr(NewFloat32(100)), convey.ShouldResemble, NewFloat32(100))
		convey.So(NewFloat32ByPtr(NewFloat32(123.123)), convey.ShouldResemble, NewFloat32(123.123))
	})
}

func TestNewFloat64ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewFloat64ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewFloat64ByPtr(nil, 0), convey.ShouldResemble, NewFloat64(0))
		convey.So(NewFloat64ByPtr(nil, 1), convey.ShouldResemble, NewFloat64(1))
		convey.So(NewFloat64ByPtr(NewFloat64(1)), convey.ShouldResemble, NewFloat64(1))
		convey.So(NewFloat64ByPtr(NewFloat64(100)), convey.ShouldResemble, NewFloat64(100))
		convey.So(NewFloat64ByPtr(NewFloat64(-100)), convey.ShouldResemble, NewFloat64(-100))
		convey.So(NewFloat64ByPtr(NewFloat64(123.123)), convey.ShouldResemble, NewFloat64(123.123))
	})
}

func TestNewStringByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewStringByPtr(nil), convey.ShouldBeNil)
		convey.So(NewStringByPtr(nil, "a"), convey.ShouldResemble, NewString("a"))
		convey.So(NewStringByPtr(NewString("1")), convey.ShouldResemble, NewString("1"))
	})
}

func TestNewComplex64ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewComplex64ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewComplex64ByPtr(nil, 0), convey.ShouldResemble, NewComplex64(0))
		convey.So(NewComplex64ByPtr(nil, 1), convey.ShouldResemble, NewComplex64(1))
		convey.So(NewComplex64ByPtr(NewComplex64(1)), convey.ShouldResemble, NewComplex64(1))
		convey.So(NewComplex64ByPtr(NewComplex64(100)), convey.ShouldResemble, NewComplex64(100))
		convey.So(NewComplex64ByPtr(NewComplex64(-100)), convey.ShouldResemble, NewComplex64(-100))
	})
}

func TestNewComplex128ByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewComplex128ByPtr(nil), convey.ShouldBeNil)
		convey.So(NewComplex128ByPtr(nil, 0), convey.ShouldResemble, NewComplex128(0))
		convey.So(NewComplex128ByPtr(nil, 1), convey.ShouldResemble, NewComplex128(1))
		convey.So(NewComplex128ByPtr(NewComplex128(1)), convey.ShouldResemble, NewComplex128(1))
		convey.So(NewComplex128ByPtr(NewComplex128(100)), convey.ShouldResemble, NewComplex128(100))
		convey.So(NewComplex128ByPtr(NewComplex128(-100)), convey.ShouldResemble, NewComplex128(-100))
	})
}

func TestNewTimeByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewTimeByPtr(nil), convey.ShouldBeNil)
		convey.So(NewTimeByPtr(nil, time.Unix(0, 0)), convey.ShouldResemble, NewTime(time.Unix(0, 0)))
		convey.So(NewTimeByPtr(nil, time.Unix(1, 0)), convey.ShouldResemble, NewTime(time.Unix(1, 0)))
		convey.So(NewTimeByPtr(NewTime(time.Unix(1, 0))), convey.ShouldResemble, NewTime(time.Unix(1, 0)))
		convey.So(NewTimeByPtr(NewTime(time.Unix(100, 0))), convey.ShouldResemble, NewTime(time.Unix(100, 0)))
		convey.So(NewTimeByPtr(NewTime(time.Unix(-100, 0))), convey.ShouldResemble, NewTime(time.Unix(-100, 0)))
	})
}

func TestNewDurationByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewDurationByPtr(nil), convey.ShouldBeNil)
		convey.So(NewDurationByPtr(nil, 0), convey.ShouldResemble, NewDuration(0))
		convey.So(NewDurationByPtr(nil, 1), convey.ShouldResemble, NewDuration(1))
		convey.So(NewDurationByPtr(NewDuration(1)), convey.ShouldResemble, NewDuration(1))
		convey.So(NewDurationByPtr(NewDuration(100)), convey.ShouldResemble, NewDuration(100))
		convey.So(NewDurationByPtr(NewDuration(-100)), convey.ShouldResemble, NewDuration(-100))
	})
}
