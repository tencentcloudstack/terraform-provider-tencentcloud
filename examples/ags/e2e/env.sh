#!/usr/bin/env bash
# Source this file for the isolated AGS build and test environment.
AGS_REPO_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
AGS_WORK_ROOT="${AGS_WORK_ROOT:-$(dirname -- "${AGS_REPO_ROOT}")/.ags-terraform-work}"
export GOCACHE="${AGS_WORK_ROOT}/cache/go-build"
export GOMODCACHE="${AGS_WORK_ROOT}/cache/mod"
export GOBIN="${AGS_TOOLS_DIR:-/tmp/ags-terraform-tools}"
export TMPDIR="${AGS_WORK_ROOT}/tmp"
export GOTMPDIR="${TMPDIR}"
export GOTOOLCHAIN=local
export GOLANGCI_LINT_CACHE="${AGS_WORK_ROOT}/cache/golangci-lint"
export GOMAXPROCS="${AGS_BUILD_CONCURRENCY:-2}"
export GOFLAGS="-p=${AGS_BUILD_CONCURRENCY:-2}"
export PATH="${GOBIN}:${HOME}/.local/bin:${PATH}"
mkdir -p "${GOCACHE}" "${GOMODCACHE}" "${GOBIN}" "${TMPDIR}" "${GOLANGCI_LINT_CACHE}"
