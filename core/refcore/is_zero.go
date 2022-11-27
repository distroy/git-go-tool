/*
 * Copyright (C) distroy
 */

package refcore

import (
	"math"
	"reflect"
)

func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}

	if vv, ok := v.(reflect.Value); ok {
		return IsValZero(vv)
	}

	vv := reflect.ValueOf(v)
	return IsValZero(vv)
}

func IsValZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true

	case reflect.Bool:
		return !v.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0

	case reflect.Array:
		return isArrayValZero(v)

	case reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()

	case reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()

	case reflect.UnsafePointer:
		return v.IsNil()

	case reflect.String:
		return v.Len() == 0

	case reflect.Struct:
		return isStructValZero(v)

	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(&reflect.ValueError{Method: "ldref.IsZero", Kind: v.Kind()})
	}
}

func isArrayValZero(v reflect.Value) bool {
	for i := 0; i < v.Len(); i++ {
		if !IsValZero(v.Index(i)) {
			return false
		}
	}
	return true
}

func isStructValZero(v reflect.Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if !IsValZero(v.Field(i)) {
			return false
		}
	}
	return true
}
