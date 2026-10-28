## ADDED Requirements

### Requirement: Framework Action Implementation Pattern

The system SHALL implement `tencentcloud_ssm_service` as a Terraform Plugin Framework `action.Action` (not an SDKv2 operation resource), following the reference implementations `NewDlcInitializeTCLake` and `NewAntiddosUnblockResources`. The action struct SHALL embed `fw.ActionWithConfigure` and the factory `NewSsmServiceAction` SHALL return `action.Action`.

#### Scenario: Action implements the framework action interface
- **WHEN** the action is defined
- **THEN** the `SsmServiceAction` struct SHALL embed `fw.ActionWithConfigure`
- **AND** the factory `NewSsmServiceAction()` SHALL return `&SsmServiceAction{}` as `action.Action`
- **AND** the struct SHALL satisfy `action.ActionWithConfigure` at compile time

#### Scenario: Not implemented as SDKv2 operation resource
- **WHEN** the change is implemented
- **THEN** there SHALL be NO `resource_tc_ssm_service_operation.go` SDKv2 resource file
- **AND** the action SHALL NOT be registered in `tencentcloud/provider.go`'s `ResourcesMap`

### Requirement: Action Metadata

The action `Metadata` handler SHALL set `resp.TypeName` to `tencentcloud_ssm_service`.

#### Scenario: Metadata sets the action type name
- **WHEN** the framework calls `Metadata`
- **THEN** `resp.TypeName` SHALL equal `tencentcloud_ssm_service`

### Requirement: Action Schema Definition

The action `Schema` handler SHALL define exactly one required input attribute:
- `pay_mode` (`schema.StringAttribute`, Required): billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.

The `client_token` required by the SSM `CreateService` API SHALL NOT be exposed as a schema attribute; it SHALL be generated internally in the `Invoke` handler via `helper.BuildToken()`. No other attributes SHALL be defined.

#### Scenario: Schema has exactly one required input attribute
- **WHEN** the framework calls `Schema`
- **THEN** `resp.Schema.Attributes` SHALL contain a Required `pay_mode` StringAttribute
- **AND** `resp.Schema.Attributes` SHALL NOT contain a `client_token` key
- **AND** the attributes map SHALL contain exactly one entry
- **AND** `resp.Schema.Description` SHALL mention SSM and the service-opening operation

### Requirement: No Output Attributes

The action schema SHALL NOT define `resource_id`, `order_id`, `service_enabled`, `invalid_type`, `access_key_escrow_enabled`, `expire_time`, `qps_limit`, `secret_limit`, `pay_model`, `renew_flag`, or `resource_region` as attributes, because the Terraform Plugin Framework action schema does not support Computed/output attributes and `InvokeResponse` has no output/state-writing mechanism. The API response fields SHALL be logged via `log.Printf` for observability only.

#### Scenario: Response fields logged not exposed
- **WHEN** the `CreateService` and `GetServiceStatus` APIs return successfully
- **THEN** the service layer SHALL log the response body via `log.Printf("[DEBUG]...")`
- **AND** the action schema SHALL NOT contain any output attribute
- **AND** no output/state SHALL be written to `InvokeResponse` beyond diagnostics

### Requirement: Action Invoke Operation

The action `Invoke` handler SHALL collect `pay_mode` from the action config, then:
1. Validate `pay_mode` is present and non-empty (null/unknown/empty → add a diagnostic error and return without calling the API).
2. When `a.Client()` is nil, add a "Provider not configured" diagnostic and return.
3. Generate `clientToken` internally via `helper.BuildToken()`; it SHALL NOT be read from the action config.
4. Call `SsmService.CreateService(ctx, payMode, clientToken)` (which wraps `CreateServiceWithContext` with `WriteRetryTimeout` retry).
5. Log the `CreateService` response `ResourceId`/`OrderId` (nil-safe).
6. Because `CreateService` is asynchronous, poll `SsmService.GetServiceStatus(ctx)` (retry-wrapped with `ReadRetryTimeout`) until `ServiceEnabled == true`:
   - On each poll, if `ServiceEnabled` is nil → non-retryable error (surfaced as a diagnostic).
   - If `ServiceEnabled != nil && *ServiceEnabled == true` → activation confirmed, return success.
   - Otherwise → continue polling within the `ReadRetryTimeout`-bounded `resource.Retry` loop (return `RetryableError` to keep polling); if the wait is exhausted, surface a failure diagnostic.
7. On any API error, add a diagnostic error whose summary references opening SSM service (for CreateService) or polling SSM service status (for GetServiceStatus), and return. No cloud-side state SHALL be persisted (no id, no output attributes written).

