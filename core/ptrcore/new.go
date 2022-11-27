/*
 * Copyright (C) distroy
 */

package ptrcore

import "time"

func NewBool(d bool) *bool { return &d }

func NewByte(d byte) *byte { return &d }
func NewRune(d rune) *rune { return &d }

func NewInt(d int) *int       { return &d }
func NewInt8(d int8) *int8    { return &d }
func NewInt16(d int16) *int16 { return &d }
func NewInt32(d int32) *int32 { return &d }
func NewInt64(d int64) *int64 { return &d }

func NewUint(d uint) *uint          { return &d }
func NewUint8(d uint8) *uint8       { return &d }
func NewUint16(d uint16) *uint16    { return &d }
func NewUint32(d uint32) *uint32    { return &d }
func NewUint64(d uint64) *uint64    { return &d }
func NewUintptr(d uintptr) *uintptr { return &d }

func NewFloat32(d float32) *float32 { return &d }
func NewFloat64(d float64) *float64 { return &d }

func NewString(d string) *string { return &d }

func NewComplex64(d complex64) *complex64    { return &d }
func NewComplex128(d complex128) *complex128 { return &d }

func NewTime(d time.Time) *time.Time             { return &d }
func NewDuration(d time.Duration) *time.Duration { return &d }
