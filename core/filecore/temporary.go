/*
 * Copyright (C) distroy
 */

package filecore

import (
	"fmt"
	"math/rand"
	"os"
)

func MustTempFile() string {
	name, err := TempFile()
	if err != nil || len(name) == 0 {
		panic(fmt.Sprintf("create temp file fail. err:%v", err))
	}
	return name
}

func TempFile() (string, error) {
	retry := 10000
	for i := 0; i < retry; i++ {
		name := fmt.Sprintf("tmp_%d", rand.Int())
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if err == nil {
			f.Close()
			return name, nil
		}
	}

	return "", fmt.Errorf("create temp file fail")
}
