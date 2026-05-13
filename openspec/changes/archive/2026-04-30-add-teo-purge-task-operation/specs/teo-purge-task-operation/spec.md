## ADDED Requirements

### Requirement: Resource schema for tencentcloud_teo_purge_task
The resource SHALL define a Terraform schema with the following input fields:
- `zone_id` (Required, TypeString, ForceNew): Site ID
- `type` (Required, TypeString, ForceNew): Purge type, valid values: `purge_url`, `purge_prefix`, `purge_host`, `purge_all`, `purge_cache_tag`
- `method` (Optional, TypeString, ForceNew): Purge method, valid values: `invalidate`, `delete`. Default: `invalidate`
- `targets` (Optional, TypeList of TypeString, ForceNew): Resource list to purge
- `cache_tag` (Optional, TypeList, ForceNew, MaxItems: 1): Cache tag configuration block with `domains` (Optional, TypeList of TypeString) sub-field

The resource SHALL define the following computed output fields:
- `job_id` (Computed, TypeString): Task ID returned by CreatePurgeTask
- `tasks` (Computed, TypeList): List of purge task results with nested attributes: `job_id`, `target`, `type`, `method`, `status`, `create_time`, `update_time`, `fail_type`, `fail_message`

#### Scenario: Create purge URL task
- **WHEN** user creates a `tencentcloud_teo_purge_task` resource with `zone_id` set to a valid zone ID, `type` set to `purge_url`, and `targets` set to a list of URLs
- **THEN** the resource SHALL call `CreatePurgeTask` API with the provided parameters, poll `DescribePurgeTasks` until task status is terminal, and set `job_id` and `tasks` computed attributes

#### Scenario: Create purge all task
- **WHEN** user creates a `tencentcloud_teo_purge_task` resource with `zone_id` set to a valid zone ID, `type` set to `purge_all`
- **THEN** the resource SHALL call `CreatePurgeTask` API with `type=purge_all` and no targets, then poll for completion

#### Scenario: Create purge cache tag task
- **WHEN** user creates a `tencentcloud_teo_purge_task` resource with `zone_id`, `type` set to `purge_cache_tag`, and `cache_tag` block with `domains`
- **THEN** the resource SHALL call `CreatePurgeTask` API with the CacheTag parameter populated

### Requirement: Create function calls CreatePurgeTask API
The Create function SHALL call the `CreatePurgeTask` API from `teo/v20220901` SDK package. It SHALL:
1. Construct `CreatePurgeTaskRequest` with `ZoneId`, `Type`, `Method` (if set), `Targets` (if set), and `CacheTag` (if set)
2. Call the API with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and handle errors with `tccommon.RetryError()`
3. Check the response: if `JobId` is empty, return `tccommon.NonRetryableError`
4. Set `d.SetId()` with `helper.BuildToken()` after successful API call

#### Scenario: API call succeeds
- **WHEN** `CreatePurgeTask` API returns a valid `JobId`
- **THEN** the resource SHALL set the Terraform ID and proceed to poll for task status

#### Scenario: API call returns empty JobId
- **WHEN** `CreatePurgeTask` API returns an empty `JobId`
- **THEN** the resource SHALL return a non-retryable error

#### Scenario: API call fails
- **WHEN** `CreatePurgeTask` API returns an error
- **THEN** the resource SHALL retry using `tccommon.RetryError()` wrapping, and ultimately return the error if all retries fail

### Requirement: Poll DescribePurgeTasks for task completion
After `CreatePurgeTask` succeeds, the Create function SHALL poll `DescribePurgeTasks` API using the returned `JobId` to verify task completion. It SHALL:
1. Call `DescribePurgeTasks` with `ZoneId` and filter by `job-id`
2. Wait until the task status is a terminal state: `success`, `failed`, `timeout`, or `canceled`
3. Use `tccommon.ReadRetryTimeout` for the polling retry loop
4. Set the `tasks` computed attribute from the response

#### Scenario: Task completes successfully
- **WHEN** `DescribePurgeTasks` returns task status as `success`
- **THEN** the resource SHALL set the `tasks` computed attribute and complete creation

#### Scenario: Task fails
- **WHEN** `DescribePurgeTasks` returns task status as `failed`
- **THEN** the resource SHALL return an error with the failure information

#### Scenario: Task times out
- **WHEN** `DescribePurgeTasks` returns task status as `timeout`
- **THEN** the resource SHALL return an error indicating the task timed out

### Requirement: Read function is empty
The Read function SHALL return nil without performing any API calls. This is consistent with the OPERATION resource pattern where no persistent state needs to be refreshed.

#### Scenario: Terraform refresh
- **WHEN** `terraform refresh` is executed
- **THEN** the Read function SHALL return nil without calling any API

### Requirement: Delete function is empty
The Delete function SHALL return nil without performing any API calls. Purge tasks cannot be undone.

#### Scenario: Terraform destroy
- **WHEN** `terraform destroy` is executed on the resource
- **THEN** the Delete function SHALL return nil without calling any API

### Requirement: No Update method
The resource SHALL NOT define an Update method. Operation resources are immutable - any change to input parameters triggers a destroy and recreate.

#### Scenario: Parameter change
- **WHEN** user modifies any input parameter
- **THEN** Terraform SHALL force recreation of the resource (all fields are ForceNew)

### Requirement: Resource registration in provider
The resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap with key `tencentcloud_teo_purge_task` and in `tencentcloud/provider.md` documentation index.

#### Scenario: Provider loads resource
- **WHEN** Terraform provider is initialized
- **THEN** `tencentcloud_teo_purge_task` SHALL be available as a resource type

### Requirement: Unit tests with gomonkey mocks
The resource SHALL have unit tests in `resource_tc_teo_purge_task_operation_test.go` using gomonkey for mocking cloud API calls. Tests SHALL cover:
1. Successful creation with `purge_url` type
2. Successful creation with `purge_cache_tag` type
3. API error handling

Tests SHALL be runnable with `go test -gcflags=all=-l` without requiring real cloud credentials.

#### Scenario: Unit test for successful purge_url creation
- **WHEN** unit test runs with mocked `CreatePurgeTask` returning valid JobId and mocked `DescribePurgeTasks` returning `success` status
- **THEN** test SHALL pass with correct resource attributes set

#### Scenario: Unit test for API error
- **WHEN** unit test runs with mocked `CreatePurgeTask` returning an error
- **THEN** test SHALL verify the error is properly returned
