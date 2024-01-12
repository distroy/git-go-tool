/*
 * Copyright (C) distroy
 */

package goformat

type Level string

const (
	LevelError   Level = "error"
	LevelWarning Level = "warning"
	LevelInfo    Level = "info"
)

func (l Level) String() string {
	if l == "" {
		return "error"
	}
	return string(l)
}

type Issue struct {
	Filename    string `json:"file_name"`
	BeginLine   int    `json:"begin_line"`
	EndLine     int    `json:"end_line"`
	Level       Level  `json:"level"`
	Description string `json:"description"`
	// Code        []string
}

type Count struct {
	Error   int `json:"error"`
	Warning int `json:"warning"`
	Info    int `json:"info"`
}
