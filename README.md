[![GoDoc](https://godoc.org/github.com/distroy/git-go-tool?status.svg)](https://godoc.org/github.com/distroy/git-go-tool)

# 接入

## GitLab CI/CD 接入

- 参考文档：[`.gitlab-ci.yml` file | GitLab](https://docs.gitlab.com/ee/ci/yaml/gitlab_ci_yaml.html)
- 中文文档：[`.gitlab-ci.yml` 文件 | GitLab](https://docs.gitlab.cn/jh/ci/yaml/gitlab_ci_yaml.html)

请在 .gitlab-ci.yml 添加以下配置（请根据需要配置不同的环境变量）

```yml
include:
  - project: 'shopee/marketing/git-go-tool'
    ref: 'master'
    file: '/gitlab-ci/main.yml'

stages: # 特别重要
  - go-check-stage # 如果已配置stages，请在stages中添加 go-check-stage，保证 go-check-stage 能够识别到

# 环境变量
variables:
  # # 用来运行的镜像
  # GO_CHECK_IMAGE: 'ubuntu-latest'

  # *** git flags begin ***
  # # git 的比较模式
  # #   diff: 增量模式，只检查新增代码的认知复杂度和单测覆盖率
  # #   all: 全量模式，检查所有代码的认知复杂度和单测覆盖率
  # GIT_COMPARE_MODE: "diff"
  # *** git flags end ***

  # *** go flags begin ***
  # # go 编译和单测的选项
  # GO_FLAGS: "-mod=vendor"
  # # go 单测的选项
  # GO_TEST_FLAGS: "-gcflags=all=-l"
  # # 检查 go 代码的认知复杂度和单测覆盖率时，排除的文件选项
  # # 格式: -exclude <regexp 1> ... -exclude <regexp N>
  # # 例如: -exclude '\.microkit\.go$' -exclude '^tool/'
  # GO_EXCLUDE_OPTIONS: -exclude '^vendor/'
  # *** go flags end ***

  # *** go buile variables begin ***
  # # 会检查这个目录下的子目录（只向下检查一级）的go代码是否可以编译通过
  # # 设置成空字符，不检查
  # GO_BUILD_TARGET_DIR: "$CI_PROJECT_DIR/app"
  # *** go buile variables end ***

  # *** go format variables begin ***
  # # 检查 go format 度的目标分支名
  # # 默认: MR 的目标分支
  # GO_FORMAT_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 检查 go format 的源分支名
  # GO_FORMAT_DIFF_SOURCE: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
  # # go 代码文件函数限制
  # GO_FORMAT_FILE_LINE_LIMIT: "1000"
  # # 是否检查 import 格式
  # GO_FORMAT_CHECK_IMPORT: "true"
  # # 是否检查源代码被 go 提供的工具格式化了
  # GO_FORMAT_CHECK_SOURCE_FORMATED: "true"
  # # 是否检查 package 命名
  # GO_FORMAT_CHECK_PACKAGE_NAME: "true"
  # 函数参数的数量限制，0-不检查
  # GO_FORMAT_FUNC_INPUT_NUM_LIMIT: "3"
  # # 函数返回值的数量限制，0-不检查
  # GO_FORMAT_FUNC_OUTPUT_NUM_LIMIT: "3"
  # # 是否检查函数返回值，是否需要被命名
  # GO_FORMAT_FUNC_OUTPUT_IF_NEED_NAMED: "true"
  # # 检查函数参数数量限制时，是否需要排除 context
  # GO_FORMAT_FUNC_INPUT_NUM_WITHOUT_CONTEXT: "true"
  # # 检查函数返回值数量限制时，是否需要排除 error
  # GO_FORMAT_FUNC_OUTPUT_NUM_WITHOUT_ERROR: "true"
  # # 是否检查 context 时函数的第一个参数
  # GO_FORMAT_FUNC_INPUT_CONTEXT_FIRST: "true"
  # # 是否检查 error 时函数的最后一个返回值
  # GO_FORMAT_FUNC_OUTPUT_ERROR_LAST: "true"
  # # 是否检查函数参数的 context 和 返回值的 error，是否同时为 go 标准的类型，或者同时为自定义的类型
  # GO_FORMAT_FUNC_CONTEXT_ERROR_MATCH: "false"
  # # 检查 go 代码格式，排除的文件选项
  # # 请参照 GO_EXCLUDE_OPTIONS 说明
  # # 检查单测覆盖率时，优先选择 GO_FORMAT_EXCLUDES
  # # 如果  GO_FORMAT_EXCLUDES 为空字符，则使用GO_EXCLUDE_OPTIONS
  # GO_FORMAT_EXCLUDES: ""
  # *** go format variables end ***

  # *** go coverage variables begin ***
  # # 检查单测覆盖率的目标分支名
  # # 默认: MR 的目标分支
  # # 也可以选择固定 release 分支，比如在提测之后修改Bug，单测可能覆盖不到改动的代码
  # GO_COVERAGE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 检查单测覆盖率的源分支名
  # GO_COVERAGE_DIFF_SOURCE: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
  # GO_COVERAGE_FILE: "$CI_PROJECT_DIR/log/coverage.out"
  # # 单测覆盖率的阈值，取值范围 0.0 ~ 1.0
  # GO_COVERAGE_RATE: "0.65"
  # # 检查 go 代码的单测覆盖率时，排除的文件选项
  # # 请参照 GO_EXCLUDE_OPTIONS 说明
  # # 检查单测覆盖率时，优先选择GO_COVERAGE_EXCLUDES
  # # 如果 GO_COVERAGE_EXCLUDES 为空字符，则使用GO_EXCLUDE_OPTIONS
  # GO_COVERAGE_EXCLUDES: ""
  # *** go coverage variables end ***

  # *** go cognitive variables begin ***
  # # 检查认知复杂度的目标分支名
  # # 默认: MR 的目标分支
  # GO_COGNITIVE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 检查认知复杂度的源分支名
  # GO_COGNITIVE_DIFF_SOURCE: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
  # # 函数认知复杂度的阈值
  # GO_COGNITIVE_OVER: "15"
  # # 检查 go 代码的认知复杂度时，排除的文件选项
  # # 请参照 GO_EXCLUDE_OPTIONS 说明
  # # 检查单测覆盖率时，优先选择GO_COGNITIVE_EXCLUDES
  # # 如果 GO_COGNITIVE_EXCLUDES 为空字符，则使用GO_EXCLUDE_OPTIONS
  # GO_COGNITIVE_EXCLUDES: ""
  # *** go cognitive variables end ***

  # *** check for test variables begin ***
  # # 检查 test 代码的关键字，防止合入 test 代码
  # FOR_TEST_KEYWORD: todo\s*for\s*test
  # *** check for test variables end ***
```

## git hook 接入

### Leader 先把 git-go-tool 加到项目的 submodule

请在项目的根目录执行以下命令

```shell
git submodule add https://github.com/distroy/git-go-tool.git git-go-tool
```

### Member 把项目的修改拉取到本地，再设置git hook

请在项目的根目录执行以下命令

```shell
git submodule init
git submodule update
git config core.hooksPath "git-go-tool/git-hook"
```

可以参考 [makefile](doc/template/makefile) 的 setup

### git hook 更新

请进入到项目的git-go-tool目录执行

```shell
git pull
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

# 认知复杂度介绍

- 参考文档：[CognitiveComplexity.pdf (sonarsource.com)](https://www.sonarsource.com/docs/CognitiveComplexity.pdf)
- 中文文档：[认知复杂度—估算项目代码的理解成本](https://blog.csdn.net/tjgykhulj/article/details/106569894)
- 计算示例：[core/gocognitive/example_for_test.go](core/gocognitive/example_for_test.go)
