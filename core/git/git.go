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
	return parseNewLinesByCommand(cmd)
}

func parseNewLinesByCommand(cmd *exec.Cmd) ([]Different, error) {
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("start pipe for exec fail. cmd:%s, err:%v", getCommandString(cmd), err)
	}

	defer func() {
		cmd.Wait()
	}()

	res, err := ParesNewLinesByReader(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		switch v := err.(type) {
		default:
			log.Fatalf("exec command fail. cmd:%s, err:%v", getCommandString(cmd), v.Error())

		case *exec.ExitError:
			log.Fatalf("exec command fail. cmd:%s, code:%d, err:%v",
				getCommandString(cmd), v.ExitCode(), v.Error())
		}
		return nil, err
	}

	return res, nil
}

func ParesNewLinesByReader(r io.Reader) ([]Different, error) {
	res := make([]Different, 0, 1024)

	reader := iocore.NewLineReader(r)
	for {
		fileLines, err := readFileLines(reader)
		if err != nil {
			if err == io.EOF {
				return res, io.EOF
			}
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

func getCommandString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := &strings.Builder{}
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
