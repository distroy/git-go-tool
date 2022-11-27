/*
 * Copyright (C) distroy
 */

package git

import (
	"strings"

	"github.com/distroy/git-go-tool/core/execcore"
)

type SubModule struct {
	CommitId string
	Path     string // relative path
}

func MustGetSubModules() []*SubModule {
	res, err := GetSubModules()
	if err != nil {
		panic(err)
	}
	// log.Printf("submodules: %s", jsoncore.MustMarshal(res))
	return res
}

func GetSubModules() ([]*SubModule, error) {
	out, err := execcore.GetOutput("git", "submodule", "status")
	if err != nil {
		return nil, err
	}

	return parseSubModules(out), nil
}

func parseSubModules(s string) []*SubModule {
	// s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	res := make([]*SubModule, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		vv := strings.Split(line, " ")
		if len(vv) < 2 {
			continue
		}

		sub := &SubModule{
			CommitId: strings.TrimLeft(vv[0], "-"),
			Path:     vv[1],
		}
		// log.Printf(" === 111 %#v", line)
		// log.Printf(" === 222 %#v", sub)

		res = append(res, sub)
	}

	return res
}
