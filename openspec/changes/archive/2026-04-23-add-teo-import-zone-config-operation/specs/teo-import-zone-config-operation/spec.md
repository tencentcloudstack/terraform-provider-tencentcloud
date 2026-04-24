## ADDED Requirements

### Requirement: Import zone config operation resource
The system SHALL provide a Terraform resource `tencentcloud_teo_import_zone_config` of type RESOURCE_KIND_OPERATION that calls the `ImportZoneConfig` cloud API to import TEO site configuration. The resource SHALL support only Create operation; Read, Update, and Delete SHALL be no-ops.

#### Scenario: Successful zone config import
- **WHEN** user creates a `tencentcloud_teo_import_zone_config` resource with valid `zone_id` and `content`
- **THEN** the system SHALL call `ImportZoneConfig` API with the provided parameters, poll `DescribeZoneConfigImportResult` until the task status is "success", and set computed attributes (`task_id`, `status`, `message`, `import_time`, `finish_time`)

#### Scenario: Zone config import failure
- **WHEN** user creates a `tencentcloud_teo_import_zone_config` resource and the async task status becomes "failure"
- **THEN** the system SHALL return an error containing the `message` from the import result

#### Scenario: Zone config import polling timeout
- **WHEN** the import task does not reach a terminal state (success/failure) within the polling timeout
- **THEN** the system SHALL return a retryable error indicating the operation is still processing

### Requirement: Schema definition for import zone config
The resource SHALL define the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): TEO zone ID
- `content` (Required, ForceNew, TypeString): Configuration content in JSON format to import
- `task_id` (Computed, TypeString): Async task ID returned by ImportZoneConfig
- `status` (Computed, TypeString): Import task status (success/failure/doing)
- `message` (Computed, TypeString): Status message from import result
- `import_time` (Computed, TypeString): Import start time
- `finish_time` (Computed, TypeString): Import finish time

#### Scenario: Required fields validation
- **WHEN** user creates the resource without `zone_id` or `content`
- **THEN** Terraform SHALL produce a validation error indicating the missing required field

#### Scenario: All fields are ForceNew
- **WHEN** user modifies any schema field after creation
- **THEN** Terraform SHALL trigger a resource recreation (destroy and create)

### Requirement: Async task polling
After calling `ImportZoneConfig`, the system SHALL poll `DescribeZoneConfigImportResult` using `resource.Retry` with a timeout of `6 * tccommon.ReadRetryTimeout`. The polling SHALL use `zone_id` and `task_id` as query parameters.

#### Scenario: Polling with doing status
- **WHEN** the DescribeZoneConfigImportResult returns status "doing"
- **THEN** the system SHALL continue polling with a retryable error

#### Scenario: Polling with success status
- **WHEN** the DescribeZoneConfigImportResult returns status "success"
- **THEN** the system SHALL stop polling and set all computed attributes

#### Scenario: Polling with failure status
- **WHEN** the DescribeZoneConfigImportResult returns status "failure"
- **THEN** the system SHALL return a non-retryable error with the message from the response

### Requirement: Resource registration
The resource SHALL be registered in `tencentcloud/provider.go` with key `tencentcloud_teo_import_zone_config` and the corresponding provider documentation SHALL be updated in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_import_zone_config` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with mock
The resource SHALL have unit tests using gomonkey to mock the cloud API calls, testing the Create logic including async polling behavior.

#### Scenario: Unit test for successful import
- **WHEN** running unit tests
- **THEN** the Create function SHALL be tested with mocked ImportZoneConfig returning a TaskId and mocked DescribeZoneConfigImportResult returning success status

#### Scenario: Unit test for failed import
- **WHEN** running unit tests
- **THEN** the Create function SHALL be tested with mocked DescribeZoneConfigImportResult returning failure status, verifying error is returned
