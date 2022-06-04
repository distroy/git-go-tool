/*
 * Copyright (C) distroy
 */

package gocoverage

type filter = func(file string, begin, end int) bool

func doFilters(file string, begin, end int, filters []filter) bool {
	for _, filter := range filters {
		if !filter(file, begin, end) {
			return false
		}
	}
	return true
}
