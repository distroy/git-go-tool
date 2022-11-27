/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"reflect"
	"time"

	"github.com/distroy/git-go-tool/core/ptrcore"
)

var (
	typeDuration = reflect.TypeOf(time.Duration(0))
	typeFunc     = reflect.TypeOf((func(string) error)(nil))
)
var (
	typeBool   = reflect.TypeOf(bool(false))
	typeString = reflect.TypeOf(string(""))

	typeInt    = reflect.TypeOf(int(0))
	typeInt64  = reflect.TypeOf(int64(0))
	typeUint   = reflect.TypeOf(uint(0))
	typeUint64 = reflect.TypeOf(uint64(0))

	typeFloat32 = reflect.TypeOf(float32(0))
	typeFloat64 = reflect.TypeOf(float64(0))
)

var (
	typeDurationPtr = reflect.TypeOf(ptrcore.NewDuration(0))

	typeBoolPtr   = reflect.TypeOf(ptrcore.NewBool(false))
	typeStringPtr = reflect.TypeOf(ptrcore.NewString(""))

	typeIntPtr    = reflect.TypeOf(ptrcore.NewInt(0))
	typeInt64Ptr  = reflect.TypeOf(ptrcore.NewInt64(0))
	typeUintPtr   = reflect.TypeOf(ptrcore.NewUint(0))
	typeUint64Ptr = reflect.TypeOf(ptrcore.NewUint64(0))

	typeFloat32Ptr = reflect.TypeOf(ptrcore.NewFloat32(0))
	typeFloat64Ptr = reflect.TypeOf(ptrcore.NewFloat64(0))
)

var (
	typeBools   = reflect.TypeOf([]bool(nil))
	typeStrings = reflect.TypeOf([]string(nil))

	typeInts    = reflect.TypeOf([]int(nil))
	typeInt64s  = reflect.TypeOf([]int64(nil))
	typeUints   = reflect.TypeOf([]uint(nil))
	typeUint64s = reflect.TypeOf([]uint64(nil))

	typeFloat32s = reflect.TypeOf([]float32(nil))
	typeFloat64s = reflect.TypeOf([]float64(nil))
)

type fillFlagFuncType = func(val reflect.Value) Value

var fillFlagFuncMap = map[reflect.Type]fillFlagFuncType{
	typeDuration: func(val reflect.Value) Value { return newDurationValue(val.Addr().Interface().(*time.Duration)) },
	typeFunc:     func(val reflect.Value) Value { return newFuncValue(val.Interface().(func(s string) error)) },

	typeBool:   func(val reflect.Value) Value { return newBoolValue(val.Addr().Interface().(*bool)) },
	typeString: func(val reflect.Value) Value { return newStringValue(val.Addr().Interface().(*string)) },

	typeInt:     func(val reflect.Value) Value { return newIntValue(val.Addr().Interface().(*int)) },
	typeInt64:   func(val reflect.Value) Value { return newInt64Value(val.Addr().Interface().(*int64)) },
	typeUint:    func(val reflect.Value) Value { return newUintValue(val.Addr().Interface().(*uint)) },
	typeUint64:  func(val reflect.Value) Value { return newUint64Value(val.Addr().Interface().(*uint64)) },
	typeFloat32: func(val reflect.Value) Value { return newFloat32Value(val.Addr().Interface().(*float32)) },
	typeFloat64: func(val reflect.Value) Value { return newFloat64Value(val.Addr().Interface().(*float64)) },

	typeStrings:  func(val reflect.Value) Value { return newStringsValue(val.Addr().Interface().(*[]string)) },
	typeInts:     func(val reflect.Value) Value { return newIntsValue(val.Addr().Interface().(*[]int)) },
	typeInt64s:   func(val reflect.Value) Value { return newInt64sValue(val.Addr().Interface().(*[]int64)) },
	typeUints:    func(val reflect.Value) Value { return newUintsValue(val.Addr().Interface().(*[]uint)) },
	typeUint64s:  func(val reflect.Value) Value { return newUint64sValue(val.Addr().Interface().(*[]uint64)) },
	typeFloat32s: func(val reflect.Value) Value { return newFloat32sValue(val.Addr().Interface().(*[]float32)) },
	typeFloat64s: func(val reflect.Value) Value { return newFloat64sValue(val.Addr().Interface().(*[]float64)) },

	typeDurationPtr: func(val reflect.Value) Value { return newDurationPtrValue(val.Addr().Interface().(**time.Duration)) },

	typeBoolPtr:   func(val reflect.Value) Value { return newBoolPtrValue(val.Addr().Interface().(**bool)) },
	typeStringPtr: func(val reflect.Value) Value { return newStringPtrValue(val.Addr().Interface().(**string)) },

	typeIntPtr:     func(val reflect.Value) Value { return newIntPtrValue(val.Addr().Interface().(**int)) },
	typeInt64Ptr:   func(val reflect.Value) Value { return newInt64PtrValue(val.Addr().Interface().(**int64)) },
	typeUintPtr:    func(val reflect.Value) Value { return newUintPtrValue(val.Addr().Interface().(**uint)) },
	typeUint64Ptr:  func(val reflect.Value) Value { return newUint64PtrValue(val.Addr().Interface().(**uint64)) },
	typeFloat32Ptr: func(val reflect.Value) Value { return newFloat32PtrValue(val.Addr().Interface().(**float32)) },
	typeFloat64Ptr: func(val reflect.Value) Value { return newFloat64PtrValue(val.Addr().Interface().(**float64)) },
}
