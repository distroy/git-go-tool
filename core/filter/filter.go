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
	Includes Values `flag:"name:include; meta:regexp; usage:the regexp for include pathes"`
	Excludes Values `flag:"name:exclude; meta:regexp; usage:the regexp for exclude pathes"`
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
