/*
 * Copyright (C) distroy
 */

package goformat

type Level int

const (
	LevelError Level = iota
	LevelWarning
	LevelInfo
)

func (l Level) String() string {
	switch l {
	case LevelError:
		return "error"
	case LevelWarning:
		return "warning"
	case LevelInfo:
		return "info"
	}
	return "error"
}

type Issue struct {
	Filename    string
	BeginLine   int
	EndLine     int
	Level       Level
	Description string
	// Code        []string
}

type Count struct {
	Error   int
	Warning int
	Info    int
}
