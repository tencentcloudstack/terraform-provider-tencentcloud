# mongodb-audit-log-file Specification

## Purpose
TBD - created by archiving change add-mongodb-audit-log-file. Update Purpose after archive.
## Requirements
### Requirement: Resource schema definition
The `tencentcloud_mongodb_audit_log_file` resource SHALL define the following schema:
- `instance_id` (Required, ForceNew, String): MongoDB instance ID.
- `start_time` (Required, String): Start time for audit log, format "2021-07-12 10:29:20".
- `end_time` (Required, String): End time for audit log, format "2021-07-12 10:39:20".
- `order` (Optional, String): Sort order, values: "ASC" or "DESC".
- `order_by` (Optional, String): Sort field, values: "timestamp", "affectRows", "execTime".
- `filter` (Optional, List, MaxItems: 1): Filter conditions with sub-fields:
  - `host` (Optional, List of String): Client addresses.
  - `user` (Optional, List of String): Usernames.
  - `exec_time` (Optional, Int): Minimum execution time in ms.
  - `affect_rows` (Optional, Int): Minimum affected rows.
  - `atype` (Optional, List of String): Operation types.
  - `result` (Optional, List of String): Execution results.
  - `param` (Optional, List of String): Keywords to filter logs.
- `file_name` (Computed, String): The generated audit log file name.
- `items` (Computed, List): Audit log file details with sub-fields:
  - `file_name` (String): File name.
  - `create_time` (String): Creation time.
  - `status` (String): File status ("creating", "failed", "success").
  - `file_size` (Int): File size in KB.
  - `download_url` (String): Download URL.
  - `err_msg` (String): Error message.
  - `progress_rate` (Int): Download progress.

#### Scenario: Schema validates required fields
- **WHEN** a user defines a `tencentcloud_mongodb_audit_log_file` resource without `instance_id`, `start_time`, or `end_time`
- **THEN** Terraform SHALL return a validation error indicating the missing required fields

#### Scenario: Schema accepts optional filter
- **WHEN** a user defines the resource with a `filter` block containing `host` and `exec_time`
- **THEN** Terraform SHALL accept the configuration without error

### Requirement: Create audit log file
The resource SHALL call the `CreateAuditLogFile` API to create an audit log file. It MUST pass all specified input parameters. After successful creation, it SHALL store the composite ID (`instance_id` + separator + `file_name`) and call Read to populate computed attributes.

#### Scenario: Successful creation
- **WHEN** the `CreateAuditLogFile` API returns a non-empty `FileName`
- **THEN** the resource SHALL set the composite ID and invoke Read to populate state

#### Scenario: API returns empty FileName
- **WHEN** the `CreateAuditLogFile` API returns a nil or empty `FileName`
- **THEN** the resource SHALL return a non-retryable error indicating creation failed

### Requirement: Read audit log file
The resource SHALL call the `DescribeAuditLogFiles` API with `InstanceId` and `FileName` to retrieve the audit log file details. It SHALL set the `items` computed attribute with the returned file information.

#### Scenario: File exists
- **WHEN** `DescribeAuditLogFiles` returns items matching the file name
- **THEN** the resource SHALL populate the `items` attribute with file details

#### Scenario: File not found
- **WHEN** `DescribeAuditLogFiles` returns no items for the given file name
- **THEN** the resource SHALL remove the resource from Terraform state (set ID to "")

### Requirement: Delete audit log file
The resource SHALL call the `DeleteAuditLogFile` API with `InstanceId` and `FileName` extracted from the composite ID.

#### Scenario: Successful deletion
- **WHEN** the `DeleteAuditLogFile` API returns successfully
- **THEN** the resource SHALL be removed from Terraform state

### Requirement: Update returns error for immutable fields
Since no Update API exists, the resource Update function SHALL check if any non-ForceNew input fields changed and return an error indicating those fields are immutable.

#### Scenario: Attempt to modify immutable field
- **WHEN** a user changes `start_time`, `end_time`, `order`, `order_by`, or `filter` in the configuration
- **THEN** the resource SHALL return an error stating the field is immutable and cannot be updated

### Requirement: Import support with composite ID
The resource SHALL support import using the composite ID format `instance_id#file_name` (where `#` is `tccommon.FILED_SP`).

#### Scenario: Import existing audit log file
- **WHEN** a user runs `terraform import tencentcloud_mongodb_audit_log_file.example cmgo-xxxxx#audit_log_file_name`
- **THEN** the resource SHALL parse the composite ID, call Read, and populate the state

### Requirement: Retry with timeout
All API calls (Create, Read, Delete) SHALL be wrapped with retry logic using `tccommon.ReadRetryTimeout`. Errors SHALL be wrapped with `tccommon.RetryError`.

#### Scenario: Transient API failure during read
- **WHEN** `DescribeAuditLogFiles` returns a transient error
- **THEN** the resource SHALL retry the call until timeout is reached

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_mongodb_audit_log_file` SHALL be available as a resource type

