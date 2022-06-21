/*
 * Copyright (C) distroy
 */

package goformat

type Level string

const (
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

type Issue struct {
	Filename  string
	BeginLine int
	EndLine   int
	Level     Level
	Message   string
}
