# ssm-service-action Specification

## Purpose
Activates (opens) the Tencent Cloud SSM (Secrets Manager) service for the current account in PostPay (pay-as-you-go) mode via a one-time operation resource `tencentcloud_ssm_service_operation`.
## Requirements
### Requirement: SDKv2 Operation Resource Implementation Pattern

The system SHALL implement `tencentcloud_ssm_service_operation` as an SDKv2 operation resource (`RESOURCE_KIND_OPERATION`) following the reference implementations `resource_tc_dnspod_modify_record_group_operation.go` and `resource_tc_sts_assume_role_operation.go`. The resource SHALL expose only a `Create`/`Read`/`Delete` lifecycle (no `Update`), and the factory `ResourceTencentCloudSsmServiceOperation` SHALL return `*schema.Resource`.

#### Scenario: Operation resource implements the SDKv2 schema.Resource interface

- **WHEN** the resource is defined
- **THEN** `ResourceTencentCloudSsmServiceOperation()` SHALL return a `*schema.Resource`
- **AND** the resource SHALL set `Create`, `Read`, and `Delete` handlers
- **AND** the resource SHALL NOT set an `Update` handler

#### Scenario: Not implemented as framework action

- **WHEN** the change is implemented
- **THEN** there SHALL be NO `action_tc_ssm_service.go` framework action file
- **AND** the action factory `ssm.NewSsmServiceAction` SHALL NOT be registered in `tencentcloud/framework/registry.go`'s `actionFactories`

### Requirement: Resource Schema Definition

The resource schema SHALL define exactly one required input attribute:
- `pay_mode` (`schema.TypeString`, Required, ForceNew): billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.

The `client_token` required by the SSM `CreateService` API SHALL NOT be exposed as a schema field; it SHALL be generated internally via `helper.BuildToken()`. No other attributes SHALL be defined.

#### Scenario: Schema has required pay_mode attribute

- **WHEN** the resource schema is inspected
- **THEN** `Schema["pay_mode"]` SHALL be a `schema.TypeString` that is `Required` and `ForceNew`
- **AND** `Schema` SHALL NOT contain a `client_token` key

### Requirement: Resource Create Operation

The resource `Create` handler SHALL:
1. Read `pay_mode` from the resource data.
2. Generate `clientToken` internally via `helper.BuildToken()` (not exposed as a schema field).
3. Construct an `SsmService` and call `SsmService.CreateService(ctx, payMode, clientToken)` (which wraps `CreateServiceWithContext` with `WriteRetryTimeout` retry).
4. Log the `CreateService` response `ResourceId`/`OrderId` (nil-safe).
5. Because `CreateService` is asynchronous, set `d.SetId(helper.BuildToken())` BEFORE polling so the state retains the id even if polling ultimately fails, allowing `terraform destroy` to clean up.
6. Poll `SsmService.GetServiceStatus(ctx)` (retry-wrapped with `ReadRetryTimeout`) until `ServiceEnabled == true`:
   - On each poll, if `ServiceEnabled` is nil → non-retryable error.
   - If `ServiceEnabled != nil && *ServiceEnabled == true` → activation confirmed, return success.
   - Otherwise → continue polling within the `ReadRetryTimeout`-bounded `resource.Retry` loop (return `RetryableError` to keep polling); if the wait is exhausted, surface a failure error.
7. On any API error, log and return the error. On `Create` error, the id is not set (the operation did not succeed).

#### Scenario: Successful create with polling

- **WHEN** the resource is created with `pay_mode = "PostPay"`, `CreateServiceWithContext` returns a successful response, and a subsequent `GetServiceStatusWithContext` poll returns `ServiceEnabled == true`
- **THEN** the system SHALL log the `ResourceId`/`OrderId` and the enabled status
- **AND** `d.Id()` SHALL be set
- **AND** no error SHALL be returned

#### Scenario: Polling eventually confirms enabled

- **WHEN** the first `GetServiceStatusWithContext` poll returns `ServiceEnabled == false` and a later poll (within `ReadRetryTimeout`) returns `ServiceEnabled == true`
- **THEN** the system SHALL continue polling until enabled and return success

#### Scenario: CreateService API error

