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
- `tc_rce_option` (Optional, TypeList): TC-RCE authentication configuration, required when attester_source is TC-RCE
- `tc_captcha_option` (Optional, TypeList): TC-CAPTCHA authentication configuration, required when attester_source is TC-CAPTCHA
- `tc_eo_captcha_option` (Optional, TypeList): TC-EO-CAPTCHA authentication configuration, required when attester_source is TC-EO-CAPTCHA
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
