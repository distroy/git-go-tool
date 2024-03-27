/*
 * Copyright (C) distroy
 */

package validate

import (
	"os"
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/obj/resultobj"
)

func TestResult(t *testing.T) {
	type args struct {
		obj *resultobj.Result
		env map[string]string
	}
	type want struct {
		ok  bool
		obj *resultobj.Result
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil input",
			args: args{obj: nil},
			want: want{obj: nil},
		},
		{
			name: "nil data",
			args: args{
				obj: &resultobj.Result{
					Mode:           "",
					Type:           "go-format",
					PipelineId:     "1001",
					PipelineUrl:    "1002",
					JobId:          "2001",
					JobUrl:         "2002",
					ProjectId:      "3001",
					MergeRequestId: "4001",
					ProjectUrl:     "5001",
					TargetBranch:   "5002",
					SourceBranch:   "",
					Data:           nil,
				},
				env: map[string]string{
					"CI_PIPELINE_ID":                      "x-1001",
					"CI_MERGE_REQUEST_SOURCE_BRANCH_NAME": "x-5003",
				},
			},
			want: want{
				ok: false,
				obj: &resultobj.Result{
					Mode:           "diff",
					Type:           "go-format",
					PipelineId:     "1001",
					PipelineUrl:    "1002",
					JobId:          "2001",
					JobUrl:         "2002",
					ProjectId:      "3001",
					MergeRequestId: "4001",
					ProjectUrl:     "5001",
					TargetBranch:   "5002",
					SourceBranch:   "x-5003",
					Data:           nil,
				},
			},
		},
		{
			name: "zero str",
			args: args{
				obj: &resultobj.Result{
					Mode:           "",
					Type:           "go-format",
					PipelineId:     "",
					PipelineUrl:    "1002",
					JobId:          "2001",
					JobUrl:         "",
					ProjectId:      "3001",
					MergeRequestId: "4001",
					ProjectUrl:     "5001",
					TargetBranch:   "5002",
					SourceBranch:   "5003",
					Data:           struct{}{},
				},
				env: map[string]string{
					"CI_PIPELINE_ID":                      "x-1001",
					"CI_MERGE_REQUEST_SOURCE_BRANCH_NAME": "x-5003",
				},
			},
			want: want{
				ok: false,
				obj: &resultobj.Result{
					Mode:           "diff",
					Type:           "go-format",
					PipelineId:     "x-1001",
					PipelineUrl:    "1002",
					JobId:          "2001",
					JobUrl:         "",
					ProjectId:      "3001",
					MergeRequestId: "4001",
					ProjectUrl:     "5001",
					TargetBranch:   "5002",
					SourceBranch:   "5003",
					Data:           struct{}{},
				},
			},
		},
		{
			name: "succ",
			args: args{
				obj: &resultobj.Result{
					Mode:           "",
					Type:           "go-format",
					PipelineId:     "",
					PipelineUrl:    "1002",
					JobId:          "2001",
					JobUrl:         "",
					ProjectId:      "3001",
					MergeRequestId: "4001",
					ProjectUrl:     "5001",
					TargetBranch:   "5002",
					SourceBranch:   "5003",
					Data:           struct{}{},
				},
				env: map[string]string{
					"CI_PIPELINE_ID":                      "x-1001",
					"CI_JOB_URL":                          "x-2002",
					"CI_MERGE_REQUEST_TITLE":              "x-4002",
					"CI_MERGE_REQUEST_SOURCE_BRANCH_NAME": "x-5003",
				},
			},
			want: want{
				ok: true,
				obj: &resultobj.Result{
					Mode:              "diff",
					Type:              "go-format",
					PipelineId:        "x-1001",
					PipelineUrl:       "1002",
					JobId:             "2001",
					JobUrl:            "x-2002",
					ProjectId:         "3001",
					MergeRequestId:    "4001",
					MergeRequestTitle: "x-4002",
					ProjectUrl:        "5001",
					TargetBranch:      "5002",
					SourceBranch:      "5003",
					Data:              struct{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.env {
				os.Setenv(k, v)
			}

			if got := Result(tt.args.obj); got != tt.want.ok {
				t.Errorf("Result() = %v, want %v", got, tt.want.ok)
			}
			if !reflect.DeepEqual(tt.args.obj, tt.want.obj) {
				t.Errorf("output res = %v, want %v", tt.args.obj, tt.want.obj)
			}

			for k := range tt.args.env {
				os.Setenv(k, "")
			}
		})
	}
}
