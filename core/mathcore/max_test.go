/*
 * Copyright (C) distroy
 */

package mathcore

import (
	"testing"
)

func TestMax(t *testing.T) {
	const (
		arg0 = 3
		arg1 = 4
		want = 4
	)
	if got := MaxInt(arg0, arg1); got != want {
		t.Errorf("MaxInt() = %v, want %v", got, want)
	}
	if got := MaxInt8(arg0, arg1); got != want {
		t.Errorf("MaxInt() = %v, want %v", got, want)
	}
	if got := MaxInt16(arg0, arg1); got != want {
		t.Errorf("MaxInt() = %v, want %v", got, want)
	}
	if got := MaxInt32(arg0, arg1); got != want {
		t.Errorf("MaxInt() = %v, want %v", got, want)
	}
	if got := MaxInt64(arg0, arg1); got != want {
		t.Errorf("MaxInt() = %v, want %v", got, want)
	}
}
