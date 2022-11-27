/*
 * Copyright (C) distroy
 */

package refcore

import (
	"testing"
	"unsafe"

	"github.com/distroy/git-go-tool/3rd/convey"
)

func TestIsZero(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("invalid", func() {
			convey.So(IsZero(nil), convey.ShouldEqual, true)
		})

		convey.Convey("bool", func() {
			convey.So(IsZero(false), convey.ShouldEqual, true)
			convey.So(IsZero(true), convey.ShouldEqual, false)
		})

		convey.Convey("int", func() {
			convey.So(IsZero(int(0)), convey.ShouldEqual, true)
			convey.So(IsZero(int(1)), convey.ShouldEqual, false)
		})

		convey.Convey("uint", func() {
			convey.So(IsZero(uint(0)), convey.ShouldEqual, true)
			convey.So(IsZero(uint(1)), convey.ShouldEqual, false)
		})

		convey.Convey("float64", func() {
			convey.So(IsZero(float64(0)), convey.ShouldEqual, true)
			convey.So(IsZero(float64(1)), convey.ShouldEqual, false)
		})

		convey.Convey("*int", func() {
			n := 1
			convey.So(IsZero((*int)(nil)), convey.ShouldEqual, true)
			convey.So(IsZero((*int)(&n)), convey.ShouldEqual, false)
		})

		convey.Convey("[]int", func() {
			convey.So(IsZero(([]int)(nil)), convey.ShouldEqual, true)
			convey.So(IsZero([]int{}), convey.ShouldEqual, false)
		})

		convey.Convey("[...]int", func() {
			convey.So(IsZero([...]int{}), convey.ShouldEqual, true)
			convey.So(IsZero([...]int{0, 0, 1}), convey.ShouldEqual, false)
		})

		convey.Convey("string", func() {
			convey.So(IsZero(""), convey.ShouldEqual, true)
			convey.So(IsZero("1"), convey.ShouldEqual, false)
		})

		convey.Convey("complex128", func() {
			convey.So(IsZero(complex(0, 0)), convey.ShouldEqual, true)
			convey.So(IsZero(complex(0, 1)), convey.ShouldEqual, false)
			convey.So(IsZero(complex(1, 0)), convey.ShouldEqual, false)
		})

		convey.Convey("unsafe.Pointer", func() {
			n := 1
			convey.So(IsZero((unsafe.Pointer)(nil)), convey.ShouldEqual, true)
			convey.So(IsZero((unsafe.Pointer)(&n)), convey.ShouldEqual, false)
		})

		convey.Convey("map[int]struct{}", func() {
			convey.So(IsZero((map[int]struct{})(nil)), convey.ShouldEqual, true)
			convey.So(IsZero(map[int]struct{}{}), convey.ShouldEqual, false)
		})

		convey.Convey("func ()", func() {
			convey.So(IsZero((func())(nil)), convey.ShouldEqual, true)
			convey.So(IsZero(func() {}), convey.ShouldEqual, false)
		})

		convey.Convey("chan struct{}", func() {
			convey.So(IsZero((chan struct{})(nil)), convey.ShouldEqual, true)
			convey.So(IsZero(make(chan struct{})), convey.ShouldEqual, false)
		})

		convey.Convey("struct", func() {
			type testStruct = struct {
				Int int
			}
			convey.So(IsZero(testStruct{}), convey.ShouldEqual, true)
			convey.So(IsZero(testStruct{Int: 1}), convey.ShouldEqual, false)
		})
	})
}
