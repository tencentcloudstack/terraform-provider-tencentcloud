# Change: Add Chinese Character Check GitHub Workflow

## Why
The Terraform provider is an international project that must ensure all public-facing code and documentation does not contain Chinese characters. Recent feedback indicates that some documentation contains Chinese characters, which violates the international project standards. We need automated validation to prevent Chinese characters from being introduced in future PRs.

## What Changes
- Add a new GitHub workflow `.github/workflows/chinese-char-check.yml` that validates PR content
- **BREAKING**: PRs containing Chinese characters in specified directories will fail CI checks
- Validation scope limited to:
  - `/tencentcloud/services/` directory (provider code)
  - `/website/` directory (documentation)
- Check only applies to changed files in PR, not historical content
- Workflow triggers on PR events (opened, synchronize, edited)

## Impact
- Affected specs: ci-workflows (new capability)
- Affected code: 
  - `.github/workflows/` (new workflow file)
  - CI pipeline behavior (new validation step)
- Breaking change: PRs with Chinese characters will be blocked until fixed
- Developers must ensure code/docs are in English before submitting PRs