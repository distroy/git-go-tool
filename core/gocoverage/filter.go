/*
 * Copyright (C) distroy
 */

package gocoverage

type filter = func(file string, lineNo int) bool

func doFilters(file string, lineNo int, filters []filter) bool {
	for _, filter := range filters {
		if !filter(file, lineNo) {
			return false
		}
	}
	return true
}
