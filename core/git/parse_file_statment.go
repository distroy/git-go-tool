/*
 * Copyright (C) distroy
 */

package git

import (
	"log"
	"strconv"
	"strings"
)

// +++ b/core/git/git.go
func parseFilenameFromNewFileLine(line string) (string, error) {
	prefix := "+++ "
	if !strings.HasPrefix(line, prefix) {
		log.Fatalf("parse filename fail. line:%s", line)
	}

	file := line[len(prefix):]

	if file != "/dev/null" {
		if !strings.HasPrefix(file, "b/") {
			log.Fatalf("parse filename fail. line:%s", line)
		}
		file = file[2:]
	}

	return file, nil
}

// @@ -0,0 +1,32 @@
// @@ -52 +52 @@
func parsePositionFromSummaryLine(summary string) (begin int, end int, err error) {
	items := strings.Split(summary, " ")
	pos := strings.Split(items[2], ",")

	begin, err = strconv.Atoi(pos[0])
	if err != nil {
		log.Fatalf("pasre line fail. line:%s, err:%v", summary, err.Error())
	}

	end = begin
	if len(pos) > 1 {
		n, err := strconv.Atoi(pos[1])
		if err != nil {
			log.Fatalf("pasre line fail. line:%s, err:%v", summary, err.Error())
		}
		end = begin + n - 1
	}

	return begin, end, nil
}

func indexLines(lines []string, f func(line string) bool) int {
	for i, line := range lines {
		if f(line) {
			return i
		}
	}
	return len(lines)
}

func parseNewLinesFromFileLines(lines []string) ([]Different, error) {
	i, l := 0, len(lines)

	var file string

	fileLineIdx := indexLines(lines, func(s string) bool { return strings.HasPrefix(s, "+++ ") })
	if fileLineIdx < l {
		tmp, err := parseFilenameFromNewFileLine(lines[fileLineIdx])
		if err != nil {
			return nil, err
		}

		file = tmp
	}

	// diff --git a/script/complexity/core/__init__.py b/script/git-tool/core/__init__.py
	// similarity index 100%
	// rename from script/complexity/core/__init__.py
	// rename to script/git-tool/core/__init__.py diff --git a/script/complexity/core/exec.py b/script/git-tool/core/exec.py
	// ...
	if len(file) == 0 {
		return []Different{}, nil
	}

	// skip the header lines
	i = indexLines(lines, func(line string) bool { return strings.HasPrefix(line, "@@ ") })

	res := make([]Different, 0, 32)
	for i < l {
		j := indexLines(lines[i:], func(line string) bool { return strings.HasPrefix(line, "@@ ") })

		statment := lines[i:j]
		news, err := parseNewLinesFromStatmentLines(file, statment)
		if err != nil {
			return nil, err
		}

		i = j
		res = append(res, news...)
	}
	return res, nil
}

func parseNewLinesFromStatmentLines(filename string, lines []string) ([]Different, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	line := lines[0]

	begin, end, err := parsePositionFromSummaryLine(line)
	if err != nil {
		return nil, err
	}

	diff := Different{
		Filename:  filename,
		BeginLine: begin,
		EndLine:   end,
	}

	blankLineNos := parseBlankLineNosFromStatmentLines(lines, diff)

	return removeLineNoFromDifferent(diff, blankLineNos), nil
}

func parseBlankLineNosFromStatmentLines(lines []string, diff Different) []int {
	begin := diff.BeginLine

	blankLineNos := make([]int, 0, 32)
	i := 0
	for _, line := range lines {
		if !strings.HasPrefix(line, "+") {
			continue
		}
		i++

		line = line[1:]
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			blankLineNos = append(blankLineNos, begin+i-1)
		}
	}

	return blankLineNos
}

func removeLineNoFromDifferent(diff Different, lineNos []int) []Different {
	if len(lineNos) == 0 {
		return []Different{diff}
	}

	file := diff.Filename
	pos := diff.BeginLine
	end := diff.EndLine

	res := make([]Different, 0, len(lineNos))

	lastIdx := pos
	for _, idx := range lineNos {
		if lastIdx == idx {
			lastIdx = idx + 1
			continue
		}

		res = append(res, Different{
			Filename:  file,
			BeginLine: lastIdx,
			EndLine:   idx - 1,
		})
		lastIdx = idx + 1
	}

	if lastIdx <= end {
		res = append(res, Different{
			Filename:  file,
			BeginLine: lastIdx,
			EndLine:   end,
		})
	}
	return res
}