#### Scenario: Missing pay_mode
- **WHEN** the action is invoked with `pay_mode` null/unknown/empty
- **THEN** the system SHALL add an error diagnostic with summary "Missing pay_mode"
- **AND** SHALL NOT call the CreateService API

#### Scenario: client_token generated internally
- **WHEN** the action is invoked with a valid `pay_mode`
- **THEN** the `CreateService` request `ClientToken` SHALL be a non-empty token generated by `helper.BuildToken()`
- **AND** the action config SHALL NOT provide any `client_token` value

#### Scenario: Provider client not configured
- **WHEN** the provider client (`a.Client()`) is nil at invoke time
- **THEN** the system SHALL add a diagnostic error with summary "Provider not configured" and return without calling the API

#### Scenario: Successful activation with polling
- **WHEN** the action is invoked with a valid `pay_mode`, `CreateServiceWithContext` returns a successful response, and a subsequent `GetServiceStatusWithContext` poll returns `ServiceEnabled == true`
- **THEN** the system SHALL log the `ResourceId`/`OrderId` and the enabled status
- **AND** no error diagnostic SHALL be produced

#### Scenario: Polling eventually confirms enabled
- **WHEN** the first `GetServiceStatusWithContext` poll returns `ServiceEnabled == false` and a later poll (within `ReadRetryTimeout`) returns `ServiceEnabled == true`
- **THEN** the system SHALL continue polling until enabled and return success with no error diagnostic

#### Scenario: CreateService API error
- **WHEN** `CreateServiceWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the Invoke handler SHALL add it to `resp.Diagnostics` with a summary referencing opening SSM service

#### Scenario: GetServiceStatus API error during polling
- **WHEN** `GetServiceStatusWithContext` returns a non-retryable error
- **THEN** the Invoke handler SHALL add an error diagnostic with a summary referencing polling SSM service status

#### Scenario: Polling timeout / never enabled
- **WHEN** `GetServiceStatusWithContext` keeps returning `ServiceEnabled != true` until `ReadRetryTimeout` is exhausted
- **THEN** the Invoke handler SHALL add an error diagnostic with a summary referencing polling SSM service status

#### Scenario: No state persisted after successful invoke
- **WHEN** the invoke completes successfully
- **THEN** the system SHALL NOT set any id, SHALL NOT write any output/computed attribute, and SHALL NOT implement Read/Update/Delete lifecycle methods

### Requirement: SSM Service Layer Methods

The system SHALL add two methods to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`:

1. `CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)`:
   - Construct `ssm.NewCreateServiceRequest()`; set `request.PayMode = helper.String(payMode)`, `request.ClientToken = helper.String(clientToken)`.
   - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`.
   - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().CreateServiceWithContext(ctx, request)`.
   - On error: return `tccommon.RetryError(e)`.
   - On nil response/response.Response: return `resource.NonRetryableError(...)`.
   - On success: `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)`.
   - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", ...)`.
   - Return `*ssm.CreateServiceResponse`.

2. `GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)`:
   - Construct `ssm.NewGetServiceStatusRequest()` (empty).
   - Wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`.
   - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().GetServiceStatusWithContext(ctx, request)`.
   - On error: return `tccommon.RetryError(e)`.
   - On response==nil or response.Response==nil: return `resource.NonRetryableError(fmt.Errorf("GetServiceStatus failed, Response is nil"))`.
   - On success: `log.Printf("[DEBUG]...")`; set `response = result.Response`.
   - Return `*ssm.GetServiceStatusResponseParams`.

The polling loop itself lives in the Invoke handler, reusing the retry-wrapped `GetServiceStatus`: while `ServiceEnabled != true` return `RetryableError` so the Read-retry continues; nil-safe access is mandatory.

#### Scenario: CreateService wraps API with retry
- **WHEN** `CreateService` is called on the service
- **THEN** it SHALL use `tccommon.WriteRetryTimeout` as the retry timeout
- **AND** on API error it SHALL return `tccommon.RetryError(e)` from within the retry closure
- **AND** it SHALL log failures via `log.Printf("[CRITAL]...")` in the deferred handler

#### Scenario: CreateService logs success
- **WHEN** the CreateService call succeeds
- **THEN** the service method SHALL log the request and response via `log.Printf("[DEBUG]...")`

#### Scenario: GetServiceStatus wraps API with retry
- **WHEN** `GetServiceStatus` is called on the service
- **THEN** it SHALL use `tccommon.ReadRetryTimeout` as the retry timeout
- **AND** on a nil response it SHALL return a non-retryable error
- **AND** on API error it SHALL return `tccommon.RetryError(e)` from within the retry closure

### Requirement: Vendored SDK Support

