## ADDED Requirements

### Requirement: Create shared CNAME resource
The system SHALL provide a Terraform resource `tencentcloud_teo_shared_cname` that creates a shared CNAME via the `CreateSharedCNAME` API. The resource SHALL accept `zone_id` (Required, ForceNew), `shared_cname_prefix` (Required, ForceNew), and `description` (Optional) as input parameters. Upon successful creation, the resource SHALL store the returned `shared_cname` value and set the resource ID as the composite of `zone_id` and `shared_cname` joined by `tccommon.FILED_SP`. If `ipssl_setting` is provided, the system SHALL immediately call the `ModifySharedCNAME` API after creation to apply the IP SSL setting.

#### Scenario: Successful creation of shared CNAME
- **WHEN** user applies a Terraform config with `tencentcloud_teo_shared_cname` specifying `zone_id`, `shared_cname_prefix`, and `description`
- **THEN** the system calls `CreateSharedCNAME` API, stores the returned `shared_cname` in state, and sets the resource ID as `zone_id + FILED_SP + shared_cname`

#### Scenario: Creation with ipssl_setting
- **WHEN** user applies a Terraform config with `tencentcloud_teo_shared_cname` specifying `zone_id`, `shared_cname_prefix`, `description`, and `ipssl_setting`
- **THEN** the system calls `CreateSharedCNAME` API first, then immediately calls `ModifySharedCNAME` API with the `IPSSLSetting` containing `Status` (mapped to `Operation`: bound/binding→bind, unbound/unbinding→unbind) and `AssociatedDomain` to apply the IP SSL setting

#### Scenario: API returns empty shared_cname on creation
- **WHEN** the `CreateSharedCNAME` API returns a nil or empty `SharedCNAME` value
- **THEN** the system SHALL return a `NonRetryableError` to prevent writing an empty ID to state

### Requirement: Read shared CNAME resource
The system SHALL read the shared CNAME resource using the `DescribeSharedCNAME` API with a filter on `shared-cname` matching the stored `shared_cname` value. The Read method SHALL parse the composite ID to extract `zone_id` and `shared_cname`, then set all computed attributes from the API response.

#### Scenario: Successful read of existing shared CNAME
- **WHEN** Terraform refreshes state for an existing `tencentcloud_teo_shared_cname` resource
- **THEN** the system calls `DescribeSharedCNAME` with `ZoneId` and `Filters` containing `shared-cname` filter, and updates state with the returned `SharedCNAMEInfo` fields

#### Scenario: Resource no longer exists (deleted externally)
- **WHEN** the `DescribeSharedCNAME` API returns an empty `SharedCNAMEInfo` list
- **THEN** the system SHALL log the resource ID for debugging, then call `d.SetId("")` to remove the resource from state

### Requirement: Update shared CNAME resource
The system SHALL support updating the `description` and `ipssl_setting` fields of an existing shared CNAME via the `ModifySharedCNAME` API. The `zone_id` and `shared_cname` fields SHALL be sent as identifiers in the modify request.

#### Scenario: Update description
- **WHEN** user changes the `description` field in the Terraform config
- **THEN** the system calls `ModifySharedCNAME` API with the new `Description` value

#### Scenario: Update ipssl_setting
- **WHEN** user sets or changes the `ipssl_setting` block in the Terraform config
- **THEN** the system calls `ModifySharedCNAME` API with the `IPSSLSetting` containing `Status` (mapped to `Operation`: bound/binding→bind, unbound/unbinding→unbind) and `AssociatedDomain`

### Requirement: Delete shared CNAME resource
The system SHALL delete the shared CNAME resource using the `DeleteSharedCNAME` API with `zone_id` and `shared_cname` extracted from the composite resource ID.

#### Scenario: Successful deletion
- **WHEN** user destroys the `tencentcloud_teo_shared_cname` resource
- **THEN** the system calls `DeleteSharedCNAME` API with the correct `ZoneId` and `SharedCNAME` values

### Requirement: Import shared CNAME resource
The system SHALL support importing an existing shared CNAME resource using the composite ID format `zone_id#shared_cname` (where `#` is `tccommon.FILED_SP`).

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_teo_shared_cname.example zone-xxx#example.com.xxx.share.dnse2.com`
- **THEN** the system parses the composite ID, calls `DescribeSharedCNAME` to read the resource, and populates the state

### Requirement: Schema definition for shared CNAME resource
The resource schema SHALL include the following fields:
- `zone_id`: Required, String, ForceNew - the zone ID
- `shared_cname_prefix`: Required, String, ForceNew - the CNAME prefix (max 50 chars)
- `description`: Optional, String - description (1-50 chars)
- `shared_cname`: Computed, String - the full shared CNAME returned by API
- `ipssl_setting`: Optional, List(MaxItems:1) - IP SSL setting block containing:
  - `status`: Required, String - association status. Valid values: `bound` (IP SSL configuration bound), `binding` (IP SSL configuration binding), `unbinding` (IP SSL configuration unbinding), `unbound` (IP SSL configuration unbound)
  - `associated_domain`: Required, String - the domain associated with IP SSL. This field is empty when Status is `unbound`

#### Scenario: Schema validation
- **WHEN** user provides a valid Terraform configuration for `tencentcloud_teo_shared_cname`
- **THEN** the schema validates that `zone_id` and `shared_cname_prefix` are provided, and `ipssl_setting` block (if present) contains both `status` and `associated_domain`

### Requirement: Retry and error handling
All API calls SHALL be wrapped with `tccommon.ReadRetryTimeout` retry logic. Errors SHALL be wrapped with `tccommon.RetryError()`. The resource SHALL use `defer tccommon.LogElapsed()` for performance logging and `defer tccommon.InconsistentCheck()` for consistency checks.

#### Scenario: Transient API failure during read
- **WHEN** the `DescribeSharedCNAME` API returns a transient error
- **THEN** the system retries the request within the configured timeout period

#### Scenario: Non-retryable error during creation
- **WHEN** the `CreateSharedCNAME` API returns a non-retryable error
- **THEN** the system wraps the error with `tccommon.RetryError()` using `NonRetryableError` and returns immediately

### Requirement: Provider registration
The resource `tencentcloud_teo_shared_cname` SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_shared_cname` is available as a resource type for use in Terraform configurations
