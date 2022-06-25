/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"flag"
	"reflect"
	"strconv"
	"time"
)

var (
	typeDuration = reflect.TypeOf(time.Duration(0))
	typeFunc     = reflect.TypeOf((func(string) error)(nil))
)
var (
	typeBool = reflect.TypeOf(bool(false))

	typeInt    = reflect.TypeOf(int(0))
	typeInt64  = reflect.TypeOf(int64(0))
	typeUint   = reflect.TypeOf(uint(0))
	typeUint64 = reflect.TypeOf(uint64(0))

	typeFloat32 = reflect.TypeOf(float32(0))
	typeFloat64 = reflect.TypeOf(float64(0))

	typeString = reflect.TypeOf(string(""))
)

var (
	typeBools = reflect.TypeOf([]bool(nil))

	typeInts    = reflect.TypeOf([]int(nil))
	typeInt64s  = reflect.TypeOf([]int64(nil))
	typeUints   = reflect.TypeOf([]uint(nil))
	typeUint64s = reflect.TypeOf([]uint64(nil))

	typeFloat32s = reflect.TypeOf([]float32(nil))
	typeFloat64s = reflect.TypeOf([]float64(nil))

	typeStrings = reflect.TypeOf([]string(nil))
)

type fillFlagFuncType = func(cmd *flag.FlagSet, val reflect.Value, m *fieldTags)

var fillFlagFuncMap = map[reflect.Type]fillFlagFuncType{
	typeDuration: fillFlagDuration,
	typeFunc:     fillFlagFunc,

	typeBool: fillFlagBool,

	typeString:  fillFlagString,
	typeInt:     fillFlagInt,
	typeInt64:   fillFlagInt64,
	typeUint:    fillFlagUint,
	typeUint64:  fillFlagUint64,
	typeFloat32: fillFlagFloat32,
	typeFloat64: fillFlagFloat64,

	typeStrings:  fillFlagStrings,
	typeInts:     fillFlagInts,
	typeInt64s:   fillFlagInt64s,
	typeUints:    fillFlagUints,
	typeUint64s:  fillFlagUint64s,
	typeFloat32s: fillFlagFloat32s,
	typeFloat64s: fillFlagFloat64s,
}

// bool
func fillFlagBool(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newBoolValue(val.Addr().Interface().(*bool))
	if len(m.Default) != 0 {
		p.Set(m.Default)
	}
	// log.Printf(" === %s, %d, %s", m.Name, *p, m.Default)
	cmd.Var(p, m.Name, m.Usage)
}

// string
func fillFlagString(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*string)
	if len(m.Default) != 0 {
		*p = m.Default
	}
	cmd.StringVar(p, m.Name, *p, m.Usage)
}

// int
func fillFlagInt(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*int)
	if len(m.Default) != 0 {
		v, err := strconv.Atoi(m.Default)
		if err == nil {
			*p = v
		}
	}
	// log.Printf(" === %s, %d, %s", m.Name, *p, m.Default)
	cmd.IntVar(p, m.Name, *p, m.Usage)
}

func fillFlagInt64(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*int64)
	if len(m.Default) != 0 {
		v, err := strconv.ParseInt(m.Default, 0, 64)
		if err == nil {
			*p = v
		}
	}
	cmd.Int64Var(p, m.Name, *p, m.Usage)
}

// uint
func fillFlagUint(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*uint)
	if len(m.Default) != 0 {
		v, err := strconv.ParseUint(m.Default, 0, strconv.IntSize)
		if err == nil {
			*p = uint(v)
		}
	}
	cmd.UintVar(p, m.Name, *p, m.Usage)
}

func fillFlagUint64(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*uint64)
	if len(m.Default) != 0 {
		v, err := strconv.ParseUint(m.Default, 0, 64)
		if err == nil {
			*p = v
		}
	}
	cmd.Uint64Var(p, m.Name, *p, m.Usage)
}

// float
func fillFlagFloat32(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newFloat32Value(val.Addr().Interface().(*float32))
	if len(m.Default) != 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

func fillFlagFloat64(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*float64)
	if len(m.Default) != 0 {
		v, err := strconv.ParseFloat(m.Default, 64)
		if err == nil {
			*p = v
		}
	}
	cmd.Float64Var(p, m.Name, *p, m.Usage)
}

// duration
func fillFlagDuration(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Addr().Interface().(*time.Duration)
	if len(m.Default) != 0 {
		v, err := time.ParseDuration(m.Default)
		if err == nil {
			*p = v
		}
	}
	cmd.DurationVar(p, m.Name, *p, m.Usage)
}

// func
func fillFlagFunc(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := val.Interface().(func(string) error)
	if len(m.Default) != 0 {
		p(m.Default)
	}
	cmd.Func(m.Name, m.Usage, p)
}

// strings
func fillFlagStrings(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newStringsValue(val.Addr().Interface().(*[]string))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

// ints
func fillFlagInts(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newIntsValue(val.Addr().Interface().(*[]int))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

func fillFlagInt64s(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newInt64sValue(val.Addr().Interface().(*[]int64))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

// uints
func fillFlagUints(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newUintsValue(val.Addr().Interface().(*[]uint))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

func fillFlagUint64s(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newUint64sValue(val.Addr().Interface().(*[]uint64))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

// floats
func fillFlagFloat32s(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newFloat32sValue(val.Addr().Interface().(*[]float32))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}

func fillFlagFloat64s(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	p := newFloat64sValue(val.Addr().Interface().(*[]float64))
	if len(m.Default) != 0 && len(*p) == 0 {
		p.Set(m.Default)
	}
	cmd.Var(p, m.Name, m.Usage)
}
