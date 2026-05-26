TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=tencentcloud
CHANGED_FILES=$$(git diff --name-only master -- $(PKG_NAME) | grep '.go$$')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PLATFORMS=darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm openbsd/amd64 openbsd/386 solaris/amd64 windows/386 windows/amd64
GO_VER ?= go

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

# check-mux validates SDKv2 + framework dual-stack mux compatibility without
# making any remote calls. It boots tf5muxserver to merge schemas (including
# an actual GetProviderSchema call so byte-level schema drift that
# NewMuxServer cannot detect is caught here), and asserts that the resource /
# data source type-name sets across the two stacks are disjoint. Any PR that
# adds framework references or changes the SDKv2 / framework provider top
# Schema / MetaSchema MUST keep this target green.
check-mux: fmtcheck
	@echo "==> Checking SDKv2/framework provider mux compatibility..."
	@go test -count=1 -run 'TestMuxServer_NoStartupError|TestMuxServer_GetProviderSchemaNoError|TestFrameworkProvider_NoTypeNameCollision' ./$(PKG_NAME)/framework/


testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

fmt-faster:
	@if [[ -z $(CHANGED_FILES) ]]; then \
		echo "skip the fmt cause the CHANGED_FILES is null."; \
		exit 0; \
	else \
		echo "==> [Faster]Fixing source code with gofmt...\n $(CHANGED_FILES) \n"; \
		gofmt -s -w $(CHANGED_FILES); \
	fi

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

deltatest:
	@sh -c "'$(CURDIR)/scripts/delta-test.sh'"

changelogtest:
	@sh -c "'$(CURDIR)/scripts/changelog-test.sh'"

dispachertest:
	@sh -c "'$(CURDIR)/scripts/dispacher-test.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=30 GOPACKAGESPRINTGOLISTERRORS=1 golangci-lint run --timeout=30m ./$(PKG_NAME)
	@tfproviderlint \
		-c 1 \
		-AT002 \
		-AT005 \
		-AT006 \
		-AT007 \
		-AT008 \
		-R003 \
		-R012 \
		-R013 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		-S020 \
		-S021 \
		-S023 \
		-S024 \
		-S025 \
		-S026 \
		-S027 \
		-S028 \
		-S029 \
		-S030 \
		-S031 \
		-S032 \
		-S033 \
		-S034 \
		-S035 \
		-S036 \
		-S037 \
		-V002 \
		-V003 \
		-V004 \
		-V005 \
		-V006 \
		-V007 \
		-V008 \
		./$(PKG_NAME)

tools:
	GO111MODULE=on cd .ci/tools && go install github.com/bflad/tfproviderlint/cmd/tfproviderlint && cd ../..
	GO111MODULE=on cd .ci/tools && go install github.com/client9/misspell/cmd/misspell && cd ../..
	GO111MODULE=on cd .ci/tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2 && cd ../..

# tfproviderlint-local mirrors .github/workflows/tfprovider-lint.yml byte for
# byte locally: it installs the same tfproviderlint v0.31.0 used by CI
# (already verified to compile under Go 1.25) and runs the identical set of
# checkers (AT/R/S/V) against ./tencentcloud. Whenever the CI yaml changes,
# please keep this target in sync.
TFPROVIDERLINT_VERSION ?= v0.31.0
TFPROVIDERLINT_BIN ?= $(shell go env GOPATH)/bin/tfproviderlint

tfproviderlint-install:
	@echo "==> Installing tfproviderlint $(TFPROVIDERLINT_VERSION)..."
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@$(TFPROVIDERLINT_VERSION)

tfproviderlint-local: tfproviderlint-install
	@echo "==> Running tfproviderlint against ./$(PKG_NAME) (mirrors CI)..."
	@$(TFPROVIDERLINT_BIN) \
		-c 1 \
		-AT002 -AT005 -AT006 -AT007 -AT008 \
		-R003 -R012 -R013 \
		-S001 -S002 -S003 -S004 -S005 -S007 -S008 -S009 -S010 -S011 -S012 \
		-S013 -S014 -S015 -S016 -S017 -S019 -S020 -S021 -S023 -S024 -S025 \
		-S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033 -S034 -S035 -S036 -S037 \
		-V002 -V003 -V004 -V005 -V006 -V007 -V008 \
		./$(PKG_NAME)

