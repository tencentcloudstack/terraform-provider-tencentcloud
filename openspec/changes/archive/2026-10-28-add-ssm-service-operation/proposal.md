## Why

The SSM (Secrets Manager) cloud API exposes a `CreateService` operation that activates (opens) the SSM service (PostPay / pay-as-you-go) for the current account. This is a one-time activation step that practitioners currently cannot trigger through Terraform — they must open the service manually in the console before managing any SSM secret resources. Adding an **operation resource** lets users open SSM as part of their IaC workflow alongside other SSM resources.

## What Changes

- Add a new SDKv2 operation resource `tencentcloud_ssm_service_operation` implemented via the `RESOURCE_KIND_OPERATION` pattern (one-time operation), following the reference implementations `resource_tc_dnspod_modify_record_group_operation.go` and `resource_tc_sts_assume_role_operation.go`. The resource exposes only `Create`/`Read`/`Delete` (no `Update`).
- Implement the `Create` handler to:
  - Read the required input `pay_mode` (`request.PayMode`) from the resource data.
  - Generate `client_token` (`request.ClientToken`) internally via `helper.BuildToken()` — it is NOT exposed as a schema field (per requirement).
  - Call the SSM `CreateService` API (`CreateServiceWithContext`) via the service layer to open the service.
  - Because `CreateService` is **asynchronous** (the API doc states "返回时资源可能仍在发货，不代表服务已可用" — the resource may still be provisioning when the call returns), set `d.SetId(helper.BuildToken())` BEFORE polling, then poll the SSM `GetServiceStatus` API until `ServiceEnabled == true` (with `tccommon.ReadRetryTimeout`-bounded retry loop) to confirm activation has taken effect.
- `Read` and `Delete` are no-ops (return nil) — this is a one-time operation resource; no cloud-side state is read or destroyed.
- Add a `CreateService` method and a `GetServiceStatus` polling helper to the SSM service layer (`service_tencentcloud_ssm.go`); the existing `DescribeSsmServiceStatusByFilter` already wraps `GetServiceStatus` but without retry, so a dedicated retry-wrapped method is added for the operation resource. Wrap both API calls with `tccommon.WriteRetryTimeout`/`tccommon.ReadRetryTimeout` retry and `tccommon.RetryError` error wrapping.
- Register the resource `tencentcloud_ssm_service_operation` → `ssm.ResourceTencentCloudSsmServiceOperation()` in `tencentcloud/provider.go`'s `ResourcesMap`. Remove the previous framework action registration `ssm.NewSsmServiceAction` from `tencentcloud/framework/registry.go`'s `actionFactories` slice (this change replaces the earlier framework-action approach with the SDKv2 operation resource approach).
- Provide a markdown documentation file `resource_tc_ssm_service_operation.md` with a one-line description mentioning SSM (Secrets Manager) and an Example Usage section using the SDKv2 resource block syntax.
- Add unit tests in `resource_tc_ssm_service_operation_test.go` using gomonkey mocks (no Terraform acceptance test suite), covering successful create (with polling), polling-across-not-yet-enabled, CreateService API error, Read/Delete no-op, and schema validation.

## Capabilities

### New Capabilities
- `ssm-service-action`: Adds an SDKv2 operation resource `tencentcloud_ssm_service_operation` that activates the SSM service (PostPay) via the SSM `CreateService` API and polls `GetServiceStatus` until the service is enabled. This is a one-time operation resource; `Read`/`Delete` are no-ops. The `client_token` is generated internally via `helper.BuildToken()` and is not exposed as a schema field.

### Modified Capabilities
<!-- None: brand-new operation resource, no existing requirements change. -->

## Impact

### Affected Code
- `tencentcloud/services/ssm/resource_tc_ssm_service_operation.go` — New operation resource implementation (factory, schema, Create/Read/Delete).
- `tencentcloud/services/ssm/service_tencentcloud_ssm.go` — New `CreateService` and retry-wrapped `GetServiceStatus` service methods.
- `tencentcloud/services/ssm/resource_tc_ssm_service_operation_test.go` — New gomonkey unit tests.
- `tencentcloud/services/ssm/resource_tc_ssm_service_operation.md` — New resource documentation.
- `tencentcloud/provider.go` — Register `tencentcloud_ssm_service_operation` in `ResourcesMap`.
- `tencentcloud/framework/registry.go` — Remove `ssm.NewSsmServiceAction` from `actionFactories` (replaced by SDKv2 operation resource).

### API References
- **SSM API** (`ssm/v20190923`):
  - `CreateService` (`CreateServiceWithContext`): opens the SSM service. Request fields: `PayMode *string` (PostPay only, case-sensitive required), `ClientToken *string` (client-generated idempotent token, required, 1–64 visible ASCII chars). Response fields: `ResourceId *string` (SSM service resource id, e.g. `ssm_<AppId>`), `OrderId *string` (billing order id). The API doc notes the returned resource may still be provisioning — the call is asynchronous.
  - `GetServiceStatus` (`GetServiceStatusWithContext`): returns the SSM service activation status. Request: no parameters. Response fields: `ServiceEnabled *bool` (true = opened), `InvalidType *int64`, `AccessKeyEscrowEnabled *bool`, `ExpireTime *string`, `QPSLimit *int64`, `SecretLimit *int64`, `PayModel *string`, `RenewFlag *int64`, `ResourceId *string`, `TotalCount *int64`, `ResourceRegion *int64`.

### Backward Compatibility
- **Fully Compatible**: This is a pure addition of a new operation resource. No existing resources, data sources, or schemas are modified. No state migration is required.