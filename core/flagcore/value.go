/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"bytes"
	"encoding/json"
	"flag"
	"strconv"
	"time"

	"github.com/distroy/git-go-tool/core/strcore"
)

type Value interface {
	flag.Value
}

type valueWithDefault interface {
	Value

	Default() string
}

func mustMarshalJson(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.Encode(v)

	s := b.Bytes()
	if l := len(s) - 1; l >= 0 && s[l] == '\n' {
		s = s[:l]
	}

	return strcore.BytesToStrUnsafe(s)
}

// duration
type durationValue time.Duration

func newDurationValue(p *time.Duration) *durationValue { return (*durationValue)(p) }
func (p *durationValue) String() string                { return time.Duration(*p).String() }
func (p *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*p = durationValue(v)
	return nil
}

// func
type funcValue func(string) error

func newFuncValue(f func(string) error) funcValue { return funcValue(f) }
func (f funcValue) Set(s string) error            { return f(s) }
func (f funcValue) String() string                { return "" }

// string
type stringValue string

func newStringValue(p *string) *stringValue { return (*stringValue)(p) }
func (p *stringValue) String() string       { return string(*p) }
func (p *stringValue) Set(s string) error {
	*p = stringValue(s)
	return nil
}

// bool
type boolValue bool

func newBoolValue(p *bool) *boolValue { return (*boolValue)(p) }
func (p *boolValue) String() string   { return strconv.FormatBool(bool(*p)) }
func (p *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	// log.Printf(" === %v, %v, %v", *p, v, err)
	if err != nil {
		return err
	}
	*p = boolValue(v)
	return nil
}

type boolFlag bool

func newBoolFlag(p *boolValue) *boolFlag { return (*boolFlag)(p) }
func (p *boolFlag) String() string       { return (*boolValue)(p).String() }
func (p *boolFlag) Set(s string) error   { return (*boolValue)(p).Set(s) }
func (p *boolFlag) IsBoolFlag() bool     { return true }

// func (p *boolValue) Get() interface{} { return bool(*p) }
// func (p *boolValue) IsBoolFlag() bool { return true }

// int
type intValue int

func newIntValue(p *int) *intValue { return (*intValue)(p) }
func (p *intValue) String() string { return strconv.Itoa(int(*p)) }
func (p *intValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*p = intValue(v)
	return nil
}

type int64Value int64

func newInt64Value(p *int64) *int64Value { return (*int64Value)(p) }
func (p *int64Value) String() string     { return strconv.FormatInt(int64(*p), 10) }
func (p *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*p = int64Value(v)
	return nil
}

// uint
type uintValue uint

func newUintValue(p *uint) *uintValue { return (*uintValue)(p) }
func (p *uintValue) String() string   { return strconv.FormatUint(uint64(*p), 10) }
func (p *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*p = uintValue(v)
	return nil
}

type uint64Value uint64

func newUint64Value(p *uint64) *uint64Value { return (*uint64Value)(p) }
func (p *uint64Value) String() string       { return strconv.FormatUint(uint64(*p), 10) }
func (p *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*p = uint64Value(v)
	return nil
}

// float
type float32Value float32

func newFloat32Value(p *float32) *float32Value { return (*float32Value)(p) }
func (p *float32Value) String() string         { return strconv.FormatFloat(float64(*p), 'g', -1, 64) }
func (p *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	*p = float32Value(v)
	return nil
}

type float64Value float64

func newFloat64Value(p *float64) *float64Value { return (*float64Value)(p) }
func (p *float64Value) String() string         { return strconv.FormatFloat(float64(*p), 'g', -1, 64) }
func (p *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*p = float64Value(v)
	return nil
}

// strings
type stringsValue []string

func newStringsValue(p *[]string) *stringsValue { return (*stringsValue)(p) }
func (p *stringsValue) String() string          { return mustMarshalJson(*p) }
func (p *stringsValue) Set(s string) error {
	*p = append(*p, s)
	return nil
}

// ints
type intsValue []int

func newIntsValue(p *[]int) *intsValue { return (*intsValue)(p) }
func (p *intsValue) String() string    { return mustMarshalJson(*p) }
func (p *intsValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*p = append(*p, v)
	return nil
}

type int64sValue []int64

func newInt64sValue(p *[]int64) *int64sValue { return (*int64sValue)(p) }
func (p *int64sValue) String() string        { return mustMarshalJson(*p) }
func (p *int64sValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*p = append(*p, v)
	return nil
}

// uints
type uintsValue []uint

func newUintsValue(p *[]uint) *uintsValue { return (*uintsValue)(p) }
func (p *uintsValue) String() string      { return mustMarshalJson(*p) }
func (p *uintsValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*p = append(*p, uint(v))
	return nil
}

type uint64sValue []uint64

func newUint64sValue(p *[]uint64) *uint64sValue { return (*uint64sValue)(p) }
func (p *uint64sValue) String() string          { return mustMarshalJson(*p) }
func (p *uint64sValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*p = append(*p, v)
	return nil
}

// floats
type float32sValue []float32

func newFloat32sValue(p *[]float32) *float32sValue { return (*float32sValue)(p) }
func (p *float32sValue) String() string            { return mustMarshalJson(*p) }
func (p *float32sValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	*p = append(*p, float32(v))
	return nil
}

type float64sValue []float64

func newFloat64sValue(p *[]float64) *float64sValue { return (*float64sValue)(p) }
func (p *float64sValue) String() string            { return mustMarshalJson(*p) }
func (p *float64sValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*p = append(*p, v)
	return nil
}
