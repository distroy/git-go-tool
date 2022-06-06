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

func ParseNewLines(branch string) ([]Different, error) {
	cmd := exec.Command("git", "diff", "--unified=0", branch)
	return parseNewLinesFromCommand(cmd)
}

func parseNewLinesFromCommand(cmd *exec.Cmd) ([]Different, error) {
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("start pipe for exec fail. cmd:%s, err:%v", getCommandString(cmd), err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), err)
	}

	defer func() {
		cmd.Wait()
	}()

	res, err := ParseNewLinesFromReader(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		switch v := err.(type) {
		default:
			log.Fatalf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), v)

		case *exec.ExitError:
			log.Fatalf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v)
		}
		return nil, err
	}

	return res, nil
}

func ParseNewLinesFromReader(r io.Reader) ([]Different, error) {
	res := make([]Different, 0, 1024)

	reader := iocore.NewLineReader(r)
	for {
		fileLines, err := readFileLines(reader)
		if err != nil {
			if err == io.EOF {
				return res, nil
			}
			// log.Fatalf("parse new lines from reader fail. err:%v", err)
			return nil, err
		}

		news, err := parseNewLinesFromFileLines(fileLines)
		if err != nil {
			return nil, err
		}

		res = append(res, news...)
	}
}

func readFileLines(r iocore.LineReader) ([]string, error) {
	line, err := r.ReadString()
	if err != nil {
		// if err == io.EOF {
		// 	return nil, io.EOF
		// }
		return nil, err
	}

	if !strings.HasPrefix(line, "diff ") {
		log.Fatalf("unexpected line prefix for file begin. line:%s", line)
	}

	buffer := make([]string, 0, 32)
	buffer = append(buffer, line)
	for {
		line, err = r.PeekString()
		if err == io.EOF {
			return buffer, nil

		} else if err != nil {
			return nil, err

		} else if strings.HasPrefix(line, "diff ") {
			return buffer, nil
		}

		r.Read()
		buffer = append(buffer, line)
	}
}