The vendored SSM SDK SHALL expose `CreateService` (`CreateServiceWithContext`) and `GetServiceStatus` (`GetServiceStatusWithContext`). The repository SHALL vendor `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm v1.3.199` (with the matching `common v1.3.199`) because the previously vendored `v1.3.164` does not contain `CreateService`.

#### Scenario: Vendored SDK exposes CreateService
- **WHEN** the code is compiled
- **THEN** `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923/client.go` SHALL contain `CreateServiceWithContext`
- **AND** `models.go` SHALL contain `CreateServiceRequest` / `CreateServiceResponse`
- **AND** `go.mod` and `vendor/modules.txt` SHALL reference `tencentcloud-sdk-go/tencentcloud/ssm v1.3.199`

### Requirement: Action Registration

The system SHALL register the `tencentcloud_ssm_service` action factory `ssm.NewSsmServiceAction` in `tencentcloud/framework/registry.go`'s `actionFactories` slice. The `ssm` services package import is ALREADY present in `registry.go` (used by the framework resource + ephemeral), so no new import line is needed. It SHALL NOT be registered in the SDKv2 `provider.go` `ResourcesMap`, and `tencentcloud/provider.md` SHALL list `tencentcloud_ssm_service` under the `Action` sub-section of the SSM product node only.

#### Scenario: Action factory registered
- **WHEN** the framework provider collects action factories
- **THEN** `actionFactories` SHALL include `ssm.NewSsmServiceAction`
- **AND** `registry.go` SHALL already import `tencentcloud/services/ssm` (no duplicate import added)

#### Scenario: Not registered in SDKv2 provider
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain an entry keyed `tencentcloud_ssm_service`
- **AND** it SHALL NOT contain an entry keyed `tencentcloud_ssm_service_operation`

### Requirement: Unit Tests

The system SHALL provide unit tests in `tencentcloud/services/ssm/action_tc_ssm_service_test.go` using gomonkey to mock the SSM cloud API, testing only business logic (no Terraform acceptance test suite), mirroring `action_tc_antiddos_unblock_resources_test.go`.

#### Scenario: Successful invoke test
- **WHEN** a test invokes the action with a valid `pay_mode`, mocks `UseSsmClient`, `CreateServiceWithContext` returning a successful response, and `GetServiceStatusWithContext` returning `ServiceEnabled == true`
- **THEN** the test SHALL assert no error diagnostic is produced
- **AND** SHALL assert the `CreateService` request `PayMode` equals `PostPay` and `ClientToken` is non-empty (generated internally)

#### Scenario: Polling across a not-yet-enabled poll test
- **WHEN** the first `GetServiceStatusWithContext` poll returns `ServiceEnabled == false` and the second returns `ServiceEnabled == true`
- **THEN** the test SHALL assert no error diagnostic is produced (polling succeeded)
- **AND** SHALL assert at least two polls occurred

#### Scenario: Missing pay_mode test
- **WHEN** a test invokes the action with `pay_mode` null
- **THEN** the test SHALL assert an error diagnostic with summary "Missing pay_mode" is produced
- **AND** SHALL assert the CreateService API was not called

#### Scenario: CreateService API error test
- **WHEN** a test invokes the action and mocks `CreateServiceWithContext` returning an error
- **THEN** the test SHALL assert an error diagnostic is produced containing the API error message

#### Scenario: Polling error test
- **WHEN** a test invokes the action and mocks `GetServiceStatusWithContext` returning an error
- **THEN** the test SHALL assert an error diagnostic whose summary references polling SSM service status is produced

#### Scenario: Provider client not configured test
- **WHEN** a test invokes the action without configuring the provider client
- **THEN** the test SHALL assert an error diagnostic is produced with summary "Provider not configured"

#### Scenario: Metadata and schema validation test
- **WHEN** a test calls `Metadata` and `Schema` on the action
- **THEN** the test SHALL assert `TypeName` equals `tencentcloud_ssm_service`
- **AND** SHALL assert the attributes map contains exactly `pay_mode` (Required) and does NOT contain `client_token`
- **AND** SHALL assert no diagnostics are produced

### Requirement: Action Documentation

The system SHALL provide a markdown documentation file `action_tc_ssm_service.md` with a one-line description mentioning SSM, an Example Usage section using the framework action block syntax (`action "tencentcloud_ssm_service" "example" { config { pay_mode = "PostPay" } }`), and a NOTE about Terraform version support. The example SHALL NOT include `client_token`. It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections
- **WHEN** the action is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning SSM
- **AND** an Example Usage block using `action "tencentcloud_ssm_service" "example" { config { pay_mode = "PostPay" } }` syntax without `client_token`
- **AND** a NOTE about Terraform 1.14+ support