# golangci-lint-local mirrors .github/workflows/golangci-lint.yml byte for
# byte locally: it installs the same golangci-lint v2.4.0 used by CI (built
# with the Go 1.25 toolchain) and runs from inside ./tencentcloud, picking up
# the repo-root .golangci.yml and reporting only issues introduced relative
# to master.
GOLANGCI_LINT_VERSION ?= v2.4.0
GOLANGCI_LINT_BIN ?= $(shell go env GOPATH)/bin/golangci-lint

golangci-lint-install:
	@echo "==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

golangci-lint-local: golangci-lint-install
	@echo "==> Running golangci-lint against ./$(PKG_NAME) (mirrors CI)..."
	@cd ./$(PKG_NAME) && $(GOLANGCI_LINT_BIN) run --new-from-rev=origin/master -v

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

test-build-all:
	@$(foreach platform, $(PLATFORMS), \
		echo GOOS=$(firstword $(subst /, ,$(platform))) GOARCH=$(lastword $(subst /, ,$(platform))) \
		go build -o terraform-provider-tencentcloud; \
		GOOS=$(firstword $(subst /, ,$(platform))) GOARCH=$(lastword $(subst /, ,$(platform))) \
		go build -o terraform-provider-tencentcloud; \
		rm -f terraform-provider-tencentcloud; \
	)

test-build: test-build-x64 test-build-x86

test-build-x64:
	GOARCH=amd64 go build -o terraform-provider-tencentcloud-amd64
	rm -f terraform-provider-tencentcloud-amd64

test-build-x86:
	GOARCH=386 go build -o terraform-provider-tencentcloud-386
	rm -f terraform-provider-tencentcloud-386

doc:
	cd gendoc && go run ./... && cd ..

doc-faster:
	@echo "==> [Faster]Generating doc..."
	@if [ ! -f gendoc/gendoc ]; then \
		$(MAKE) doc-bin-build; \
	fi
	@$(MAKE) doc-with-bin

doc-with-bin:
	cd gendoc && ./gendoc ./... && cd ..

doc-bin-build:
	@echo "==> Building gendoc binary..."
	cd gendoc && go build ./... && cd ..

hooks: tools
	@find .git/hooks -type l -exec rm {} \;
	@find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
	@echo "==> Install hooks done."

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	ln -sf ../../../../ext/providers/tencentcloud/website/docs $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/docs/providers/tencentcloud
	ln -sf ../../../ext/providers/tencentcloud/website/tencentcloud.erb $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/layouts/tencentcloud.erb
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/ || (echo; \
		echo "Unexpected mispelling found in website files."; \
		echo "To automatically fix the misspelling, run 'make website-lint-fix' and commit the changes."; \
		exit 1)
	@terrafmt diff ./website --check --pattern '*.markdown' --quiet || (echo; \
		echo "Unexpected differences in website HCL formatting."; \
		echo "To see the full differences, run: terrafmt diff ./website --pattern '*.markdown'"; \
		echo "To automatically fix the formatting, run 'make website-lint-fix' and commit the changes."; \
		exit 1)

website-lint-fix:
	@echo "==> Applying automatic website linter fixes..."
	@misspell -w -source=text website/
	@terrafmt fmt ./website --pattern '*.markdown'

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	ln -sf ../../../../ext/providers/tencentcloud/website/docs $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/docs/providers/tencentcloud
	ln -sf ../../../ext/providers/tencentcloud/website/tencentcloud.erb $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/layouts/tencentcloud.erb
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

changelog:
	./scripts/generate-changelog.sh

.PHONY: build sweep test testacc fmt fmtcheck lint tools tfproviderlint-install tfproviderlint-local golangci-lint-install golangci-lint-local test-compile doc hooks website website-lint website-test check-mux

ready: doc fmt-faster
