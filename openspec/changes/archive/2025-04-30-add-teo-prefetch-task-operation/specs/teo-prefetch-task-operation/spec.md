## ADDED Requirements

### Requirement: Create TEO prefetch task resource
The system SHALL provide a `tencentcloud_teo_prefetch_task` Terraform resource of type RESOURCE_KIND_OPERATION that allows users to submit URL prefetch tasks to TEO EdgeOne. The resource SHALL call CreatePrefetchTask API and poll DescribePrefetchTasks API until the task status is no longer "processing".

#### Scenario: Successful prefetch task creation
- **WHEN** user creates a `tencentcloud_teo_prefetch_task` resource with required `zone_id` and `targets` parameters
- **THEN** the system SHALL call CreatePrefetchTask API with the provided parameters, receive a JobId, poll DescribePrefetchTasks using the job-id filter until the task status is not "processing", and set the resource ID to `zoneId:jobId`

#### Scenario: Prefetch task creation with all optional parameters
- **WHEN** user creates a `tencentcloud_teo_prefetch_task` resource with `zone_id`, `targets`, `mode`, `headers`, and `prefetch_media_segments` parameters
- **THEN** the system SHALL call CreatePrefetchTask API with all provided parameters and poll until completion

#### Scenario: Prefetch task fails
- **WHEN** the polled task status from DescribePrefetchTasks is "failed"
- **THEN** the system SHALL return a NonRetryableError with the fail type and fail message

#### Scenario: Prefetch task times out during polling
- **WHEN** the task status remains "processing" beyond the polling timeout
- **THEN** the system SHALL return a retryable error

### Requirement: Resource schema definition
The system SHALL define the following schema for `tencentcloud_teo_prefetch_task`:

Required parameters (ForceNew):
- `zone_id` (TypeString): Zone ID
- `targets` (TypeList, element TypeString): List of URLs to prefetch

Optional parameters (ForceNew):
- `mode` (TypeString): Prefetch mode, valid values: default, edge
- `headers` (TypeList, element TypeMap): HTTP headers to carry during prefetch, each element has `name` and `value` keys
- `prefetch_media_segments` (TypeString): Media segment prefetch control, valid values: on, off

Computed parameters:
- `job_id` (TypeString): Task job ID returned by CreatePrefetchTask
- `tasks` (TypeList): Task result list, each element contains: job_id, target, type, method, status, create_time, update_time, fail_type, fail_message

#### Scenario: Schema validates required fields
- **WHEN** user omits `zone_id` or `targets` in the resource configuration
- **THEN** Terraform SHALL produce a validation error indicating the missing required field

#### Scenario: All schema fields are ForceNew
- **WHEN** user changes any input parameter after creation
- **THEN** Terraform SHALL force resource recreation

### Requirement: Empty Read/Update/Delete methods
The system SHALL implement empty Read, Update, and Delete methods for `tencentcloud_teo_prefetch_task` resource, as this is a RESOURCE_KIND_OPERATION type.

#### Scenario: Read method returns nil
- **WHEN** Terraform calls the Read method
- **THEN** the system SHALL return nil without making any API calls

#### Scenario: Delete method returns nil
- **WHEN** Terraform calls the Delete method
- **THEN** the system SHALL return nil without making any API calls

### Requirement: Resource registration in provider
The system SHALL register `tencentcloud_teo_prefetch_task` resource in `tencentcloud/provider.go` with the resource name "tencentcloud_teo_prefetch_task" and add corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** user references `tencentcloud_teo_prefetch_task` in Terraform configuration
- **THEN** the provider SHALL recognize and process the resource

### Requirement: Unit tests with gomonkey mock
The system SHALL include unit tests in `resource_tc_teo_prefetch_task_operation_test.go` using gomonkey to mock the TEO API client calls. Tests SHALL cover the Create operation with successful task completion and failed task scenarios.

#### Scenario: Test successful prefetch task creation
- **WHEN** unit test runs with mocked CreatePrefetchTask returning a JobId and mocked DescribePrefetchTasks returning status "success"
- **THEN** the test SHALL verify the resource ID is set correctly and computed fields are populated

#### Scenario: Test failed prefetch task
- **WHEN** unit test runs with mocked CreatePrefetchTask returning a JobId and mocked DescribePrefetchTasks returning status "failed"
- **THEN** the test SHALL verify that an error is returned

### Requirement: Resource documentation
The system SHALL include a `resource_tc_teo_prefetch_task_operation.md` documentation file with a description, example usage, and import section (not applicable for OPERATION type).

#### Scenario: Documentation file exists
- **WHEN** the resource is added
- **THEN** a corresponding .md file SHALL exist in the teo service directory with proper format following the gendoc README guidelines
