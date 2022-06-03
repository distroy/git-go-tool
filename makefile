#
# Copyright (C) distroy
#

# variables
PROJECT_ROOT=$(patsubst %/,%,$(abspath $(dir $$PWD)))
COMMANDS=$(sort $(notdir $(patsubst %/,%,$(dir $(wildcard $(PROJECT_ROOT)/cmd/*/*.go)))))
$(info PROJECT_ROOT: $(PROJECT_ROOT))
$(info COMMANDS: $(COMMANDS))

GO=go
GO_FLAGS=${flags}
GO_VERSION=$(shell go version | cut -d" " -f 3)
GO_MAJOR_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 1)
GO_SUB_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 2)
export GO111MODULE=on
# ifeq ($(shell expr ${GO_SUB_VERSION} '>' 10), 1)
# 	GO_FLAGS+=-mod=vendor
# endif
$(info GO_VERSION: $(GO_MAJOR_VERSION).$(GO_SUB_VERSION))
$(info GO_FLAGS: $(GO_FLAGS))

# go test
ifeq (${test_report},)
	export test_report=$(PROJECT_ROOT)/log
endif
GO_TEST_FLAGS+=-v
GO_TEST_FLAGS+=-gcflags="all=-l"
GO_TEST_REPORT_DIR=${test_report}

# git
GIT_REVISION=$(shell git rev-parse HEAD 2> /dev/null)
GIT_BRANCH=$(shell git symbolic-ref HEAD 2> /dev/null | sed -e 's/refs\/heads\///')
GIT_TAG=$(shell git describe --exact-match --tags 2> /dev/null)
$(info GIT_REVISION: $(GIT_REVISION))
$(info GIT_BRANCH: $(GIT_BRANCH))
$(info GIT_TAG: $(GIT_TAG))

mk_command = ( \
	echo "=== building service: $(PROJECT_ROOT)/cmd/$(1)"; \
	cd $(PROJECT_ROOT)/cmd/$(1); \
	echo $(GO) build $(GO_FLAGS) -o $(1) .; \
	$(GO) build $(GO_FLAGS) -o $(1) . || exit $$?; \
	cd $(PROJECT_ROOT); \
	);

rm_command = \
	rm -f $(PROJECT_ROOT)/cmd/$(1)/$(1);

.PHONY: all
all: setup $(COMMANDS)

.PHONY: $(COMMANDS)
$(COMMANDS):
	@$(call mk_command,$@)

.PHONY: clean
clean:
	@$(foreach service, $(COMMANDS), $(call rm_command,$(service)))

.PHONY: dep
dep: setup
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: update
update: dep
	$(MAKE) protocol

.PHONY: build-test
build-test: $(COMMANDS) clean

.PHONY: go-test-report-dir
go-test-report-dir:
	mkdir $(GO_TEST_REPORT_DIR) -pv

.PHONY: go-test-coverage
go-test-coverage: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/coverage.out"

.PHONY: go-test-report
go-test-report: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/coverage.out" \
		-json > "$(GO_TEST_REPORT_DIR)/test.json"

.PHONY: go-test
go-test:
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./...

.PHONY: setup
setup:
	git config core.hooksPath "script/git-hook"
	type go-cognitive \
		|| go install github.com/distroy/git-go-tool/cmd/go-cognitive \
		|| go install github.com/distroy/git-go-tool/cmd/go-cognitive@latest
	@echo $$'\E[32;1m'"setup succ"$$'\E[0m'

.PHONY: complexity
complexity: setup
	go-cognitive -over 15 .
	go-cognitive -top 10 .
