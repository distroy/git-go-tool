/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"strconv"
	"time"

	"github.com/distroy/git-go-tool/core/ptrcore"
)

// duration
type durationPtrValue struct{ V **time.Duration }

func newDurationPtrValue(p **time.Duration) durationPtrValue { return durationPtrValue{V: p} }
func (p durationPtrValue) String() string {
	if *p.V == nil {
		return ""
	}
	return time.Duration(**p.V).String()
}
func (p durationPtrValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

// string
type stringPtrValue struct{ V **string }

func newStringPtrValue(p **string) stringPtrValue { return stringPtrValue{V: p} }
func (p stringPtrValue) String() string           { return ptrcore.GetString(*p.V) }
func (p stringPtrValue) Set(s string) error {
	*p.V = &s
	return nil
}

// bool
type boolPtrValue struct{ V **bool }

func newBoolPtrValue(p **bool) boolPtrValue { return boolPtrValue{V: p} }
func (p boolPtrValue) String() string       { return strconv.FormatBool(ptrcore.GetBool(*p.V)) }
func (p boolPtrValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	// log.Printf(" === %v, %v, %v", *p, v, err)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

type boolPtrFlag boolPtrValue

func newBoolPtrFlag(p boolPtrValue) boolPtrFlag { return boolPtrFlag(p) }
func (p boolPtrFlag) String() string            { return (boolPtrValue)(p).String() }
func (p boolPtrFlag) Set(s string) error        { return (boolPtrValue)(p).Set(s) }
func (p boolPtrFlag) IsBoolFlag() bool          { return true }

// func (p *boolPtrValue) Get() interface{} { return bool(*p) }
// func (p *boolPtrValue) IsBoolFlag() bool { return true }

// int
type intPtrValue struct{ V **int }

func newIntPtrValue(p **int) intPtrValue { return intPtrValue{V: p} }
func (p intPtrValue) String() string     { return strconv.Itoa(ptrcore.GetInt(*p.V)) }
func (p intPtrValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

type int64PtrValue struct{ V **int64 }

func newInt64PtrValue(p **int64) int64PtrValue { return int64PtrValue{V: p} }
func (p int64PtrValue) String() string         { return strconv.FormatInt(ptrcore.GetInt64(*p.V), 10) }
func (p int64PtrValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

// uint
type uintPtrValue struct{ V **uint }

func newUintPtrValue(p **uint) uintPtrValue { return uintPtrValue{V: p} }
func (p uintPtrValue) String() string       { return strconv.FormatUint(uint64(ptrcore.GetUint(*p.V)), 10) }
func (p uintPtrValue) Set(s string) error {
	u64, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	v := uint(u64)
	*p.V = &v
	return nil
}

type uint64PtrValue struct{ V **uint64 }

func newUint64PtrValue(p **uint64) uint64PtrValue { return uint64PtrValue{V: p} }
func (p uint64PtrValue) String() string           { return strconv.FormatUint(ptrcore.GetUint64(*p.V), 10) }
func (p uint64PtrValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

// float
type float32PtrValue struct{ V **float32 }

func newFloat32PtrValue(p **float32) float32PtrValue { return float32PtrValue{V: p} }
func (p float32PtrValue) String() string {
	return strconv.FormatFloat(float64(ptrcore.GetFloat32(*p.V)), 'g', -1, 64)
}
func (p float32PtrValue) Set(s string) error {
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	v := float32(f64)
	*p.V = &v
	return nil
}

type float64PtrValue struct{ V **float64 }

func newFloat64PtrValue(p **float64) float64PtrValue { return float64PtrValue{V: p} }
func (p float64PtrValue) String() string {
	return strconv.FormatFloat(ptrcore.GetFloat64(*p.V), 'g', -1, 64)
}
func (p float64PtrValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}
