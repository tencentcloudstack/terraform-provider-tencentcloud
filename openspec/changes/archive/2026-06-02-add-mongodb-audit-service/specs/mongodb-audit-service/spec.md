## ADDED Requirements

### Requirement: Create MongoDB audit service

The system SHALL allow users to create a `tencentcloud_mongodb_audit_service` resource that opens audit service on a MongoDB instance by calling the `OpenAuditService` API.

The resource SHALL accept the following parameters:
- `instance_id` (Required, ForceNew, String): The MongoDB instance ID
- `log_expire_day` (Required, Int): Audit log retention days, valid values: 7, 30, 90, 180, 365, 1095, 1825
- `audit_all` (Required, Bool): true for full audit, false for rule-based audit
- `rule_filters` (Optional, List): Audit filter rules, only needed when `audit_all` is false
  - `type` (Required, String): Filter condition name, valid values: SrcIp, DB, Collection, User, SqlType
  - `compare` (Required, String): Filter match type, MUST be "EQ"
  - `value` (Required, List of String): Filter match values

After calling `OpenAuditService`, the system SHALL poll `DescribeAuditConfig` until `IsOpening` equals `"false"` to confirm the operation is complete.

The resource ID SHALL be set to the `instance_id` value.

#### Scenario: Create audit service with full audit mode
- **WHEN** user creates a `tencentcloud_mongodb_audit_service` resource with `audit_all = true` and `log_expire_day = 30`
- **THEN** the system calls `OpenAuditService` with the specified parameters and polls until the audit service is fully opened

#### Scenario: Create audit service with rule-based audit
- **WHEN** user creates a `tencentcloud_mongodb_audit_service` resource with `audit_all = false` and `rule_filters` specified
- **THEN** the system calls `OpenAuditService` with the specified parameters including rule filters and polls until the audit service is fully opened

#### Scenario: Create audit service returns empty response
- **WHEN** the `OpenAuditService` API returns a nil response or empty response
- **THEN** the system SHALL return a non-retryable error

### Requirement: Read MongoDB audit service configuration

The system SHALL read the audit service configuration by calling `DescribeAuditConfig` API with the `instance_id`.

The system SHALL set the following computed attributes from the response:
- `instance_name` (Computed, String): Instance name
- `create_time` (Computed, String): Time when audit was enabled
- `log_type` (Computed, String): Audit log storage type
- `is_closing` (Computed, String): Whether audit is being closed
- `is_opening` (Computed, String): Whether audit is being opened

The system SHALL also refresh `audit_all` and `log_expire_day` from the response.

If the API returns an error indicating the resource does not exist, the system SHALL remove the resource from state by calling `d.SetId("")`.

#### Scenario: Read existing audit service configuration
- **WHEN** the resource exists and `DescribeAuditConfig` returns valid data
- **THEN** the system sets all computed and configurable attributes in state

#### Scenario: Read non-existent audit service
- **WHEN** `DescribeAuditConfig` returns an error indicating the audit service is not enabled
- **THEN** the system removes the resource from state

### Requirement: Update MongoDB audit service configuration

The system SHALL allow users to update `log_expire_day`, `audit_all`, and `rule_filters` by calling the `ModifyAuditService` API.

The `instance_id` field SHALL NOT be updatable (ForceNew).

#### Scenario: Update log expire day
- **WHEN** user changes `log_expire_day` from 30 to 90
- **THEN** the system calls `ModifyAuditService` with the new value and reads back the updated configuration

#### Scenario: Switch from full audit to rule-based audit
- **WHEN** user changes `audit_all` from true to false and adds `rule_filters`
- **THEN** the system calls `ModifyAuditService` with the updated audit mode and rule filters

### Requirement: Delete MongoDB audit service

The system SHALL delete the audit service by calling `CloseAuditService` API with the `instance_id`.

After calling `CloseAuditService`, the system SHALL poll `DescribeAuditConfig` until `IsClosing` equals `"false"` to confirm the operation is complete.

#### Scenario: Delete audit service successfully
- **WHEN** user destroys the `tencentcloud_mongodb_audit_service` resource
- **THEN** the system calls `CloseAuditService` and polls until the audit service is fully closed

### Requirement: Import MongoDB audit service

The system SHALL support importing an existing audit service configuration using the `instance_id` as the import ID.

#### Scenario: Import existing audit service
- **WHEN** user runs `terraform import tencentcloud_mongodb_audit_service.example cmgo-xxxxxx`
- **THEN** the system reads the audit configuration and populates the state

### Requirement: Resource registration

The system SHALL register `tencentcloud_mongodb_audit_service` in `provider.go` and `provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_mongodb_audit_service` is available as a resource type

### Requirement: Resource documentation

The system SHALL provide a `.md` documentation file for the resource with:
- One-line description mentioning MongoDB product name
- Example Usage section showing both full audit and rule-based audit configurations
- Import section showing the import command

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** a documentation file exists at `tencentcloud/services/mongodb/resource_tc_mongodb_audit_service.md`

### Requirement: Unit tests with gomonkey mock

The system SHALL provide unit tests using gomonkey to mock cloud API calls, verifying:
- Create flow (OpenAuditService + polling)
- Read flow (DescribeAuditConfig)
- Update flow (ModifyAuditService)
- Delete flow (CloseAuditService + polling)

Tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Unit tests pass
- **WHEN** `go test -gcflags=all=-l` is run on the test file
- **THEN** all tests pass successfully
