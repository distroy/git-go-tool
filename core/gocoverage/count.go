/*
 * Copyright (C) distroy
 */

package gocoverage

type Count struct {
	Coverages    int
	NonCoverages int
}

func (c Count) IsZero() bool {
	return c.Coverages == 0 && c.NonCoverages == 0
}

func (c Count) GetRate() float64 {
	total := c.Coverages + c.NonCoverages
	if total == 0 {
		return 1
	}

	return float64(c.Coverages) / float64(total)
}
