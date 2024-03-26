/*
 * Copyright (C) distroy
 */

package validate

import (
	"os"
	"reflect"

	"github.com/distroy/git-go-tool/core/refcore"
	"github.com/distroy/git-go-tool/obj/resultobj"
	"github.com/distroy/git-go-tool/service/modeservice"
)

func checkZeraAndAssignStr(p *string, def string) {
	if *p == "" {
		*p = def
	}
}

func Result(p *resultobj.Result) bool {
	if p == nil {
		return false
	}

	checkZeraAndAssignStr(&p.Mode, modeservice.ModeDefault)

	checkZeraAndAssignStr(&p.PipelineId, os.Getenv("CI_PIPELINE_ID"))
	checkZeraAndAssignStr(&p.PipelineUrl, os.Getenv("CI_PIPELINE_URL"))

	checkZeraAndAssignStr(&p.JobId, os.Getenv("CI_JOB_ID"))
	checkZeraAndAssignStr(&p.JobUrl, os.Getenv("CI_JOB_URL"))

	checkZeraAndAssignStr(&p.ProjectId, os.Getenv("CI_MERGE_REQUEST_PROJECT_ID"))
	checkZeraAndAssignStr(&p.MergeRequestId, os.Getenv("CI_MERGE_REQUEST_IID"))

	checkZeraAndAssignStr(&p.ProjectUrl, os.Getenv("CI_MERGE_REQUEST_PROJECT_URL"))
	checkZeraAndAssignStr(&p.TargetBranch, os.Getenv("CI_MERGE_REQUEST_TARGET_BRANCH_NAME"))
	checkZeraAndAssignStr(&p.SourceBranch, os.Getenv("CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"))

	v := reflect.ValueOf(p).Elem()
	for i, n := 0, v.NumField(); i < n; i++ {
		f := v.Field(i)
		if refcore.IsValZero(f) {
			return false
		}
	}

	return true
}
