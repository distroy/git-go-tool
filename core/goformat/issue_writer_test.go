/*
 * Copyright (C) distroy
 */

package goformat

import (
	"strings"
	"testing"
)

func Test_writerWrapper_WriteIssues(t *testing.T) {
	tests := []struct {
		name   string
		issues []*Issue
		count  Count
		text   string
	}{
		{
			issues: []*Issue{
				{
					Filename:    "test.go",
					Level:       LevelInfo,
					Description: "issue 1",
				},
				{
					Filename:    "test.go",
					BeginLine:   1,
					EndLine:     1,
					Level:       LevelWarning,
					Description: "issue 2",
				},
				{
					Filename:    "test.go",
					BeginLine:   1,
					EndLine:     4,
					Level:       LevelError,
					Description: "issue 3",
				},
			},
			count: Count{
				Error:   1,
				Warning: 1,
				Info:    1,
			},
			text: strings.Join([]string{
				"test.go [info] issue 1",
				"test.go:1 [warning] issue 2",
				"test.go:1,4 [error] issue 3",
				"",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &strings.Builder{}
			w := NewIssueWriter(buf)
			w.WriteIssues(tt.issues)

			if got := w.Count(); got != tt.count {
				t.Errorf("WriteIssues() count = %#v, want:%#v", got, tt.count)
			}

			if got := buf.String(); got != tt.text {
				t.Errorf("WriteIssues() text = \n%s\n    want:\n%s", got, tt.text)
				return
			}
		})
	}
}
