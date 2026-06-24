## ADDED Requirements

### Requirement: Enable PostgreSQL audit service
The system SHALL allow users to enable database audit service on a PostgreSQL instance by specifying instance_id, log_expire_day, hot_log_expire_day, and audit_type. The resource SHALL call the OpenAuditService API to enable audit.

#### Scenario: Successfully enable audit service
- **WHEN** user creates a `tencentcloud_postgres_audit_service` resource with valid instance_id, log_expire_day=30, hot_log_expire_day=7, and audit_type="simple"
- **THEN** the system calls OpenAuditService API with the provided parameters and sets the resource ID to the instance_id

#### Scenario: Enable audit service with all parameters
- **WHEN** user creates a `tencentcloud_postgres_audit_service` resource with instance_id, log_expire_day=365, hot_log_expire_day=90, audit_type="complex", and product="postgres"
- **THEN** the system calls OpenAuditService API with all parameters including Product field

#### Scenario: Create fails with empty response
- **WHEN** the OpenAuditService API returns a nil response or empty result
- **THEN** the system SHALL return a non-retryable error indicating the API response is invalid

### Requirement: Read PostgreSQL audit service status
The system SHALL read the current audit service configuration by calling DescribeAuditInstanceList API with InstanceId filter and AuditSwitch=1.

#### Scenario: Successfully read audit service configuration
- **WHEN** the system reads the resource state for an existing audit service
- **THEN** the system calls DescribeAuditInstanceList with Product="postgres", AuditSwitch=1, and Filters containing the InstanceId, and sets log_expire_day, hot_log_expire_day, and computed attributes from the AuditInstanceInfo response

#### Scenario: Audit service not found (deleted externally)
- **WHEN** the system reads the resource state but DescribeAuditInstanceList returns no matching items or the instance has AuditStatus="OFF"
- **THEN** the system SHALL remove the resource from state (d.SetId(""))

### Requirement: Modify PostgreSQL audit service configuration
The system SHALL allow users to modify log_expire_day, hot_log_expire_day, and audit_type on an existing audit service by calling ModifyAuditService API.

#### Scenario: Successfully modify audit configuration
- **WHEN** user updates log_expire_day, hot_log_expire_day, or audit_type on an existing resource
- **THEN** the system calls ModifyAuditService API with the updated parameters and the existing instance_id

#### Scenario: instance_id change forces recreation
- **WHEN** user changes the instance_id attribute
- **THEN** the system SHALL destroy the existing resource and create a new one (ForceNew behavior)

### Requirement: Disable PostgreSQL audit service
The system SHALL disable the audit service when the resource is destroyed by calling CloseAuditService API.

#### Scenario: Successfully disable audit service
- **WHEN** user destroys the `tencentcloud_postgres_audit_service` resource
- **THEN** the system calls CloseAuditService API with the instance_id and product parameters

### Requirement: Import existing audit service
The system SHALL support importing an existing audit service configuration by instance_id.

#### Scenario: Successfully import audit service
- **WHEN** user runs `terraform import tencentcloud_postgres_audit_service.example <instance_id>`
- **THEN** the system reads the audit service configuration and populates the state with all attributes

### Requirement: Retry on API errors
All API calls (OpenAuditService, DescribeAuditInstanceList, ModifyAuditService, CloseAuditService) SHALL be wrapped with retry logic using tccommon.ReadRetryTimeout.

#### Scenario: Transient API error during create
- **WHEN** OpenAuditService returns a transient error
- **THEN** the system SHALL retry the request within the configured timeout period

#### Scenario: Transient API error during read
- **WHEN** DescribeAuditInstanceList returns a transient error
- **THEN** the system SHALL retry the request within the configured timeout period
