# mongodb-restore-db-instance-action Specification

## Purpose
TBD - created by syncing change add-mongodb-restore-db-instance. Update Purpose after archive.
## Requirements
### Requirement: Action Resource Schema Definition
The system SHALL define a Terraform Plugin Framework action `tencentcloud_mongodb_restore_db_instance` with the following schema attributes:
- `instance_id` (Required, String): 待回档的 MongoDB 实例 ID，格式 `cmgo-xxxxxxxx`。
- `restore_time` (Required, String): 回档目标时间点，格式 `YYYY-MM-DD hh:mm:ss`，须处于实例的备份保留期内。
- `databases` (Required, List Nested Block): 回档的库表信息列表，每项含 `db`（Required, String）与 `collections`（Required, List Nested Block）。
- `collections` (Required, List Nested Block, located inside `databases`): 集合回档信息，每项含 `old_collection`（Required, String）与 `new_collection`（Required, String）。

#### Scenario: Schema defines required input attributes
- **WHEN** the action schema is defined
- **THEN** it SHALL include `instance_id` as a required string attribute, `restore_time` as a required string attribute, and `databases` as a required list-nested-block
- **AND** each `databases` block SHALL include `db` as a required string and `collections` as a required list-nested-block
- **AND** each `collections` block SHALL include `old_collection` and `new_collection` as required strings

#### Scenario: Required attributes enforced
- **WHEN** a user omits `instance_id`, `restore_time`, or `databases` in the action config
- **THEN** the framework SHALL report a validation error before the action is invoked

### Requirement: Action Invoke Operation
The system SHALL implement the action `Invoke` handler to call the MongoDB `RestoreDBInstance` API via `RestoreDBInstanceWithContext`, passing `InstanceId`, `RestoreTime`, and `Databases` (constructed from the nested blocks) from the action config. Because `RestoreDBInstance` is asynchronous and returns a `FlowId`, the system SHALL poll the task status via `MongodbService.DescribeAsyncRequestInfo` using the returned `FlowId` until the task succeeds (`success`), fails (`failed`), or times out. No cloud-side state SHALL be persisted (no id, no output attributes written).

#### Scenario: Successful restore invocation
- **WHEN** the action is invoked with valid `instance_id`, `restore_time`, and `databases` (containing `db` and `collections` with `old_collection`/`new_collection`)
- **THEN** the system SHALL construct a `RestoreDBInstanceRequest` with `InstanceId`, `RestoreTime`, and `Databases` (each `RestoreDatabases` having `Db` and `Collections` of `RestoreCollection`)
- **AND** call `RestoreDBInstanceWithContext` via the service layer with `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- **AND** extract the `FlowId` from the response
- **AND** poll `DescribeAsyncRequestInfo` with the `FlowId` until the task status is `success`
- **AND** return no diagnostics on success

#### Scenario: Async task failure surfaced
- **WHEN** `DescribeAsyncRequestInfo` reports task status `failed`
- **THEN** the system SHALL return an error diagnostic describing the restore task failure

#### Scenario: Async task timeout
- **WHEN** the async task does not complete within the poll timeout (`3 * tccommon.ReadRetryTimeout`)
- **THEN** the system SHALL return a timeout error diagnostic

#### Scenario: Missing required input rejected before API call
- **WHEN** `instance_id` is null/unknown/empty, or `restore_time` is null/unknown/empty, or `databases` is null/unknown/empty at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Provider client not configured
- **WHEN** the provider client (`a.Client()`) is nil at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Retryable API error on RestoreDBInstance
- **WHEN** `RestoreDBInstanceWithContext` returns a retryable error
- **THEN** the service layer SHALL retry within `tccommon.WriteRetryTimeout` via `resource.Retry`, wrapping the error with `tccommon.RetryError`

#### Scenario: Non-retryable API error surfaced
- **WHEN** `RestoreDBInstanceWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the `Invoke` handler SHALL add it to `resp.Diagnostics`

#### Scenario: No state persisted after successful invoke
- **WHEN** the invoke completes successfully
- **THEN** the system SHALL NOT set any id, SHALL NOT write any output/computed attribute (including `flow_id`), and SHALL NOT implement Read/Update/Delete lifecycle methods

