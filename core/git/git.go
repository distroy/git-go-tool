/*
 * Copyright (C) distroy
 */

package git

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/distroy/git-go-tool/core/iocore"
)

type Different struct {
	Filename  string
	BeginLine int
	EndLine   int
}

func ParesNewLinesByCommand(branch string) ([]Different, error) {
	cmd := exec.Command("git", "diff", "--unified", branch)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
	}

	defer func() {
		cmd.Wait()
	}()

	res, err := ParesNewLinesByReader(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}

func ParesNewLinesByReader(r io.Reader) ([]Different, error) {

	res := make([]Different, 0, 1024)

	buffer := make([]string, 0, 32)
	reader := iocore.NewLineReader(r)
	for {
		line, err := reader.ReadString()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if !strings.HasPrefix(line, "diff ") {
			log.Fatalf("unexpected line prefix for file begin. line:%s", line)
		}

		buffer = append(buffer, line)
		for {
			line, err = reader.PeekString()
			if strings.HasPrefix(line, "diff ") {
				break
			}
			buffer = append(buffer, line)
		}
	}

	return res, nil
}

func parseFileNewLines(lines []string) ([]Different, error) {
	i, l := 0, len(lines)

	var file string
	for i < l {
		line := lines[i]
		prefix := "+++ "
		if !strings.HasPrefix(line, prefix) {
			i++
			continue
		}

		file = line[len(prefix):]

		if file != "/dev/null" {
			if !strings.HasPrefix(file, "b/") {
				log.Fatalf("invalid git different content. content:\n%s", strings.Join(lines, "\n"))
			}
			file = file[2:]
		}
	}
	// diff --git a/script/complexity/core/__init__.py b/script/git-tool/core/__init__.py
	// similarity index 100%
	// rename from script/complexity/core/__init__.py
	// rename to script/git-tool/core/__init__.py
	// diff --git a/script/complexity/core/exec.py b/script/git-tool/core/exec.py
	// ...
}
