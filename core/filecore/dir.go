/*
 * Copyright (C) distroy
 */

package filecore

import "os"

func IsDir(filename string) bool {
	stat, err := os.Stat(filename)
	return err == nil && stat.IsDir()
}
