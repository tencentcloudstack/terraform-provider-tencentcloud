# postgres-audit-log-file Specification

## Purpose
TBD - created by archiving change add-postgres-audit-log-file-attachment. Update Purpose after archive.
## Requirements
### Requirement: Create audit log file
The system SHALL provide a Terraform resource `tencentcloud_postgres_audit_log_file` that creates a PostgreSQL audit log file by calling the `CreateAuditLogFile` API with the specified instance ID, time range, product name, and optional filter conditions.

#### Scenario: Successful creation with all parameters
- **WHEN** user applies a Terraform configuration with `tencentcloud_postgres_audit_log_file` specifying `instance_id`, `start_time`, `end_time`, `product`, and `filter` block
- **THEN** the system SHALL call `CreateAuditLogFile` API with all provided parameters, poll `DescribeAuditLogFiles` until the file status becomes `success`, and set the resource ID as `instance_id#file_name`

#### Scenario: Successful creation without filter
- **WHEN** user applies a Terraform configuration with `tencentcloud_postgres_audit_log_file` specifying only `instance_id`, `start_time`, `end_time`, and `product` (no filter)
- **THEN** the system SHALL call `CreateAuditLogFile` API without the Filter parameter and poll until completion

#### Scenario: Creation fails asynchronously
- **WHEN** the audit log file creation results in a `failed` status during polling
- **THEN** the system SHALL return an error containing the ErrMsg from the API response

### Requirement: Read audit log file
The system SHALL support reading the audit log file metadata by calling `DescribeAuditLogFiles` with the `InstanceId`, `Product`, and `FileName` parameters extracted from the resource ID.

#### Scenario: File exists and is readable
- **WHEN** Terraform performs a refresh on an existing `tencentcloud_postgres_audit_log_file` resource
- **THEN** the system SHALL call `DescribeAuditLogFiles` with the FileName filter and populate computed attributes (`file_name`, `status`, `file_size`, `create_time`, `download_url`, `err_msg`, `progress`, `finish_time`) from the response Items

#### Scenario: File no longer exists
- **WHEN** Terraform performs a refresh and `DescribeAuditLogFiles` returns no matching items
- **THEN** the system SHALL remove the resource from state (call `d.SetId("")`)

### Requirement: Delete audit log file
The system SHALL support deleting the audit log file by calling `DeleteAuditLogFile` with `InstanceId`, `Product`, and `FileName` extracted from the resource ID.

#### Scenario: Successful deletion
- **WHEN** user destroys a `tencentcloud_postgres_audit_log_file` resource
- **THEN** the system SHALL call `DeleteAuditLogFile` API with the correct parameters and remove the resource from state

### Requirement: Resource schema definition
The system SHALL define the resource schema with the following fields:

#### Scenario: Required input fields
- **WHEN** the resource schema is defined
- **THEN** the following fields SHALL be Required and ForceNew: `instance_id` (String), `start_time` (String), `end_time` (String), `product` (String)

#### Scenario: Optional input fields
- **WHEN** the resource schema is defined
- **THEN** the `filter` field SHALL be Optional and ForceNew, containing a nested block with fields: `affect_rows` (Int, Optional), `db_name` (List of String, Optional), `exec_time` (Int, Optional), `host` (List of String, Optional), `sql` (String, Optional), `user` (List of String, Optional), `sql_type` (List of String, Optional)

#### Scenario: Computed output fields
- **WHEN** the resource schema is defined
- **THEN** the following fields SHALL be Computed: `file_name` (String), `status` (String), `file_size` (Int), `create_time` (String), `download_url` (String), `err_msg` (String), `progress` (Int), `finish_time` (String)

### Requirement: Import support
The system SHALL support importing existing audit log files using the composite ID format `instance_id#file_name`.

#### Scenario: Import by composite ID
- **WHEN** user runs `terraform import tencentcloud_postgres_audit_log_file.example instance_id#file_name`
- **THEN** the system SHALL parse the composite ID, call `DescribeAuditLogFiles` to read the file metadata, and populate the state

### Requirement: Retry and error handling
The system SHALL wrap all API calls with retry logic using `tccommon.ReadRetryTimeout` and `tccommon.RetryError()`.

#### Scenario: Transient API failure during creation
- **WHEN** `CreateAuditLogFile` API returns a transient error
- **THEN** the system SHALL retry the call within the configured timeout

#### Scenario: Transient API failure during read
- **WHEN** `DescribeAuditLogFiles` API returns a transient error
- **THEN** the system SHALL retry the call within the configured timeout

#### Scenario: Transient API failure during deletion
- **WHEN** `DeleteAuditLogFile` API returns a transient error
- **THEN** the system SHALL retry the call within the configured timeout

