## ADDED Requirements

### Requirement: DTS sync job config resource exists
The Terraform Provider SHALL provide a `tencentcloud_dts_sync_job_config` resource of kind RESOURCE_KIND_CONFIG that manages the dynamic configuration of a Tencent Cloud DTS sync job. The resource id SHALL be the sync job id (`job_id`).

#### Scenario: Resource is created from an existing sync job
- **WHEN** a user applies a `tencentcloud_dts_sync_job_config` resource with `job_id` set to an existing sync job id
- **THEN** the provider SHALL set the terraform state id to `job_id` and call `ModifySyncJobConfig` to apply the dynamic configuration, then call Read to refresh state

#### Scenario: Resource is imported
- **WHEN** a user imports `tencentcloud_dts_sync_job_config` with an existing job id
- **THEN** the provider SHALL populate state by calling `DescribeSyncJobs` and return the current dynamic configuration

### Requirement: Read sync job config via DescribeSyncJobs
The resource Read function SHALL call the cloud API `DescribeSyncJobs` with `request.JobId` set to `d.Id()` and SHALL read the configuration from `response.Response.JobList[0]` (SyncJobInfo). When the response or instance is empty, the provider SHALL log the id via `log.Printf("[CRUD] dts_sync_job_config id=%s", d.Id())` before calling `d.SetId("")`.

#### Scenario: Sync job exists
- **WHEN** Read is called for an existing sync job id
- **THEN** the provider SHALL set `job_id`, `dynamic_objects`, `dynamic_options`, and the computed output fields from the returned `SyncJobInfo`

#### Scenario: Sync job not found
- **WHEN** `DescribeSyncJobs` returns an empty `JobList`
- **THEN** the provider SHALL log the id and call `d.SetId("")` so terraform removes the resource from state

### Requirement: Update sync job config via ModifySyncJobConfig
The resource Update function SHALL call the cloud API `ModifySyncJobConfig` with `request.JobId` (required), and optionally `request.DynamicObjects` and `request.DynamicOptions` when the corresponding schema blocks are configured. The API call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and errors SHALL be wrapped with `tccommon.RetryError(e)`.

#### Scenario: Update dynamic objects and options
- **WHEN** a user changes `dynamic_objects` or `dynamic_options` in the configuration
- **THEN** the provider SHALL build a `ModifySyncJobConfigRequest` from the schema and call `ModifySyncJobConfig`; on success the provider SHALL call Read to refresh state

#### Scenario: API call fails transiently
- **WHEN** `ModifySyncJobConfig` returns a retryable error
- **THEN** the provider SHALL retry within `tccommon.WriteRetryTimeout` via `tccommon.RetryError`

### Requirement: Schema covers ModifySyncJobConfig inputs
The resource schema SHALL expose `job_id` (Required, ForceNew), `dynamic_objects` (Optional, single-block), and `dynamic_options` (Optional, single-block). The `dynamic_objects` block SHALL mirror the SDK `Objects` struct (mode, databases, advanced_objects, online_ddl, databases_op_filter). The `dynamic_options` block SHALL mirror the SDK `DynamicOptions` struct (op_types, ddl_options, conflict_handle_type, conflict_handle_option, kafka_option, filter_begin_commit, filter_checkpoint, deal_of_exist_same_table, start_position).

#### Scenario: Required fields are enforced
- **WHEN** a user omits `job_id`
- **THEN** terraform SHALL report a validation error before any API call

#### Scenario: Nested required fields are enforced
- **WHEN** `dynamic_objects.online_ddl` is configured without `status`, or `dynamic_options` is configured without `op_types` and `ddl_options`
- **THEN** terraform SHALL report a validation error

### Requirement: Schema exposes DescribeSyncJobs outputs as computed fields
The resource schema SHALL expose the `DescribeSyncJobs` response fields as computed attributes so that `terraform plan` reflects drift. These include job metadata (status, run_mode, pay_mode, etc.), src/dst endpoint info, options, objects, and detail fields as listed in the proposal's output mapping.

#### Scenario: Read refreshes computed output fields
- **WHEN** Read succeeds for an existing job
- **THEN** the provider SHALL set every computed output field from the `SyncJobInfo` response whose value is non-nil

### Requirement: Provider registration
The provider SHALL register `tencentcloud_dts_sync_job_config` in `tencentcloud/provider.go` and `tencentcloud/provider.md`.

#### Scenario: Resource is usable after build
- **WHEN** the provider is built with the new resource registered
- **THEN** `terraform plan` SHALL accept `tencentcloud_dts_sync_job_config` as a valid resource type

### Requirement: Unit tests cover business logic
The resource SHALL have a unit test file `resource_tc_tencentcloud_dts_sync_job_config_test.go` that uses gomonkey to mock the cloud API calls (DescribeSyncJobs and ModifySyncJobConfig) and validates the Create/Read/Update business logic without invoking real TencentCloud endpoints.

#### Scenario: Unit tests mock cloud APIs
- **WHEN** the unit tests run
- **THEN** they SHALL mock `DescribeSyncJobs` and `ModifySyncJobConfig` via gomonkey and assert the resource CRUD behavior