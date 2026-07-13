## ADDED Requirements

### Requirement: Chinese Character Validation in PRs
The CI system SHALL validate that PR content does not contain Chinese characters in public-facing directories to maintain international project standards.

#### Scenario: PR with Chinese characters is rejected
- **WHEN** a PR is submitted with Chinese characters in `/tencentcloud/services/` or `/website/` directories
- **THEN** the CI check fails with clear error message indicating which files contain Chinese characters
- **AND** the PR cannot be merged until Chinese characters are removed

#### Scenario: PR with only English content passes validation
- **WHEN** a PR is submitted with only English content in all changed files
- **THEN** the Chinese character validation check passes
- **AND** the PR can proceed through other CI checks

#### Scenario: Chinese characters outside scope are ignored
- **WHEN** a PR contains Chinese characters in files outside `/tencentcloud/services/` and `/website/` directories
- **THEN** the Chinese character validation check passes
- **AND** the workflow does not report these characters as violations

#### Scenario: Validation only checks changed files
- **WHEN** a PR is submitted that modifies existing files
- **THEN** only the changed lines/files are checked for Chinese characters
- **AND** existing Chinese characters in unchanged files are not flagged

### Requirement: Workflow Configuration
The Chinese character check workflow SHALL be configured to run automatically on all PR events and provide clear feedback.

#### Scenario: Workflow triggers on PR events
- **WHEN** a PR is opened, synchronized, or edited
- **THEN** the Chinese character check workflow is automatically triggered
- **AND** the check runs within the standard CI pipeline

#### Scenario: Clear error reporting
- **WHEN** Chinese characters are detected in the validation scope
- **THEN** the workflow reports the specific files and line numbers containing Chinese characters
- **AND** provides guidance on how to resolve the issue

#### Scenario: Workflow performance
- **WHEN** the Chinese character check runs
- **THEN** it completes within 2 minutes for typical PR sizes
- **AND** does not significantly impact overall CI pipeline duration