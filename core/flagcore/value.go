/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"encoding/json"
	"strconv"
)

func mustMarshalJson(v interface{}) string {
	d, _ := json.Marshal(v)
	return string(d)
}

// bool
type boolValue bool

func newBoolValue(p *bool) *boolValue {
	return (*boolValue)(p)
}
func (p *boolValue) String() string { return strconv.FormatBool(bool(*p)) }
func (p *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	// log.Printf(" === %v, %v, %v", *p, v, err)
	if err != nil {
		return err
	}
	*p = boolValue(v)
	return nil
}

// func (p *boolValue) Get() interface{} { return bool(*p) }
// func (p *boolValue) IsBoolFlag() bool { return true }

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

// strings
type stringsValue []string

func newStringsValue(p *[]string) *stringsValue { return (*stringsValue)(p) }
func (p *stringsValue) String() string          { return mustMarshalJson(*p) }
func (p *stringsValue) Set(s string) error {
	*p = append(*p, s)
	return nil
}
