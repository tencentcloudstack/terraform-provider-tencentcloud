# bdrc-copy-pair-tasks-action Specification

## Purpose
TBD - created by archiving change add-bdrc-copy-pair-tasks-operation. Update Purpose after archive.
## Requirements
### Requirement: BDRC SDK Client Connectivity
The system SHALL register a BDRC SDK client accessor `UseBdrcV20260330Client()` in `tencentcloud/connectivity/client.go`, returning a `*bdrcv20260330.Client`, with the bdrc/v20260330 package import and a lazily-initialized client field on `TencentCloudClient`.

#### Scenario: BDRC client is lazily initialized
- **WHEN** `UseBdrcV20260330Client()` is called for the first time
- **THEN** the system SHALL create the bdrc client with the provider credential, region, and a 300-second timeout profile, attach a `LogRoundTripper`, cache it on the `TencentCloudClient`, and return it

#### Scenario: BDRC client is reused on subsequent calls
- **WHEN** `UseBdrcV20260330Client()` is called again on the same client
- **THEN** the system SHALL return the cached bdrc client without creating a new one

### Requirement: Action Resource Schema Definition
The system SHALL define a Terraform Plugin Framework action `tencentcloud_bdrc_run_copy_pair_tasks` with the following schema attributes:
- `copy_pair_ids` (Required, List of String): 复制对 ID 列表
- `copy_pair_type` (Required, String): 要启动复制对的类型，取值 DISK / INSTANCE / CFS

#### Scenario: Schema defines both required input attributes
- **WHEN** the action schema is defined
- **THEN** it SHALL include `copy_pair_ids` as a required list-of-string attribute and `copy_pair_type` as a required string attribute

#### Scenario: Required attributes enforced
- **WHEN** a user omits `copy_pair_ids` or `copy_pair_type` in the action config
- **THEN** the framework SHALL report a validation error before the action is invoked

### Requirement: Action Invoke Operation
The system SHALL implement the action `Invoke` handler to call the BDRC `RunCopyPairTasks` API via `RunCopyPairTasksWithContext`, passing `CopyPairIds` and `CopyPairType` from the action config. The operation is synchronous and SHALL NOT poll any Read interface afterwards. No cloud-side state SHALL be persisted (no id, no output attributes written).

#### Scenario: Successful copy pair tasks launch
- **WHEN** the action is invoked with valid `copy_pair_ids` and `copy_pair_type`
- **THEN** the system SHALL construct a `RunCopyPairTasksRequest` with `CopyPairIds` and `CopyPairType`, call `RunCopyPairTasksWithContext`, log the request/response, and return no diagnostics on success

#### Scenario: Missing required input rejected before API call
- **WHEN** `copy_pair_ids` is null/unknown/empty or `copy_pair_type` is null/unknown/empty at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Provider client not configured
- **WHEN** the provider client (`a.Client()`) is nil at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Retryable API error
- **WHEN** `RunCopyPairTasksWithContext` returns a retryable error
- **THEN** the service layer SHALL retry within `tccommon.WriteRetryTimeout` via `resource.Retry`, wrapping the error with `tccommon.RetryError`

#### Scenario: Non-retryable API error surfaced
- **WHEN** `RunCopyPairTasksWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the `Invoke` handler SHALL add it to `resp.Diagnostics`

#### Scenario: No state persisted after successful invoke
- **WHEN** the invoke completes successfully
- **THEN** the system SHALL NOT set any id, SHALL NOT write any output/computed attribute, and SHALL NOT implement Read/Update/Delete lifecycle methods

### Requirement: Action Registration
The system SHALL register the `tencentcloud_bdrc_run_copy_pair_tasks` action factory `bdrc.NewBdrcRunCopyPairTasks` in `tencentcloud/framework/registry.go`'s `actionFactories` slice, and SHALL NOT register it in the SDKv2 `provider.go` `ResourcesMap`.

#### Scenario: Action factory registered
- **WHEN** the framework provider collects action factories
- **THEN** `frameworkActions()` SHALL include `bdrc.NewBdrcRunCopyPairTasks`

#### Scenario: Not registered in SDKv2 provider
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain an entry keyed `tencentcloud_bdrc_run_copy_pair_tasks_operation`

### Requirement: Unit Tests
The system SHALL provide unit tests in `tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks_test.go` using gomonkey to mock the BDRC cloud API, testing only business logic (no Terraform acceptance test suite).

#### Scenario: Successful invoke test
- **WHEN** a test invokes the action with valid inputs and mocks `RunCopyPairTasksWithContext` to return a successful response
- **THEN** the test SHALL assert the mocked API is called with the expected `CopyPairIds` and `CopyPairType`, and that no error diagnostic is produced

#### Scenario: API error test
- **WHEN** a test invokes the action and mocks `RunCopyPairTasksWithContext` to return an error
- **THEN** the test SHALL assert an error diagnostic is produced containing the API error message

### Requirement: Action Documentation
The system SHALL provide a markdown documentation file `action_tc_bdrc_run_copy_pair_tasks.md` with a one-line description mentioning BDRC, an Example Usage section using the framework action `action` block syntax, and a NOTE about Terraform version support. It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections
- **WHEN** the action is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning BDRC, an Example Usage block using `action "tencentcloud_bdrc_run_copy_pair_tasks" "example" { config { ... } }` syntax, and a NOTE about Terraform 1.14+ support

