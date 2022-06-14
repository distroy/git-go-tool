/*
 * Copyright (C) distroy
 */

package mathcore

import "testing"

func TestMin(t *testing.T) {
	const (
		arg0 = 4
		arg1 = 3
		want = 3
	)
	if got := MinInt(arg0, arg1); got != want {
		t.Errorf("MinInt() = %v, want %v", got, want)
	}
	if got := MinInt8(arg0, arg1); got != want {
		t.Errorf("MinInt() = %v, want %v", got, want)
	}
	if got := MinInt16(arg0, arg1); got != want {
		t.Errorf("MinInt() = %v, want %v", got, want)
	}
	if got := MinInt32(arg0, arg1); got != want {
		t.Errorf("MinInt() = %v, want %v", got, want)
	}
	if got := MinInt64(arg0, arg1); got != want {
		t.Errorf("MinInt() = %v, want %v", got, want)
	}
}