- **WHEN** `CreateServiceWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the Create handler SHALL return it

#### Scenario: Polling timeout / never enabled

- **WHEN** `GetServiceStatusWithContext` keeps returning `ServiceEnabled != true` until `ReadRetryTimeout` is exhausted
- **THEN** the Create handler SHALL return the polling error

### Requirement: Read and Delete are no-ops

The resource `Read` handler SHALL be a no-op (return nil). The resource `Delete` handler SHALL be a no-op (return nil). This is a one-time operation resource; no cloud-side state is read or destroyed.

#### Scenario: Read is a no-op

- **WHEN** the resource `Read` is called
- **THEN** it SHALL return nil without calling any cloud API

#### Scenario: Delete is a no-op

- **WHEN** the resource `Delete` is called
- **THEN** it SHALL return nil without calling any cloud API

### Requirement: SSM Service Layer Methods

The system SHALL add two methods to `SsmService` in `tencentcloud/services/ssm/service_tencentcloud_ssm.go`:

1. `CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)`:
   - Construct `ssm.NewCreateServiceRequest()`; set `request.PayMode = helper.String(payMode)`, `request.ClientToken = helper.String(clientToken)`.
   - Wrap the call in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`.
   - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().CreateServiceWithContext(ctx, request)`.
   - On error: return `tccommon.RetryError(e)`.
   - On nil response/response.Response: return `resource.NonRetryableError(...)`.
   - On success: `log.Printf("[DEBUG]%s api[%s] success, ...")`.
   - Defer error logging: `log.Printf("[CRITAL]%s api[%s] fail, ...")`.
   - Return `*ssm.CreateServiceResponse`.

2. `GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)`:
   - Construct `ssm.NewGetServiceStatusRequest()` (empty).
   - Wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`.
   - Inside retry: `ratelimit.Check(request.GetAction())`, call `me.client.UseSsmClient().GetServiceStatusWithContext(ctx, request)`.
   - On error: return `tccommon.RetryError(e)`.
   - On nil response/response.Response: return `resource.NonRetryableError(...)`.
   - On success: `log.Printf("[DEBUG]...")`; set `response = result.Response`.
   - Return `*ssm.GetServiceStatusResponseParams`.

#### Scenario: CreateService wraps API with retry

- **WHEN** `CreateService` is called on the service
- **THEN** it SHALL use `tccommon.WriteRetryTimeout` as the retry timeout
- **AND** on API error it SHALL return `tccommon.RetryError(e)` from within the retry closure
- **AND** it SHALL log failures via `log.Printf("[CRITAL]...")` in the deferred handler

#### Scenario: GetServiceStatus wraps API with retry

- **WHEN** `GetServiceStatus` is called on the service
- **THEN** it SHALL use `tccommon.ReadRetryTimeout` as the retry timeout
- **AND** on a nil response it SHALL return a non-retryable error
- **AND** on API error it SHALL return `tccommon.RetryError(e)` from within the retry closure

### Requirement: Resource Registration

The system SHALL register `tencentcloud_ssm_service_operation` in `tencentcloud/provider.go`'s `ResourcesMap` mapping to `ssm.ResourceTencentCloudSsmServiceOperation()`. It SHALL NOT be registered in the framework `tencentcloud/framework/registry.go`'s `actionFactories`.

#### Scenario: Resource registered in SDKv2 provider

- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL contain an entry keyed `tencentcloud_ssm_service_operation` mapping to `ssm.ResourceTencentCloudSsmServiceOperation()`

#### Scenario: Not registered as framework action

- **WHEN** the framework provider collects action factories
- **THEN** `actionFactories` SHALL NOT include `ssm.NewSsmServiceAction`

### Requirement: Unit Tests

The system SHALL provide unit tests in `tencentcloud/services/ssm/resource_tc_ssm_service_operation_test.go` using gomonkey to mock the SSM cloud API, testing only business logic (no Terraform acceptance test suite).

#### Scenario: Successful create test

- **WHEN** a test creates the resource with `pay_mode = "PostPay"`, mocks `UseSsmClient`, `CreateServiceWithContext` returning a successful response, and `GetServiceStatusWithContext` returning `ServiceEnabled == true`
- **THEN** the test SHALL assert no error and `d.Id()` is set
- **AND** SHALL assert the `CreateService` request `PayMode` equals `PostPay` and `ClientToken` is non-empty

#### Scenario: Polling across a not-yet-enabled poll test

- **WHEN** the first `GetServiceStatusWithContext` poll returns `ServiceEnabled == false` and the second returns `ServiceEnabled == true`
- **THEN** the test SHALL assert no error and at least two polls occurred

#### Scenario: CreateService API error test

- **WHEN** a test mocks `CreateServiceWithContext` returning an error
- **THEN** the test SHALL assert the create returns an error containing the API error message

#### Scenario: Read no-op test

- **WHEN** a test calls `Read`
- **THEN** the test SHALL assert no error

#### Scenario: Delete no-op test

- **WHEN** a test calls `Delete`
- **THEN** the test SHALL assert no error

#### Scenario: Schema validation test

- **WHEN** a test inspects the resource schema
- **THEN** the test SHALL assert `Create`, `Read`, `Delete` are set, `Update` is nil
- **AND** `pay_mode` is Required/ForceNew
- **AND** `client_token` is NOT present in the schema

### Requirement: Resource Documentation

The system SHALL provide a markdown documentation file `resource_tc_ssm_service_operation.md` with a one-line description mentioning SSM (Secrets Manager), an Example Usage section using the SDKv2 resource block syntax (`resource "tencentcloud_ssm_service_operation" "example" { pay_mode = "PostPay" }`). It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections

- **WHEN** the resource is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning SSM
- **AND** an Example Usage block using `resource "tencentcloud_ssm_service_operation" "example" { pay_mode = "PostPay" }` syntax
