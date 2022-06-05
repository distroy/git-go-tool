/*
 * Copyright (C) distroy
 */

package gocognitive

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"testing"
)

// doc: https://sonarsource.com/docs/CognitiveComplexity.pdf

func getFileLine() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	return file, line
}

func complexityFromFuncName(funcName string) int {
	var pos int
	for i := len(funcName) - 1; i >= 0; i-- {
		b := funcName[i]
		if b >= '0' && b <= '9' {
			continue
		}

		pos = i + 1
		break
	}
	str := funcName[pos:]
	n, _ := strconv.Atoi(str)
	return n
}

func TestExample(t *testing.T) {
	file, _ := getFileLine()
	file = fmt.Sprintf("%s/%s", path.Dir(file), "example_for_test.go")
	t.Logf("example file: %s", file)

	res, err := AnalyzeFileByPath(file)
	if err != nil {
		t.Fatalf("analyze file fail. err:%s", err)
	}
	for _, v := range res {
		t.Run(v.FuncName, func(t *testing.T) {
			want := complexityFromFuncName(v.FuncName)
			if v.Complexity != want {
				t.Errorf("check func complexity fail. got:%d, want:%d", v.Complexity, want)
			}
		})
	}
}
