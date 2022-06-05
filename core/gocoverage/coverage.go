/*
 * Copyright (C) distroy
 */

package gocoverage

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/distroy/git-go-tool/core/gocore"
	"github.com/distroy/git-go-tool/core/iocore"
)

type Coverage struct {
	Filename  string
	BeginLine int
	EndLine   int
	Count     int
}

func (c Coverage) String() string {
	return fmt.Sprintf("%s:%d,%d#%d", c.Filename, c.BeginLine, c.EndLine, c.Count)
}

func ParseFile(filePath string) ([]Coverage, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("open coverage file fail. file:%s, err:%v", filePath, err)
		return nil, err
	}
	defer f.Close()

	res, err := ParseReader(f)
	if err != nil {
		log.Fatalf("parse coverage file fail. file:%s, err:%v", filePath, err)
		return nil, err
	}

	return res, err
}

func ParseReader(reader io.Reader) ([]Coverage, error) {
	modPrefix := gocore.GetModPrefix()
	return parseReader(modPrefix, reader)
}

func parseReader(prefix string, reader io.Reader) ([]Coverage, error) {
	r := iocore.NewLineReader(reader)

	res := make([]Coverage, 0, 32)
	for {
		line, err := r.ReadString()
		if err != nil {
			if err == io.EOF {
				return res, nil
			}
			return nil, err
		}

		c, ok := parseLine(prefix, line)
		if !ok {
			return nil, fmt.Errorf("invalid coverage line. line:%s", line)

		} else if c == nil {
			continue
		}

		res = append(res, *c)
	}
}

// format: {prefix}/name.go:line.column,line.column numberOfStatements count
func parseLine(prefix string, line string) (*Coverage, bool) {
	// format: name.go:line.column,line.column numberOfStatements count
	// {prefix}/core/iocore/abc.go:70.28,71.29 1 1
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, prefix) {
		return nil, true
	}

	// remove mod prefix
	line = line[len(prefix):]
	if strings.HasPrefix(line, "/") {
		line = line[1:]
	}

	items := strings.Split(line, ":")
	if len(items) != 2 {
		return nil, false
	}
	file := items[0]

	// format: line.column,line.column numberOfStatements count
	// 70.28,71.29 1 1
	line = items[1]

	// format: line column line column numberOfStatements count
	// 70 28 71 29 1 1
	replacer := strings.NewReplacer(".", " ", ",", " ")
	line = replacer.Replace(line)
	items = strings.Split(line, " ")

	if len(items) != 6 {
		return nil, false
	}

	pos, err := strconv.Atoi(strings.Split(items[0], ".")[0]) // format: line.column
	if err != nil {
		return nil, false
	}

	end, err := strconv.Atoi(strings.Split(items[2], ".")[0]) // format: line.column
	if err != nil {
		return nil, false
	}

	count, err := strconv.Atoi(items[5])
	if err != nil {
		return nil, false
	}

	return &Coverage{
		Filename:  file,
		BeginLine: pos,
		EndLine:   end,
		Count:     count,
	}, true
}
