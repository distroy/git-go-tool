/*
 * Copyright (C) distroy
 */

package modeservice

import "strings"

type modeBase struct{}

func (m *modeBase) isLineIgnored(line string) bool {
	line = strings.TrimSpace(line)
	return len(line) == 0
}
