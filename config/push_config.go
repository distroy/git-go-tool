/*
 * Copyright (C) distroy
 */

package config

var DefaultPush = &PushConfig{}

type PushConfig struct {
	PushUrl      string `yaml:"push_url"`
	ProjectUrl   string `yaml:"project_url"` // env: CI_MERGE_REQUEST_PROJECT_URL
	TargetBranch string `yaml:"-"`           // env: CI_MERGE_REQUEST_TARGET_BRANCH_NAME
	SourceBranch string `yaml:"-"`           // env: CI_MERGE_REQUEST_SOURCE_BRANCH_NAME
}
