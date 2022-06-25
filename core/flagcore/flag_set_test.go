/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"testing"
)

func TestNewFlagSet(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFlagSet(); tt.want != (got != nil) {
				t.Errorf("NewFlagSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
