<!--
Thanks for contributing to terraform-provider-tencentcloud!
Please fill in the sections below to help reviewers understand your change.
See CONTRIBUTING.md for project conventions.
-->

## What

<!-- Describe the change in 1-3 sentences. -->

## Why

<!-- The problem this change solves, or the new capability it enables. -->

## Type of change

- [ ] New SDKv2 resource / data source
- [ ] New framework resource / data source
- [ ] Modification to an existing resource / data source
- [ ] Bug fix
- [ ] Refactor / internal-only change
- [ ] Documentation only

## Checklist

- [ ] `make build` succeeds locally
- [ ] `make fmt` produces no diff
- [ ] `make lint` passes (or pre-existing failures are unrelated)
- [ ] `make check-mux` passes (SDKv2 + framework mux compatibility)
- [ ] **Confirmed the resource/data source type name is not already
      registered in the other stack** (CI's `make check-mux` enforces
      this; please double-check before opening the PR)
- [ ] Acceptance tests added or updated (or skipped with justification)
- [ ] Website docs updated under `website/docs/r/` or `website/docs/d/`
- [ ] `go.mod` changes are accompanied by `go mod vendor` in the same commit
- [ ] Changelog entry added under `.changelog/<PR>.txt` (if user-facing)

## Notes for reviewers

<!-- Anything reviewers should pay extra attention to: tricky logic,
unusual SDK quirks, breaking changes, etc. -->
