## ADDED Requirements

### Requirement: Resource CRUD for tencentcloud_teo_security_client_attester
The system SHALL provide a Terraform resource `tencentcloud_teo_security_client_attester` that manages the full lifecycle of TEO client attestation options through Create, Read, Update, and Delete operations.

The resource SHALL use `zone_id` as a required parameter and `client_attesters` as a required parameter representing the list of client attestation options.

The resource SHALL use a composite ID formed by `zone_id` and `client_attester_ids` joined with `tccommon.FILED_SP` as the separator.

#### Scenario: Create a new client attester resource
- **WHEN** a user creates a `tencentcloud_teo_security_client_attester` resource with `zone_id` and `client_attesters`
- **THEN** the system SHALL call `CreateSecurityClientAttester` API with the provided parameters
- **THEN** the system SHALL store the returned `client_attester_ids` as part of the composite resource ID
- **THEN** the system SHALL call Read to refresh the resource state

#### Scenario: Read client attester resource
- **WHEN** the system reads a `tencentcloud_teo_security_client_attester` resource
- **THEN** the system SHALL call `DescribeSecurityClientAttester` API with `zone_id` from the composite ID
- **THEN** the system SHALL use Limit=100 (maximum) for pagination to fetch all client attesters
- **THEN** the system SHALL flatten the returned `ClientAttesters` into the Terraform state
- **THEN** if the resource is not found, the system SHALL set the resource ID to empty string

#### Scenario: Update client attester resource
- **WHEN** a user updates mutable fields of a `tencentcloud_teo_security_client_attester` resource
- **THEN** the system SHALL check that immutable fields (`zone_id`) have not changed
- **THEN** if immutable fields changed, the system SHALL return an error
- **THEN** the system SHALL call `ModifySecurityClientAttester` API with `zone_id` and updated `client_attesters`
- **THEN** the system SHALL call Read to refresh the resource state

#### Scenario: Delete client attester resource
- **WHEN** a user deletes a `tencentcloud_teo_security_client_attester` resource
- **THEN** the system SHALL call `DeleteSecurityClientAttester` API with `zone_id` and `client_attester_ids` extracted from the composite ID
- **THEN** the system SHALL set the resource ID to empty string after successful deletion

#### Scenario: Import existing client attester resource
- **WHEN** a user imports a `tencentcloud_teo_security_client_attester` resource using the composite ID format `zone_id#client_attester_ids`
- **THEN** the system SHALL parse the composite ID
- **THEN** the system SHALL call Read to populate the resource state

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

### Requirement: Computed client_attester_ids attribute
The resource SHALL expose a `client_attester_ids` computed attribute of TypeList/TypeSet that contains the IDs of the created client attestation options returned by the CreateSecurityClientAttester API response.

#### Scenario: client_attester_ids populated after create
- **WHEN** a client attester resource is created
- **THEN** the `client_attester_ids` attribute SHALL be populated from the `CreateSecurityClientAttester` response

### Requirement: API retry handling
All CRUD operations SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` for write operations (Create, Update, Delete) and `tccommon.ReadRetryTimeout` for read operations. Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: API call with retry on transient error
- **WHEN** an API call fails with a transient error
- **THEN** the system SHALL retry the operation within the timeout period
- **THEN** if retries are exhausted, the system SHALL return the wrapped error

### Requirement: Resource registration in provider
The resource `tencentcloud_teo_security_client_attester` SHALL be registered in `provider.go` ResourcesMap and documented in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the resource `tencentcloud_teo_security_client_attester` SHALL be available for use in Terraform configurations

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