### Requirement: Service Layer RestoreDBInstance Method
The MongoDB service layer SHALL provide a `RestoreDBInstance(ctx, instanceId, restoreTime string, databases []*mongodb.RestoreDatabases) (flowId int64, errRet error)` method that constructs a `RestoreDBInstanceRequest`, fills `InstanceId`, `RestoreTime`, and `Databases`, wraps the `RestoreDBInstanceWithContext` call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `ratelimit.Check` and `tccommon.RetryError`, logs success/failure, and returns the `FlowId` from the response on success.

#### Scenario: Service method constructs and sends request
- **WHEN** `RestoreDBInstance` service method is called with instance ID, restore time, and databases
- **THEN** the system SHALL create `mongodb.NewRestoreDBInstanceRequest()`, set `InstanceId`, `RestoreTime`, and `Databases`
- **AND** call `RestoreDBInstanceWithContext` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, invoking `ratelimit.Check(request.GetAction())` before the call
- **AND** on a retryable error return `tccommon.RetryError(e)`

#### Scenario: Service method returns FlowId on success
- **WHEN** `RestoreDBInstanceWithContext` succeeds
- **THEN** the service method SHALL verify `result.Response != nil` and `result.Response.FlowId != nil`
- **AND** return the dereferenced `FlowId` (int64)

#### Scenario: Service method records failure log
- **WHEN** the service method returns an error
- **THEN** a deferred function SHALL log `[CRITAL]` with the log ID, action name, request body, and error reason

### Requirement: Action Registration
The system SHALL register the `tencentcloud_mongodb_restore_db_instance` action factory `mongodb.NewMongodbRestoreDbInstance` in `tencentcloud/framework/registry.go`'s `actionFactories` slice, and SHALL NOT register it in the SDKv2 `provider.go` `ResourcesMap`.

#### Scenario: Action factory registered
- **WHEN** the framework provider collects action factories
- **THEN** `frameworkActions()` SHALL include `mongodb.NewMongodbRestoreDbInstance`

#### Scenario: Not registered in SDKv2 provider
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain an entry keyed `tencentcloud_mongodb_restore_db_instance_operation`

### Requirement: Unit Tests
The system SHALL provide unit tests in `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance_test.go` using gomonkey to mock the MongoDB cloud API, testing only business logic (no Terraform acceptance test suite).

#### Scenario: Successful invoke test
- **WHEN** a test invokes the action with valid inputs and mocks `UseMongodbClient`, `RestoreDBInstanceWithContext` to return a `FlowId`, and `DescribeAsyncRequestInfo` to report `success`
- **THEN** the test SHALL assert the mocked API is called with the expected `InstanceId`/`RestoreTime`/`Databases`, and that no error diagnostic is produced

#### Scenario: API error test
- **WHEN** a test invokes the action and mocks `RestoreDBInstanceWithContext` to return an error
- **THEN** the test SHALL assert an error diagnostic is produced containing the API error message

#### Scenario: Async task failure test
- **WHEN** a test mocks `RestoreDBInstanceWithContext` to return a `FlowId` but `DescribeAsyncRequestInfo` to report `failed`
- **THEN** the test SHALL assert an error diagnostic is produced

#### Scenario: Missing input test
- **WHEN** a test invokes the action with `instance_id`, `restore_time`, or `databases` null/empty
- **THEN** the test SHALL assert an error diagnostic is produced and the API is not called

#### Scenario: Client not configured test
- **WHEN** the provider client is not set on the action
- **THEN** the test SHALL assert an error diagnostic mentioning "Provider not configured"

#### Scenario: Metadata and schema validation test
- **WHEN** a test calls `Metadata` and `Schema`
- **THEN** the test SHALL assert `TypeName == "tencentcloud_mongodb_restore_db_instance"` and that all required attributes/blocks exist with correct types

### Requirement: Action Documentation
The system SHALL provide a markdown documentation file `action_tc_mongodb_restore_db_instance.md` with a one-line description mentioning MongoDB, an Example Usage section using the framework action `action` block syntax (including nested `databases`/`collections` blocks), and a NOTE about Terraform version support. It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections
- **WHEN** the action is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning MongoDB, an Example Usage block using `action "tencentcloud_mongodb_restore_db_instance" "example" { config { ... } }` syntax with nested `databases`/`collections` blocks, and a NOTE about Terraform 1.14+ support
