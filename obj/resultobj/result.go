/*
 * Copyright (C) distroy
 */

package resultobj

import (
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/goformat"
)

const (
	TypeGoCognitive = "go-cognitive"
	TypeGoCoverage  = "go-coverage"
	TypeGoFormat    = "go-format"
)

type Result struct {
	Mode           string      `json:"mode"`
	Type           string      `json:"type"`
	PipelineId     string      `json:"pipeline_id"`      // env: CI_PIPELINE_ID
	PipelineUrl    string      `json:"pipeline_url"`     // env: CI_PIPELINE_URL
	JobId          string      `json:"job_id"`           // env: CI_JOB_ID
	JobUrl         string      `json:"job_url"`          // env: CI_JOB_URL
	ProjectId      string      `json:"project_id"`       // env: CI_MERGE_REQUEST_PROJECT_ID
	MergeRequestId string      `json:"merge_request_id"` // env: CI_MERGE_REQUEST_IID
	ProjectUrl     string      `json:"project_url"`      // env: CI_MERGE_REQUEST_PROJECT_URL
	TargetBranch   string      `json:"target_branch"`    // env: CI_MERGE_REQUEST_TARGET_BRANCH_NAME
	SourceBranch   string      `json:"source_branch"`    // env: CI_MERGE_REQUEST_SOURCE_BRANCH_NAME
	Data           interface{} `json:"data"`
}

type GoBaseData struct {
	ExecError        string   `json:"exec_error,omitempty"`
	ExecErrorDetails []string `json:"exec_error_details,omitempty"`
}

type GoComplexityData struct {
	GoBaseData `json:",inline"`

	Threshold              int                       `json:"threshold"`
	FunctionsOverThreshold []*gocognitive.Complexity `json:"functions_over_threshold"`
}

type GoFormatData struct {
	GoBaseData `json:",inline"`

	IssueCount goformat.Count    `json:"issue_count"`
	Issues     []*goformat.Issue `json:"issues"`
}

type GoCoverageData struct {
	GoBaseData `json:",inline"`

	Threshold            float64 `json:"threshold"`
	Rate                 float64 `json:"rate"`
	CoverageLineCount    int     `json:"coverage_line_count"`
	NonCoverageLineCount int     `json:"non_coverage_line_count"`
}
