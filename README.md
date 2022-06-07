[![GoDoc](https://godoc.org/github.com/distroy/git-go-tool?status.svg)](https://godoc.org/github.com/distroy/git-go-tool)

# GitLab CI/CD 接入

```yml
include:
  - project: 'shopee/marketing/git-go-tool'
    ref: 'master'
    file: '/gitlab-ci/main.yml'

stages: # 特别重要
  - go-check-stage # 如果已配置stages，请在stages中添加 go-check-stage，保证 go-check-stage 能够识别到

# 环境变量
variables:
  # # git 的比较模式
  # #   diff: 增量模式，只检查新增代码的认知复杂度和单测覆盖率
  # #   all: 全量模式，检查所有代码的认知复杂度和单测覆盖率
  # GIT_COMPARE_MODE: "diff"
  # # 检查 go 代码的认知复杂度和单测覆盖率时，排除的文件选项
  # # 格式: -exclude <regexp 1> ... -exclude <regexp N>
  # # 例如: -exclude '\.microkit\.go$' -exclude '^tool/'
  # GO_EXCLUDE_OPTIONS: ""

  # # 会检查这个目录下的子目录（只向下检查一级）的go代码是否可以编译通过
  # # 设置成空字符，不检查
  # GO_BUILD_TARGET_DIR: "$CI_PROJECT_DIR/app"
  # # go 编译和单测的选项
  # GO_FLAGS: "-mod=vendor"
  # # go 单测的选项
  # GO_TEST_FLAGS: "-gcflags=all=-l"

  # # 检查单测覆盖率的目标分支名
  # # 默认: MR 的目标分支
  # # 也可以选择固定 release 分支，比如在提测之后修改Bug，单测可能覆盖不到改动的代码
  # GO_COVERAGE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 单测覆盖率的阈值，取值范围 0.0 ~ 1.0
  # GO_COVERAGE_RATE: 0.65
  # # 检查 go 代码的单测覆盖率时，排除的文件选项
  # # 请参照 GO_EXCLUDE_OPTIONS 说明
  # # 检查单测覆盖率时，优先选择GO_COVERAGE_EXCLUDES
  # # 如果 GO_COVERAGE_EXCLUDES 为空字符，则使用GO_EXCLUDE_OPTIONS
  # GO_COVERAGE_EXCLUDES: ""

  # # 检查单测覆盖率的目标分支名
  # # 默认: MR 的目标分支
  # GO_COGNITIVE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 函数认知复杂度的阈值
  # GO_COGNITIVE_OVER: 15
  # # 检查 go 代码的认知复杂度时，排除的文件选项
  # # 请参照 GO_EXCLUDE_OPTIONS 说明
  # # 检查单测覆盖率时，优先选择GO_COGNITIVE_EXCLUDES
  # # 如果 GO_COGNITIVE_EXCLUDES 为空字符，则使用GO_EXCLUDE_OPTIONS
  # GO_COGNITIVE_EXCLUDES: ""
```


# Commands 介绍

## Installation

```shell
go install github.com/distroy/git-go-tool/cmd/go-cognitive@latest
go install github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive@latest
go install github.com/distroy/git-go-tool/cmd/git-diff-go-coverage@latest
```

## go-cognitive
go-cognitive analyze cognitive complexities of functions in Go source code. A measurement of how hard does the code is intuitively to understand.

> The document of cognitive complexity: [https://sonarsource.com/docs/CognitiveComplexity.pdf](https://sonarsource.com/docs/CognitiveComplexity.pdf)
> Example: these example are specific for Go, please see the [core/gocognitive/example_for_test.go](https://github.com/distroy/git-go-tool/blob/master/core/gocognitive/example_for_test.go)
