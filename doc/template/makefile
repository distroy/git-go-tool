#
# Copyright (C) distroy
#

env ?= test
export env
cid ?= sg
export cid

goc_flag ?=


# variables
ENV=${env}
REGION=${cid}
PROJECT=adp
PROJECT_ROOT=$(patsubst %/,%,$(abspath $(dir $$PWD)))
TARGETS=$(sort $(notdir $(patsubst %/,%,$(dir $(wildcard $(PROJECT_ROOT)/app/*/*.go)))))
PROTOS=$(dir $(shell test -d $(PROJECT_ROOT)/proto && find $(PROJECT_ROOT)/proto -name '*.proto'))
$(info ENV: $(ENV))
$(info REGION: $(REGION))
$(info PROJECT_ROOT: $(PROJECT_ROOT))
$(info TARGETS: $(TARGETS))
$(info PROTOS: $(PROTOS))

# go
GO_DEP=
GO=go
GOC=goc
GO_FLAGS=${flags}
GO_VERSION=$(shell go version | cut -d" " -f 3)
GO_MAJOR_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 1)
GO_SUB_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 2)
export GO111MODULE=on
ifeq ($(shell expr ${GO_SUB_VERSION} '>' 10), 1)
	GO_FLAGS+=-mod=vendor
endif
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

mk_service = ( \
	echo "=== building service: $(PROJECT_ROOT)/app/$(1)"; \
	cd $(PROJECT_ROOT)/app/$(1); \
	echo $(GO) build $(GO_FLAGS) -o $(1) .; \
	$(GO) build $(GO_FLAGS) -o $(1) . || exit $$?; \
	cp $(1) $(PROJECT_ROOT)/$(PROJECT)-$(1); \
	cd $(PROJECT_ROOT); \
	);

protobuf = ( \
	echo "=== building protobuf: $(1)"; \
	cd $(1); \
	echo protoc --go_out . --python_out . *.proto; \
	protoc --go_out . --python_out . *.proto || exit $$?; \
	cd $(PROJECT_ROOT); \
	);

rm_service = \
	rm -f $(PROJECT_ROOT)/app/$(1)/$(1) $(PROJECT_ROOT)/$(PROJECT)-$(1);

go_install = \
	go install $(1)@latest || go install $(1)

go_install_cmd =  ( \
	_work_dir="$$PWD"; \
	cd /; \
	echo go install github.com/distroy/git-go-tool/cmd/$(1) \
		"||" go install github.com/distroy/git-go-tool/cmd/$(1)@latest; \
	go install github.com/distroy/git-go-tool/cmd/$(1) \
		|| go install github.com/distroy/git-go-tool/cmd/$(1)@latest; \
	cd "$$_work_dir"; \
	)

git_ignore_service = ( \
	_ignore_file="$(PROJECT_ROOT)/.gitignore"; \
	_service_bin="/app/$(1)/$(1)"; \
	_count=$$(grep "$$_service_bin" "$$_ignore_file" | wc -l); \
	test $$_count -ne 0 || echo "$$_service_bin" >> "$$_ignore_file"; \
	);

ifneq ($(goc_flag),)
ifeq ($(ENV), $(filter $(ENV), dev test))

GO_DEP += goc-download
GO=$(GOC)
GO_FLAGS := --buildflags='$(GO_FLAGS)'
$(info GO_FLAGS: $(GO_FLAGS))

endif # ifeq ($(env), $(filter $(env), dev test))
endif # ifneq ($(goc_flag),)

.PHONY: all
all: setup $(TARGETS)

.PHONY: pb
pb:
	@$(foreach i, $(PROTOS), $(call protobuf,$(i)))

.PHONY: $(TARGETS)
$(TARGETS): git-ignore $(GO_DEP)
	@$(call mk_service,$@)

.PHONY: goc-download
goc-download:
	# go install github.com/qiniu/goc
	# chmod u+x $(GOC)

.PHONY: clean
clean: git-ignore
	@$(foreach service, $(TARGETS), $(call rm_service,$(service)))

.PHONY: dep
dep: setup
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: dep-backup
dep-backup:
	test -d vendor && mv vendor vendor.bak
	echo >> ~/.bashrc
	echo "source $(PROJECT_ROOT)/bash/bashrc" >> ~/.bashrc

.PHONY: update
update: dep pb

.PHONY: build-test
build-test: $(TARGETS) clean

.PHONY: go-test-report-dir
go-test-report-dir:
	mkdir $(GO_TEST_REPORT_DIR) -pv

.PHONY: go-test-coverage
go-test-coverage: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/go-coverage.out"
	$(GO) tool cover -html $(GO_TEST_REPORT_DIR)/go-coverage.out \
		-o $(GO_TEST_REPORT_DIR)/go-coverage.html

.PHONY: go-test-report
go-test-report: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/go-coverage.out" \
		-json > "$(GO_TEST_REPORT_DIR)/go-test.json"
	$(GO) tool cover -html $(GO_TEST_REPORT_DIR)/go-coverage.out \
		-o $(GO_TEST_REPORT_DIR)/go-coverage.html

.PHONY: go-test
go-test:
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./...

BLOCK_HOSTS_FILE=$(PROJECT_ROOT)/block-hosts
.PHONY: block-hosts
block-hosts:
	touch $(BLOCK_HOSTS_FILE)
	grep -v '^#' "$(BLOCK_HOSTS_FILE)" | \
		while read i; do \
			test -z "$$i" || echo "0.0.0.0  $$i"; \
		done >> /etc/hosts
	cat /etc/hosts

.PHONY: git-ignore
git-ignore:
	@$(foreach service, $(TARGETS), $(call git_ignore_service,$(service)))

.PHONY: setup
setup: git-ignore
	git submodule init
	git submodule update
	git config core.hooksPath "git-go-tool/git-hook"
	@$(call go_install_cmd,go-cognitive)
	@$(call go_install_cmd,go-format)
	@$(call go_install_cmd,git-diff-go-cognitive)
	@$(call go_install_cmd,git-diff-go-coverage)
	@$(call go_install_cmd,git-diff-go-format)
	@echo $$'\E[32;1m'"setup succ"$$'\E[0m'

.PHONY: go-cognitive
go-cognitive: setup
	go-cognitive

.PHONY: go-format
go-format: setup
	go-format
