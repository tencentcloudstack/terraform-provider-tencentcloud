## ADDED Requirements

### Requirement: Framework Action Implementation Pattern
The system SHALL implement `tencentcloud_dlc_initialize_tc_lake` as a Terraform Plugin Framework `action.Action` (not an SDKv2 operation resource), following the reference implementations `NewTeoConfirmOriginAclUpdate` and `NewBdrcRunCopyPairTasks`. The action struct SHALL embed `fw.ActionWithConfigure` and the factory `NewDlcInitializeTCLake` SHALL return `action.Action`.

#### Scenario: Action implements the framework action interface
- **WHEN** the action is defined
- **THEN** the `DlcInitializeTCLake` struct SHALL embed `fw.ActionWithConfigure`
- **AND** the factory `NewDlcInitializeTCLake()` SHALL return `&DlcInitializeTCLake{}` as `action.Action`
- **AND** the struct SHALL satisfy `action.ActionWithConfigure` at compile time

#### Scenario: Not implemented as SDKv2 operation resource
- **WHEN** the change is implemented
- **THEN** there SHALL be NO `resource_tc_dlc_initialize_tc_lake_operation.go` SDKv2 resource file
- **AND** the action SHALL NOT be registered in `tencentcloud/provider.go`'s `ResourcesMap`

### Requirement: Action Metadata
The action `Metadata` handler SHALL set `resp.TypeName` to `tencentcloud_dlc_initialize_tc_lake`.

#### Scenario: Metadata sets the action type name
- **WHEN** the framework calls `Metadata`
- **THEN** `resp.TypeName` SHALL equal `tencentcloud_dlc_initialize_tc_lake`

### Requirement: Action Schema Definition
The action `Schema` handler SHALL define a schema with a `Description` and an empty `Attributes` map (`map[string]schema.Attribute{}`), because the `InitializeTCLake` API takes no request parameters.

#### Scenario: Schema has no input attributes
- **WHEN** the framework calls `Schema`
- **THEN** `resp.Schema.Attributes` SHALL be an empty map
- **AND** `resp.Schema.Description` SHALL mention DLC and the TCLake initialization operation

### Requirement: No Output Attributes
The action schema SHALL NOT define `instance_id` or `is_success` as attributes, because the Terraform Plugin Framework action schema does not support Computed/output attributes and `InvokeResponse` has no output/state-writing mechanism. The API response fields `InstanceId` and `IsSuccess` SHALL be logged via `log.Printf` for observability only.

#### Scenario: Response fields logged not exposed
- **WHEN** the `InitializeTCLake` API returns successfully
- **THEN** the service layer SHALL log the response body via `log.Printf("[DEBUG]...")`
- **AND** the action schema SHALL NOT contain `instance_id` or `is_success` attributes
- **AND** no output/state SHALL be written to `InvokeResponse` beyond diagnostics

### Requirement: Action Invoke Operation
The action `Invoke` handler SHALL call the DLC `InitializeTCLake` API via the `DlcService.InitializeTCLake` service method (which calls `InitializeTCLakeWithContext`). The operation is synchronous and SHALL NOT poll any Read interface afterwards. No cloud-side state SHALL be persisted (no id, no output attributes written).

#### Scenario: Provider client not configured
- **WHEN** the provider client (`a.Client()`) is nil at invoke time
- **THEN** the system SHALL add a diagnostic error with summary "Provider not configured" and return without calling the API

#### Scenario: Successful TCLake initialization
- **WHEN** the action is invoked and the provider client is configured
- **THEN** the system SHALL call `DlcService.InitializeTCLake(ctx)` which constructs an empty `InitializeTCLakeRequest`, calls `InitializeTCLakeWithContext` with retry, logs the response, and returns no diagnostics on success

#### Scenario: Retryable API error
- **WHEN** `InitializeTCLakeWithContext` returns a retryable error
- **THEN** the service layer SHALL retry within `tccommon.WriteRetryTimeout` via `resource.Retry`, wrapping the error with `tccommon.RetryError`

