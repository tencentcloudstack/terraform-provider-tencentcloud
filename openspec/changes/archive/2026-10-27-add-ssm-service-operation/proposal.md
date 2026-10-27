## Why

The SSM (Secrets Manager) cloud API exposes a `CreateService` operation that activates (opens) the SSM service (PostPay / pay-as-you-go) for the current account. This is a one-time activation step that practitioners currently cannot trigger through Terraform — they must open the service manually in the console before managing any SSM secret resources. Adding a framework-based **action** lets users open SSM as part of their IaC workflow alongside other SSM resources.

## What Changes

- Add a new Terraform Plugin Framework **action** resource `tencentcloud_ssm_service` implemented via the `action.Action` interface (not the SDKv2 operation resource pattern), following the reference implementations `NewDlcInitializeTCLake` (`action_tc_dlc_initialize_tc_lake.go`) and `NewAntiddosUnblockResources` (`action_tc_antiddos_unblock_resources.go`, which takes input arguments).
- Implement the `Invoke` handler to:
  - Collect required inputs `pay_mode` (`request.PayMode`) and `client_token` (`request.ClientToken`) from the action config.
  - Call the SSM `CreateService` API (`CreateServiceWithContext`) to open the service.
  - Because `CreateService` is **asynchronous** (the API doc states "返回时资源可能仍在发货，不代表服务已可用" — the resource may still be provisioning when the call returns), poll the SSM `GetServiceStatus` API until `ServiceEnabled == true` (with `tccommon.ReadRetryTimeout`-bounded retry loop) to confirm activation has taken effect.
- Since the Terraform Plugin Framework `action` schema does **not** support Computed/output attributes and `InvokeResponse` only carries `Diagnostics` + `SendProgress`, none of the cloud API response fields (`resource_id`, `order_id`, `service_enabled`, etc.) are exposed as Terraform attributes. The `CreateService` response fields (`ResourceId`, `OrderId`) and the polled `GetServiceStatus` status are logged via `log.Printf` for observability only. This matches the reference implementations, which invoke an API and surface only success/error diagnostics.
- Add a `CreateService` method and a `GetServiceStatus` polling helper to the SSM service layer (`service_tencentcloud_ssm.go`); the existing `DescribeSsmServiceStatusByFilter` already wraps `GetServiceStatus` but without retry, so a dedicated retry-wrapped method is added for the action. Wrap both API calls with `tccommon.WriteRetryTimeout`/`tccommon.ReadRetryTimeout` retry and `tccommon.RetryError` error wrapping.
- Register the action factory `ssm.NewSsmServiceAction` in `tencentcloud/framework/registry.go`'s `actionFactories` slice. Do NOT register it in the SDKv2 `tencentcloud/provider.go` `ResourcesMap` (the `ssm` services package is already imported by `registry.go`).
- Provide a markdown documentation file `action_tc_ssm_service.md` with a one-line description mentioning SSM, an Example Usage section using the framework action block syntax (`config {}` with `pay_mode`/`client_token`), and a NOTE about Terraform 1.14+ support.
- Add unit tests in `action_tc_ssm_service_test.go` using gomonkey mocks (no Terraform acceptance test suite), covering successful invoke (with polling), missing-input, API error on CreateService, API error/during polling, client-not-configured, and schema/metadata validation.

## Capabilities

### New Capabilities
- `ssm-service-action`: Adds a framework action `tencentcloud_ssm_service` that activates the SSM service (PostPay) via the SSM `CreateService` API and polls `GetServiceStatus` until the service is enabled. This is a one-time operation; no cloud-side state is persisted and no output attributes are exposed (the API response `ResourceId`/`OrderId` and the polled `ServiceEnabled`/`InvalidType`/`AccessKeyEscrowEnabled`/`ExpireTime`/`QPSLimit`/`SecretLimit`/`PayModel`/`RenewFlag`/`resource_id`/`ResourceRegion` fields are logged for observability only).

### Modified Capabilities
<!-- None: brand-new action, no existing requirements change. -->

## Impact

### Affected Code
- `tencentcloud/services/ssm/action_tc_ssm_service.go` — New action implementation (factory, model, schema, Invoke).
- `tencentcloud/services/ssm/service_tencentcloud_ssm.go` — New `CreateService` and retry-wrapped `GetServiceStatus` service methods.
- `tencentcloud/services/ssm/action_tc_ssm_service_test.go` — New gomonkey unit tests.
- `tencentcloud/services/ssm/action_tc_ssm_service.md` — New action documentation.
- `tencentcloud/framework/registry.go` — Register `ssm.NewSsmServiceAction` in `actionFactories`.

### API References
- **SSM API** (`ssm/v20190923`):
  - `CreateService` (`CreateServiceWithContext`): opens the SSM service. Request fields: `PayMode *string` (PostPay only, case-sensitive required), `ClientToken *string` (client-generated idempotent token, required, 1–64 visible ASCII chars). Response fields: `ResourceId *string` (SSM service resource id, e.g. `ssm_<AppId>`), `OrderId *string` (billing order id). The API doc notes the returned resource may still be provisioning — the call is asynchronous.
  - `GetServiceStatus` (`GetServiceStatusWithContext`): returns the SSM service activation status. Request: no parameters. Response fields: `ServiceEnabled *bool` (true = opened), `InvalidType *int64`, `AccessKeyEscrowEnabled *bool`, `ExpireTime *string`, `QPSLimit *int64`, `SecretLimit *int64`, `PayModel *string`, `RenewFlag *int64`, `ResourceId *string`, `TotalCount *int64`, `ResourceRegion *int64`.

### Backward Compatibility
- **Fully Compatible**: This is a pure addition of a new action resource. No existing resources, data sources, or schemas are modified. No state migration is required.