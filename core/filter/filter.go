/*
 * Copyright (C) distroy
 */

package filter

type Values interface {
	Set(s string) error
	String() string

	Check(s string) bool
}

// type Filter interface {
// 	Includes() Values
// 	Excludes() Values
//
// 	Check(s string) bool
// }

type Filter struct {
	Includes Values
	Excludes Values
}

func (f *Filter) Check(s string) bool {
	if f.Includes.Check(s) {
		return true
	}
	if f.Excludes.Check(s) {
		return false
	}
	return true
}