#### Scenario: Non-retryable API error surfaced
- **WHEN** `InitializeTCLakeWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the `Invoke` handler SHALL add it to `resp.Diagnostics` with a summary referencing initializing TCLake

#### Scenario: No state persisted after successful invoke
- **WHEN** the invoke completes successfully
- **THEN** the system SHALL NOT set any id, SHALL NOT write any output/computed attribute, and SHALL NOT implement Read/Update/Delete lifecycle methods

### Requirement: DLC Service Layer Method
The system SHALL add a method `InitializeTCLake(ctx context.Context) (*dlc.InitializeTCLakeResponse, error)` to `DlcService` in `tencentcloud/services/dlc/service_tencentcloud_dlc.go`. The method SHALL construct an empty `dlc.NewInitializeTCLakeRequest()`, call `me.client.UseDlcClient().InitializeTCLakeWithContext(ctx, request)` inside a `resource.Retry(tccommon.WriteRetryTimeout, ...)` block with `ratelimit.Check(request.GetAction())` and `tccommon.RetryError(e)` error wrapping, and log the request/response via `log.Printf`.

#### Scenario: Service method wraps API with retry
- **WHEN** `InitializeTCLake` is called on the service
- **THEN** it SHALL use `tccommon.WriteRetryTimeout` as the retry timeout
- **AND** on API error it SHALL return `tccommon.RetryError(e)` from within the retry closure
- **AND** it SHALL log failures via `log.Printf("[CRITAL]...")` in the deferred handler

#### Scenario: Service method logs success
- **WHEN** the API call succeeds
- **THEN** the service method SHALL log the request and response via `log.Printf("[DEBUG]...")`

### Requirement: Action Registration
The system SHALL register the `tencentcloud_dlc_initialize_tc_lake` action factory `dlc.NewDlcInitializeTCLake` in `tencentcloud/framework/registry.go`'s `actionFactories` slice, and SHALL add the `dlc` services package import to `registry.go`. It SHALL NOT be registered in the SDKv2 `provider.go` `ResourcesMap`.

#### Scenario: Action factory registered
- **WHEN** the framework provider collects action factories
- **THEN** `actionFactories` SHALL include `dlc.NewDlcInitializeTCLake`
- **AND** `registry.go` SHALL import `tencentcloud/services/dlc`

#### Scenario: Not registered in SDKv2 provider
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain an entry keyed `tencentcloud_dlc_initialize_tc_lake_operation`

### Requirement: Unit Tests
The system SHALL provide unit tests in `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake_test.go` using gomonkey to mock the DLC cloud API, testing only business logic (no Terraform acceptance test suite).

#### Scenario: Successful invoke test
- **WHEN** a test invokes the action with the provider client configured and mocks `InitializeTCLakeWithContext` to return a successful response
- **THEN** the test SHALL assert no error diagnostic is produced

#### Scenario: API error test
- **WHEN** a test invokes the action and mocks `InitializeTCLakeWithContext` to return an error
- **THEN** the test SHALL assert an error diagnostic is produced containing the API error message

#### Scenario: Provider client not configured test
- **WHEN** a test invokes the action without configuring the provider client
- **THEN** the test SHALL assert an error diagnostic is produced with summary "Provider not configured"

#### Scenario: Metadata and schema validation test
- **WHEN** a test calls `Metadata` and `Schema` on the action
- **THEN** the test SHALL assert `TypeName` equals `tencentcloud_dlc_initialize_tc_lake`
- **AND** SHALL assert the attributes map is empty
- **AND** SHALL assert no diagnostics are produced

### Requirement: Action Documentation
The system SHALL provide a markdown documentation file `action_tc_dlc_initialize_tc_lake.md` with a one-line description mentioning DLC, an Example Usage section using the framework action `action` block syntax with an empty `config {}` block, and a NOTE about Terraform version support. It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections
- **WHEN** the action is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning DLC
- **AND** an Example Usage block using `action "tencentcloud_dlc_initialize_tc_lake" "example" { config {} }` syntax
- **AND** a NOTE about Terraform 1.14+ support