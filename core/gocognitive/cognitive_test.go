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

	res, err := AnalyzeFileByPath(file)
	if err != nil {
		t.Fatalf("analyze file fail. err:%s", err)
	}
	for _, stat := range res {
		want := complexityFromFuncName(stat.FuncName)
		if stat.Complexity == want {
			t.Logf("check func complexity succ. func:%s, complexity:%d, file:%s",
				stat.FuncName, stat.Complexity, stat.Filename)
		} else {
			t.Errorf("check func complexity fail. func:%s, complexity:%d, want:%d, file:%s",
				stat.FuncName, stat.Complexity, want, stat.Filename)
		}
	}
}
