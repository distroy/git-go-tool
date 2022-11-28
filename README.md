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

  # *** go flags begin ***
  # # go 编译和单测的选项
  # GO_FLAGS: "-mod=vendor"
  # # go 单测的选项
  # GO_TEST_FLAGS: "-gcflags=all=-l"
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
  # *** go format variables end ***

  # *** go coverage variables begin ***
  # # 检查单测覆盖率的目标分支名
  # # 默认: MR 的目标分支
  # # 也可以选择固定 release 分支，比如在提测之后修改Bug，单测可能覆盖不到改动的代码
  # GO_COVERAGE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 检查单测覆盖率的源分支名
  # GO_COVERAGE_DIFF_SOURCE: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
  # GO_COVERAGE_FILE: "$CI_PROJECT_DIR/log/coverage.out"
  # *** go coverage variables end ***

  # *** go cognitive variables begin ***
  # # 检查认知复杂度的目标分支名
  # # 默认: MR 的目标分支
  # GO_COGNITIVE_DIFF_TARGET: "$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
  # # 检查认知复杂度的源分支名
  # GO_COGNITIVE_DIFF_SOURCE: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
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


# 配置文件

配置文件路径: .git-go-tool/config.yaml

```yaml
# 单测覆盖率相关的配置
go-coverage:
  git-diff:
    # # git 的比较模式, default: diff
    # #   diff: 增量模式，只检查新增代码的认知复杂度和单测覆盖率
    # #   all: 全量模式，检查所有代码的认知复杂度和单测覆盖率
    # mode: diff
  # # 检查 go 代码的单测覆盖率时，包含的文件选项, 优先级高于 exclude。正则表达式的列表
  include:
  # # 检查 go 代码的单测覆盖率时，排除的文件选项。正则表达式的列表
  exclude:
    # - '(^|/)vendor/'
    # - '\.pb\.go$'
  # # 单测覆盖率的阈值，取值范围 0.0 ~ 1.0, default: 0.65
  # rate: 0.65
  # # 展示 Top N 没有被覆盖的代码文件, default: 10
  # top: 10

# go 认知复杂度相关的配置
go-cognitive:
  git-diff:
    # # git 的比较模式, default: diff
    # mode: diff
  # # 检查 go 代码的认知复杂度时，包含的文件选项, 优先级高于 exclude。正则表达式的列表
  include:
  # # 检查 go 代码的认知复杂度时，排除的文件选项。正则表达式的列表
  exclude:
  # # 函数认知复杂度的阈值, default: 15
  # over: 15
  # # 展示 Top N 复杂度过高的函数, default: 10
  # top: 10

# go-format:
go-format:
  git-diff:
    # # git 的比较模式, default: diff
    # mode: diff
  # # 检查 go 代码的格式时，包含的文件选项, 优先级高于 exclude。正则表达式的列表
  include:
  # # 检查 go 代码的格式时，排除的文件选项。正则表达式的列表
  exclude:
  # # go 代码文件函数限制, default: 1000
  # file-line: 1000
  # # 是否检查 import 格式, default: true
  # import: true
  # # 是否检查源代码被 go 提供的工具格式化了, default: true
  # formated: true
  # # 是否检查 package 命名, default: true
  # package: true
  # 函数参数的数量限制，0-不检查, default: 3
  # func-input-num: 3
  # # 函数返回值的数量限制，0-不检查, default: 3
  # func-output-num: 3
  # # 是否检查函数返回值，是否需要被命名, default: true
  # func-named-output: true
  # # 检查函数参数数量限制时，是否需要排除 context, default: true
  # func-input-num-without-context: true
  # # 检查函数返回值数量限制时，是否需要排除 error, default: true
  # func-output-num-without-error: true
  # # 是否检查 context 时函数的第一个参数, default: true
  # func-context-first: true
  # # 是否检查 error 时函数的最后一个返回值, default: true
  # func-error-last: true
  # # 是否检查函数参数的 context 和 返回值的 error，是否同时为 go 标准的类型，或者同时为自定义的类型, default: false
  # func-context-error-match: false
```


# Commands 介绍

## Installation

```shell
go install github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive@latest
go install github.com/distroy/git-go-tool/cmd/git-diff-go-coverage@latest
go install github.com/distroy/git-go-tool/cmd/git-diff-go-format@latest
go install github.com/distroy/git-go-tool/cmd/go-cognitive@latest
go install github.com/distroy/git-go-tool/cmd/go-format@latest
```

## go-cognitive
go-cognitive analyze cognitive complexities of functions in Go source code. A measurement of how hard does the code is intuitively to understand.

# 认知复杂度介绍

- 参考文档：[CognitiveComplexity.pdf (sonarsource.com)](https://www.sonarsource.com/docs/CognitiveComplexity.pdf)
- 中文文档：[认知复杂度—估算项目代码的理解成本](https://blog.csdn.net/tjgykhulj/article/details/106569894)
- 计算示例：[core/gocognitive/example_for_test.go](core/gocognitive/example_for_test.go)
