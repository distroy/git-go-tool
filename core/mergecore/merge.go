/*
 * Copyright (C) distroy
 */

package mergecore

import (
	"fmt"
	"reflect"

	"github.com/distroy/git-go-tool/core/refcore"
)

func valueOf(v interface{}) reflect.Value {
	if vv, ok := v.(reflect.Value); ok {
		return vv
	}
	return reflect.ValueOf(v)
}

func Merge(target, source interface{}) error {
	targetV := valueOf(target)
	sourceV := valueOf(source)

	if targetV.Kind() == reflect.Ptr {
		targetV = targetV.Elem()
	}
	if sourceV.Kind() == reflect.Ptr {
		sourceV = sourceV.Elem()
	}

	if !targetV.CanAddr() {
		return fmt.Errorf("target is not addressable. %s", targetV.Type().String())
	}

	if targetV.Type() != sourceV.Type() {
		return fmt.Errorf("the types of target and source is different. target:%s, source:%s",
			targetV.Type().String(), sourceV.Type().String())
	}

	// if targetV.Kind() != reflect.Struct {
	// 	return fmt.Errorf("target should be a pointer to struct. %T", valueOf(target).Type().String())
	// }

	return merge(targetV, sourceV)
}

func merge(target, source reflect.Value) error {
	switch target.Kind() {
	case reflect.Ptr:
		return mergePtr(target, source)

	case reflect.Interface, reflect.Func:
		return mergeInterface(target, source)

	case reflect.Slice:
		return mergeSlice(target, source)

	case reflect.Map:
		return mergeMap(target, source)

	case reflect.Struct:
		return mergeStruct(target, source)

	case reflect.Bool, reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.UnsafePointer:

		return mergeBaseType(target, source)
	}

	return fmt.Errorf("Unsuported type for merge. %s", target.Type().String())
}

func mergePtr(target, source reflect.Value) error {
	if source.IsNil() {
		return nil
	}

	if target.IsNil() {
		target.Set(source)
		return nil
	}

	return merge(target.Elem(), source.Elem())
}

func mergeInterface(target, source reflect.Value) error {
	if !source.IsNil() {
		target.Set(source)
	}
	return nil
}

func mergeBaseType(target, source reflect.Value) error {
	if !refcore.IsValZero(source) {
		target.Set(source)
	}
	return nil
}

func mergeSlice(target, source reflect.Value) error {
	res := target
	for i, l := 0, source.Len(); i < l; i++ {
		item := source.Index(i)
		res = reflect.Append(res, item)
	}
	target.Set(res)
	return nil
}

func mergeMap(target, source reflect.Value) error {
	it := source.MapRange()
	for it.Next() {
		key := it.Key()
		sVal := it.Value()

		tVal := target.MapIndex(key)
		if !tVal.IsValid() {
			target.SetMapIndex(key, sVal)
			continue
		}

		err := merge(tVal, sVal)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeStruct(target, source reflect.Value) error {
	for i, l := 0, source.NumField(); i < l; i++ {
		tField := target.Field(i)
		sField := source.Field(i)
		if err := merge(tField, sField); err != nil {
			return err
		}
	}

	return nil
}
