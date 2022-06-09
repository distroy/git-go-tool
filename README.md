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
  # # 执行 go-check-stage 的镜像
  # GO_CHECK_IMAGE: 'harbor.shopeemobile.com/shopee/golang-base'

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
