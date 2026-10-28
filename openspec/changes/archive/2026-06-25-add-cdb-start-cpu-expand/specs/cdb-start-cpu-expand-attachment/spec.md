## ADDED Requirements

### Requirement: Resource schema defines CPU elastic expansion attachment parameters
The resource `tencentcloud_cdb_start_cpu_expand` SHALL define a schema with the following top-level fields:
- `instance_id` (Required, ForceNew, TypeString): CDB instance ID
- `type` (Required, ForceNew, TypeString): Expansion type — one of `auto`, `manual`, `timeInterval`, `period`
- `expand_cpu` (Optional, ForceNew, TypeInt): CPU cores to expand (required when type is `manual`, `timeInterval`, or `period`)
- `auto_strategy` (Optional, ForceNew, TypeList, MaxItems:1): Auto expansion strategy block (required when type is `auto`)
- `time_interval_strategy` (Optional, ForceNew, TypeList, MaxItems:1): Time interval expansion strategy block (required when type is `timeInterval`)
- `period_strategy` (Optional, ForceNew, TypeList, MaxItems:1): Period expansion strategy block (required when type is `period`)
- `async_request_id` (Computed, TypeString): Async request ID returned by Create/Delete APIs

#### Scenario: Schema validation for auto type
- **WHEN** user sets `type` to `auto` and provides `auto_strategy` block
- **THEN** the resource SHALL accept the configuration and call `StartCpuExpand` with `AutoStrategy` parameter

#### Scenario: Schema validation for manual type
- **WHEN** user sets `type` to `manual` and provides `expand_cpu` value
- **THEN** the resource SHALL accept the configuration and call `StartCpuExpand` with `ExpandCpu` and `Type=manual` parameters

#### Scenario: Schema validation for timeInterval type
- **WHEN** user sets `type` to `timeInterval` and provides `time_interval_strategy` block with `start_time` and `end_time`
- **THEN** the resource SHALL accept the configuration and call `StartCpuExpand` with `TimeIntervalStrategy` parameter

#### Scenario: Schema validation for period type
- **WHEN** user sets `type` to `period` and provides `period_strategy` block
- **THEN** the resource SHALL accept the configuration and call `StartCpuExpand` with `PeriodStrategy` parameter

### Requirement: Auto strategy block schema
The `auto_strategy` block SHALL contain the following fields:
- `expand_threshold` (Required, TypeInt): Auto expansion threshold (40, 50, 60, 70, 80, 90)
- `shrink_threshold` (Required, TypeInt): Auto shrink threshold (10, 20, 30)
- `expand_second_period` (Optional, TypeInt): Expansion observation period in seconds (15, 30, 45, 60, 180, 300, 600, 900, 1800)
- `shrink_second_period` (Optional, TypeInt): Shrink observation period in seconds (300, 600, 900, 1800)

#### Scenario: Auto strategy with all parameters
- **WHEN** user provides `auto_strategy` with `expand_threshold`, `shrink_threshold`, `expand_second_period`, and `shrink_second_period`
- **THEN** the resource SHALL map all fields to the `AutoStrategy` struct in the `StartCpuExpand` request

### Requirement: Time interval strategy block schema
The `time_interval_strategy` block SHALL contain the following fields:
- `start_time` (Required, TypeInt): Start expansion time as integer timestamp (seconds)
- `end_time` (Required, TypeInt): End expansion time as integer timestamp (seconds)

#### Scenario: Time interval strategy configuration
- **WHEN** user provides `time_interval_strategy` with `start_time` and `end_time`
- **THEN** the resource SHALL map the fields to `TimeIntervalStrategy` struct in the `StartCpuExpand` request

