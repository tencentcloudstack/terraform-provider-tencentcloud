## MODIFIED Requirements

### Requirement: Client attesters schema definition
The resource SHALL define a `client_attesters` parameter as a TypeList with nested block schema containing the following fields:

- `name` (Required, TypeString): Attestation option name
- `attester_source` (Required, TypeString): Authentication method, values: TC-RCE, TC-CAPTCHA, TC-EO-CAPTCHA
- `attester_duration` (Optional, TypeString): Authentication validity duration
- `usage_limit` (Optional, TypeString): Whether to enable authentication result usage limit. Values: `on` (enable), `off` (disable). Default `off`.
- `max_usage_count` (Optional, TypeInt): Maximum usage count of authentication result. Required when `usage_limit` is `on`. Range [1, 100000].
- `tc_rce_option` (Optional, TypeList): TC-RCE authentication configuration, required when attester_source is TC-RCE
- `tc_captcha_option` (Optional, TypeList): TC-CAPTCHA authentication configuration, required when attester_source is TC-CAPTCHA
- `tc_eo_captcha_option` (Optional, TypeList): TC-EO-CAPTCHA authentication configuration, required when attester_source is TC-EO-CAPTCHA. Contains `captcha_mode` (Required, TypeString) and `mobile_captcha_mode` (Optional, TypeString, values: `Invisible`, `Visible`)
- `id` (Computed, TypeString): Attestation option ID, returned by the API
- `type` (Computed, TypeString): Rule type, returned by the API (PRESET/CUSTOM)

#### Scenario: Client attester with TC-RCE authentication
- **WHEN** a user creates a client attester with `attester_source` = "TC-RCE"
- **THEN** the `tc_rce_option` block SHALL be required with `channel` and `region` fields

#### Scenario: Client attester with TC-CAPTCHA authentication
- **WHEN** a user creates a client attester with `attester_source` = "TC-CAPTCHA"
- **THEN** the `tc_captcha_option` block SHALL be required with `captcha_app_id` and `app_secret_key` fields

#### Scenario: Client attester with TC-EO-CAPTCHA authentication
- **WHEN** a user creates a client attester with `attester_source` = "TC-EO-CAPTCHA"
- **THEN** the `tc_eo_captcha_option` block SHALL be required with `captcha_mode` field
- **AND** the `mobile_captcha_mode` field SHALL be optional with values `Invisible` or `Visible`

#### Scenario: Client attester with usage limit enabled
- **WHEN** a user sets `usage_limit` = "on" and `max_usage_count` = 100 in the `client_attesters` block
- **THEN** the schema SHALL accept the values and pass them to the `ClientAttesters[].UsageLimit` and `ClientAttesters[].MaxUsageCount` API fields on Create

#### Scenario: Client attester without usage limit configured
- **WHEN** a user does not set `usage_limit` or `max_usage_count` in the configuration
- **THEN** the schema SHALL treat them as optional and NOT send them to the cloud API (omitnil behavior)

### Requirement: Resource Create/Read/Update with new fields
The system SHALL handle the new `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` fields consistently across Create, Read, and Update operations.

- On Create, the `buildClientAttesterFromMap` helper SHALL populate `UsageLimit`, `MaxUsageCount`, and `TCEOCaptchaOption.MobileCaptchaMode` only when their values are non-zero (non-empty for strings, non-zero for int).
- On Update, the same helper SHALL be used, since Create and Modify share the `ClientAttester` request structure.
- On Read, the system SHALL backfill `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` from the `DescribeSecurityClientAttester` response into state, but only when the corresponding response fields are non-nil.

#### Scenario: Create passes new fields to API
- **WHEN** a user creates a resource with `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` configured
- **THEN** the `CreateSecurityClientAttester` request SHALL include these fields in the `ClientAttesters` element

#### Scenario: Read backfills new fields from API
- **WHEN** the system reads a resource and the API response contains `UsageLimit`, `MaxUsageCount`, and `MobileCaptchaMode`
- **THEN** the state SHALL be populated with `usage_limit`, `max_usage_count`, and `mobile_captcha_mode`

#### Scenario: Read skips nil response fields
- **WHEN** the system reads a resource and the API response has nil `UsageLimit`, `MaxUsageCount`, or `MobileCaptchaMode`
- **THEN** the system SHALL NOT set the corresponding state key (avoid panic)

#### Scenario: Update passes new fields to API
- **WHEN** a user updates `usage_limit`, `max_usage_count`, or `mobile_captcha_mode`
- **THEN** the `ModifySecurityClientAttester` request SHALL include the updated values in the `ClientAttesters` element

## ADDED Requirements

### Requirement: Unit tests for new client attester fields
The system SHALL provide unit tests in `resource_tc_teo_security_client_attester_test.go` using gomonkey to mock cloud API calls, covering the new `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` fields.

#### Scenario: Unit tests pass for new fields
- **WHEN** `go test` is run with `-gcflags=all=-l` on the test file
- **THEN** all test cases for Create, Read, Update with `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` SHALL pass

### Requirement: Resource documentation for new fields
The system SHALL update the markdown documentation `resource_tc_teo_security_client_attester.md` to include the new `usage_limit`, `max_usage_count`, and `mobile_captcha_mode` fields in the example usage.

#### Scenario: Documentation includes new fields
- **WHEN** the resource documentation is generated
- **THEN** the example usage SHALL demonstrate `usage_limit`, `max_usage_count` (when applicable) and `mobile_captcha_mode` (within `tc_eo_captcha_option`)