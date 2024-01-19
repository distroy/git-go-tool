/*
 * Copyright (C) distroy
 */

package validate

import (
	"os"

	"github.com/distroy/git-go-tool/obj/resultobj"
	"github.com/distroy/git-go-tool/service/modeservice"
)

func Result(p *resultobj.Result) bool {
	if p.Type == "" || p.Data == nil {
		return false
	}

	if p.Mode == "" {
		p.Mode = modeservice.ModeDefault
	}

	if p.ProjectUrl == "" {
		p.ProjectUrl = os.Getenv("CI_MERGE_REQUEST_PROJECT_URL")
	}
	if p.TargetBranch == "" {
		p.TargetBranch = os.Getenv("CI_MERGE_REQUEST_TARGET_BRANCH_NAME")
	}
	if p.SourceBranch == "" {
		p.SourceBranch = os.Getenv("CI_MERGE_REQUEST_SOURCE_BRANCH_NAME")
	}

	if p.ProjectUrl == "" || p.TargetBranch == "" || p.SourceBranch == "" {
		return false
	}

	return true
}