### Requirement: Period strategy block schema
The `period_strategy` block SHALL contain the following nested blocks:
- `time_cycle` (Optional, TypeList, MaxItems:1): Weekly cycle configuration with fields: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday` (all TypeBool)
- `time_interval` (Optional, TypeList, MaxItems:1): Daily time range with fields: `start_time` (TypeString), `end_time` (TypeString)

#### Scenario: Period strategy with weekly cycle and time interval
- **WHEN** user provides `period_strategy` with `time_cycle` (selecting specific days) and `time_interval` (specifying start/end time strings)
- **THEN** the resource SHALL map the fields to `PeriodStrategy` struct containing `TImeCycle` and `TimeInterval` sub-structs

### Requirement: Resource Create operation
The resource Create operation SHALL call the `StartCpuExpand` cloud API with all configured parameters. Since `StartCpuExpand` is an async interface that returns `AsyncRequestId`, the Create operation SHALL:
1. Call `StartCpuExpand` with retry using `tccommon.WriteRetryTimeout`
2. Validate that the response is not nil and `AsyncRequestId` is not empty; if empty, return `NonRetryableError`
3. Set the resource ID to the `instance_id` value
4. Poll the `DescribeCPUExpandStrategyInfo` API until the expansion strategy is confirmed as active (the `Type` field is not nil/null)

#### Scenario: Successful creation with async polling
- **WHEN** `StartCpuExpand` API succeeds and returns a valid `AsyncRequestId`
- **THEN** the resource SHALL set the ID to `instance_id` and poll `DescribeCPUExpandStrategyInfo` until the expansion strategy is confirmed

#### Scenario: StartCpuExpand returns nil response
- **WHEN** `StartCpuExpand` API returns nil response or nil `AsyncRequestId`
- **THEN** the resource SHALL return `NonRetryableError` to prevent writing an empty ID

### Requirement: Resource Read operation
The resource Read operation SHALL call `DescribeCPUExpandStrategyInfo` with the `instance_id` extracted from `d.Id()`. The Read operation SHALL:
1. Use retry with `tccommon.ReadRetryTimeout`
2. If the response is nil or `Type` field is nil/null, log the empty response with `log.Printf("[CRUD] cdb_start_cpu_expand id=%s", d.Id())` before calling `d.SetId("")`
3. If the strategy is active, set all fields from the response: `type`, `expand_cpu`, `auto_strategy`, `time_interval_strategy`, `period_strategy`

#### Scenario: Reading an active expansion strategy
- **WHEN** `DescribeCPUExpandStrategyInfo` returns a valid strategy with `Type` field present
- **THEN** the resource SHALL set all corresponding schema fields from the response

#### Scenario: Expansion strategy no longer exists
- **WHEN** `DescribeCPUExpandStrategyInfo` returns nil `Type` (expansion not enabled)
- **THEN** the resource SHALL log the ID and call `d.SetId("")` to mark the resource as gone

### Requirement: Resource Delete operation
The resource Delete operation SHALL call the `StopCpuExpand` cloud API with `instance_id` extracted from `d.Id()`. Since `StopCpuExpand` is an async interface, the Delete operation SHALL:
1. Call `StopCpuExpand` with retry using `tccommon.WriteRetryTimeout`
2. Poll `DescribeCPUExpandStrategyInfo` until the expansion strategy is confirmed as removed (the `Type` field is nil/null)

#### Scenario: Successful deletion with async polling
- **WHEN** `StopCpuExpand` API succeeds and returns a valid `AsyncRequestId`
- **THEN** the resource SHALL poll `DescribeCPUExpandStrategyInfo` until the expansion strategy is confirmed as removed

### Requirement: Resource has no Update operation (immutable attachment)
The resource SHALL NOT support update operations. All top-level fields besides `instance_id` are immutable. If the `Update` method is invoked with changes to any immutable field, it SHALL return an error indicating that the resource must be recreated.

#### Scenario: Attempting to update an immutable field
- **WHEN** user changes `type`, `expand_cpu`, `auto_strategy`, `time_interval_strategy`, or `period_strategy`
- **THEN** the Update function SHALL detect changes in `immutableArgs` and return an error requiring resource recreation

### Requirement: Resource supports Import
The resource SHALL support Terraform Import via `schema.ImportStatePassthrough`. The import ID SHALL be the `instance_id` of the CDB instance.

#### Scenario: Importing an existing expansion configuration
- **WHEN** user imports the resource using an `instance_id`
- **THEN** the resource SHALL call `DescribeCPUExpandStrategyInfo` with the imported ID and populate the Terraform state

### Requirement: Resource registered in provider
The resource SHALL be registered in `tencentcloud/provider.go` with the key `tencentcloud_cdb_start_cpu_expand` and corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cdb_start_cpu_expand` SHALL be available as a valid Terraform resource type

### Requirement: Unit tests using gomonkey mock
The resource SHALL have unit tests in `resource_tc_cdb_start_cpu_expand_attachment_test.go` that mock the cloud API calls using gomonkey. Tests SHALL be executable with `go test -gcflags=all=-l` and SHALL NOT use Terraform acceptance test suite.

#### Scenario: Create operation test
- **WHEN** unit test mocks `StartCpuExpand` API to return a valid response
- **THEN** the test SHALL verify that the Create function correctly sets the resource ID and fields

#### Scenario: Read operation test
- **WHEN** unit test mocks `DescribeCPUExpandStrategyInfo` API
- **THEN** the test SHALL verify that the Read function correctly populates the schema fields

#### Scenario: Delete operation test
- **WHEN** unit test mocks `StopCpuExpand` API to return a valid response
- **THEN** the test SHALL verify that the Delete function correctly removes the resource

### Requirement: Documentation file
The resource SHALL have a `.md` documentation file following the gendoc format with:
- One-sentence description mentioning the CDB cloud product
- Example Usage section showing all four expansion types
- Import section (since this is RESOURCE_KIND_ATTACHMENT)

#### Scenario: Documentation completeness
- **WHEN** the documentation file is generated
- **THEN** it SHALL contain description, example usage for auto/manual/timeInterval/period types, and import instructions
