## Context

The SSM (Secrets Manager) cloud API exposes two account-level operations in the vendored SDK `ssm/v20190923` (both verified present — present in `vendor/.../ssm/v20190923/`, no SDK upgrade needed):

- `CreateService` (`CreateServiceWithContext`) — activates (opens) the SSM service for the calling account. Request fields: `PayMode *string` (PostPay only, case-sensitive, required), `ClientToken *string` (client-generated idempotent token, required, 1–64 visible ASCII chars, no spaces/Chinese/control chars; recommended `ssm-<UUID>`). Response fields: `ResourceId *string` (SSM service resource id, e.g. `ssm_<AppId>`), `OrderId *string` (billing order id). Per the API doc the returned resource "可能仍在发货，不代表服务已可用" — the call is **asynchronous**; activation is not guaranteed to have taken effect when the call returns.
- `GetServiceStatus` (`GetServiceStatusWithContext`) — returns the SSM service activation status. Request: no parameters. Response fields: `ServiceEnabled *bool` (true = opened), `InvalidType *int64`, `AccessKeyEscrowEnabled *bool`, `ExpireTime *string`, `QPSLimit *int64`, `SecretLimit *int64`, `PayModel *string`, `RenewFlag *int64`, `ResourceId *string`, `TotalCount *int64`, `ResourceRegion *int64`.

This is a **one-time operation** (RESOURCE_KIND_OPERATION): it triggers service activation, no cloud-side state needs to be persisted after the action completes.

The user has explicitly requested that this be implemented with the **Terraform Plugin Framework `action.Action`** pattern (not the SDKv2 operation resource pattern), following the existing reference actions `NewDlcInitializeTCLake` (`tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.go`, parameterless) and `NewAntiddosUnblockResources` (`tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.go`, takes input arguments).

### Vendored cloud API struct definitions (excerpted for implementation reference)

From `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923/models.go`:

`CreateServiceRequestParams` (lines 449–455) / `CreateServiceRequest` (lines 457–465):
```go
type CreateServiceRequestParams struct {
    // 只支持 PostPay，区分大小写。枚举值：PostPay：后付费 按量计费模式
    PayMode *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`
    // 客户端生成的创建请求标识。同一次逻辑创建发生重试时须使用相同的值。
    // 格式：1～64 个可见 ASCII 字符（! 至 ~）。必填；不能含空格、中文或控制字符。建议 ssm-<UUID>。
    // SSM 透传给计费侧，不存储或返回。订单状态不明时先核单，再决定是否用原值重试。
    ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}
type CreateServiceRequest struct {
    *tchttp.BaseRequest
    PayMode     *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`
    ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}
```

`CreateServiceResponseParams` (lines 488–497) / `CreateServiceResponse` (lines 499–502):
```go
type CreateServiceResponseParams struct {
    // SSM 服务资源 ID，供调用方保存并标识该服务资源。格式：一期后付费为 ssm_ 加 AppId，如 ssm_1250000000。
    // 返回时资源可能仍在发货，不代表服务已可用。
    ResourceId *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`
    // 计费交易订单号，用于核对本次开通订单
    OrderId *string `json:"OrderId,omitnil,omitempty" name:"OrderId"`
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateServiceResponse struct {
    *tchttp.BaseResponse
    Response *CreateServiceResponseParams `json:"Response"`
}
```

`GetServiceStatusRequestParams` (lines 1379–1381) / `GetServiceStatusRequest` (lines 1383–1386): **empty** (no fields).

`GetServiceStatusResponseParams` (lines 1408–1444) / `GetServiceStatusResponse` (lines 1446–1449):
```go
type GetServiceStatusResponseParams struct {
    ServiceEnabled         *bool   `json:"ServiceEnabled,omitnil,omitempty" name:"ServiceEnabled"`         // true=服务已开通
    InvalidType            *int64  `json:"InvalidType,omitnil,omitempty" name:"InvalidType"`              // 0未购买 1正常 2欠费停服 3资源释放
    AccessKeyEscrowEnabled *bool   `json:"AccessKeyEscrowEnabled,omitnil,omitempty" name:"AccessKeyEscrowEnabled"`
    ExpireTime             *string `json:"ExpireTime,omitnil,omitempty" name:"ExpireTime"`
    QPSLimit               *int64  `json:"QPSLimit,omitnil,omitempty" name:"QPSLimit"`
    SecretLimit            *int64  `json:"SecretLimit,omitnil,omitempty" name:"SecretLimit"`
    PayModel               *string `json:"PayModel,omitnil,omitempty" name:"PayModel"`
    RenewFlag              *int64  `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`                      // 0手动 1自动 2到期不续
    ResourceId             *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`
    TotalCount             *int64  `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    ResourceRegion         *int64  `json:"ResourceRegion,omitnil,omitempty" name:"ResourceRegion"`
    RequestId              *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

Client methods from `vendor/.../ssm/v20190923/client.go`:
```go
func (c *Client) CreateServiceWithContext(ctx context.Context, request *CreateServiceRequest) (response *CreateServiceResponse, err error)
//   InitBaseRequest("ssm", APIVersion, "CreateService"); requires credential.
func (c *Client) GetServiceStatusWithContext(ctx context.Context, request *GetServiceStatusRequest) (response *GetServiceStatusResponse, err error)
//   InitBaseRequest("ssm", APIVersion, "GetServiceStatus"); requires credential.
```

### Current SSM service layer state

`tencentcloud/services/ssm/service_tencentcloud_ssm.go`:
- `SsmService` struct holds `client *connectivity.TencentCloudClient`; client accessor `me.client.UseSsmClient()` returns `*ssm.Client`.
- `NewSsmService(client *connectivity.TencentCloudClient) SsmService` is the constructor used by BOTH SDKv2 resources and framework resources (`resource_tc_ssm_secret_version_v2.go` calls `NewSsmService(r.Client())`).
- Existing `DescribeSsmServiceStatusByFilter(ctx)` wraps `GetServiceStatus` WITHOUT a retry block (returns `*ssm.GetServiceStatusResponseParams`). We add a NEW retry-wrapped method `GetServiceStatus` so the action's polling loop benefits from `tccommon.ReadRetryTimeout`-bounded retry and `tccommon.RetryError` wrapping.

### Reference action implementation structure (key structure, from `action_tc_dlc_initialize_tc_lake.go` & `action_tc_antiddos_unblock_resources.go`)

```go
package <product>

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	... sdk import ...
)

var _ action.ActionWithConfigure = &XxxYyyAction{}

func NewXxxYyyAction() action.Action { return &XxxYyyAction{} }

type XxxYyyAction struct { fw.ActionWithConfigure }

type XxxYyyActionModel struct {
	// inputs as types.String etc.
}

func (a *XxxYyyAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_xxx_yyy"
}

func (a *XxxYyyAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "...",
		Attributes: map[string]schema.Attribute{
			"input_foo": schema.StringAttribute{ Required: true, Description: "..." },
			...
		},
	}
}

func (a *XxxYyyAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data XxxYyyActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() { return }
	// validate inputs (null/unknown/empty guards) -> AddError, return
	if a.Client() == nil { resp.Diagnostics.AddError("Provider not configured", "..."); return }
	service := NewXxxService(a.Client())
	... call API; on err AddError; on success return (no state) ...
}
```

Key framework pieces:
- `fw.ActionWithConfigure` (`tencentcloud/framework/fw/action_with_configure.go`) embeds `withMeta` and provides `Client()` accessor + `Configure` method that reads `*sharedmeta.ProviderMeta`.
- Registration is ONLY in `tencentcloud/framework/registry.go`'s `actionFactories` slice (NOT in SDKv2 `provider.go`). `frameworkActions()` → `Provider.Actions()` exposes them.
- Framework action test pattern (`action_tc_antiddos_unblock_resources_test.go`): build `action.InvokeRequest` with raw `tftypes.Value` config matching the schema; inject mock client via `action.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}`; mock the SDK client method via gomonkey `ApplyMethodFunc`/`ApplyMethodReturn`.

### Framework Action Constraint (Critical)

The Terraform Plugin Framework `action.Action` schema attributes **cannot be Computed**. The `action.InvokeResponse` struct only carries `Diagnostics` and a `SendProgress` function; it has no mechanism to write output state back to Terraform. Consequently, even though `CreateService` returns `ResourceId`/`OrderId` and `GetServiceStatus` returns many status fields, NONE can be exposed as Terraform output attributes. The reference implementations confirm this: they invoke an API and surface only success/error diagnostics. All response fields are logged via `log.Printf` for observability but not persisted as Terraform state.

## Goals / Non-Goals

**Goals:**
- Add a framework action `tencentcloud_ssm_service` whose `Invoke` collects required inputs `pay_mode` and `client_token`, calls SSM `CreateService` to open the service, then polls `GetServiceStatus` until `ServiceEnabled == true` to confirm the asynchronous activation took effect.
- Follow the reference actions (`NewDlcInitializeTCLake`, `NewAntiddosUnblockResources`) precisely for the action struct, factory, model, metadata, schema, and invoke flow.
- Wrap `CreateService` with `tccommon.WriteRetryTimeout` retry and `GetServiceStatus` with `tccommon.ReadRetryTimeout` retry, both with `resource.Retry` + `tccommon.RetryError(e)` and `ratelimit.Check(request.GetAction())`, in the service layer.
- Poll `GetServiceStatus` in a bounded loop until `ServiceEnabled == true`; cap the loop by `tccommon.ReadRetryTimeout` so it does not run forever. Log the status each iteration and the final confirmed-enabled status.
- Register the action factory in `tencentcloud/framework/registry.go`'s `actionFactories` (the `ssm` services package is ALREADY imported there, so no new import is needed).
- Provide gomonkey-based unit tests (no Terraform acceptance test suite) covering: successful invoke (CreateService ok + polling reaches ServiceEnabled==true), missing `pay_mode`, missing `client_token`, `CreateService` API error, polling error (GetServiceStatus returns error throughout / never becomes enabled — bounded), client-not-configured, and schema/metadata validation — mirroring `action_tc_antiddos_unblock_resources_test.go`.
- Provide a `.md` documentation file following the reference action docs format.

**Non-Goals:**
- Do NOT implement this as an SDKv2 operation resource (`resource_tc_ssm_service_operation.go`) — the user explicitly requested the framework action pattern.
- Do NOT expose any cloud API response field (`resource_id`, `order_id`, `service_enabled`, etc.) as a Terraform attribute (framework actions do not support Computed/output attributes and `InvokeResponse` has no output mechanism). All are logged for observability only.
- Do NOT register the action in `tencentcloud/provider.go`'s `ResourcesMap` — framework actions are registered only in `tencentcloud/framework/registry.go`.
- Do NOT add a Read/Update/Delete lifecycle — actions only implement `Invoke`.
- Do NOT modify any existing SSM resource, data source, or schema.

## Decisions

### 1. Implementation Pattern: Framework Action (not SDKv2 operation)
**Decision**: Implement `tencentcloud_ssm_service` as a Terraform Plugin Framework `action.Action`, following `NewDlcInitializeTCLake` / `NewAntiddosUnblockResources`.

**Rationale**:
- The user explicitly requested the framework action pattern.
- Framework actions are the correct primitive for one-time operations that do not persist state, which matches SSM service-activation semantics exactly.
- `NewAntiddosUnblockResources` is the close reference because it takes required input arguments (`resources`) — our action takes `pay_mode` + `client_token`.

### 2. Output Attribute Handling
**Decision**: Do NOT expose `resource_id`/`order_id`/`service_enabled`/... as schema attributes. Log the `CreateService` response (`ResourceId`, `OrderId`) and the polled `GetServiceStatus` body via `log.Printf` inside the service methods (consistent with existing SSM service methods). The Invoke handler only acts on errors.

**Rationale**:
- Framework action schema attributes cannot be `Computed`.
- `action.InvokeResponse` only contains `Diagnostics` and `SendProgress` — there is no output/state write path.
- Both reference implementations invoke their API and surface only diagnostics; neither exposes response fields as attributes.
- The requirement's output mappings (`resource_id`, `order_id`, `service_enabled`, ...) are captured in the Go response structs and logged, fulfilling observability without violating framework constraints.

### 3. Schema Definition (two required input attributes)
**Decision**: The action schema defines exactly two required `StringAttribute`s: `pay_mode` (the value `PostPay`) and `client_token` (1–64 visible ASCII chars). These map to `CreateServiceRequest.PayMode` and `CreateServiceRequest.ClientToken`.

**Rationale**:
- The cloud API `CreateServiceRequestParams` has exactly these two fields, both required per the API doc (PayMode required/PostPay-only; ClientToken required idempotency token).
- `NewAntiddosUnblockResources` demonstrates the Required-attributes pattern; we mirror it with two `schema.StringAttribute{Required: true}`.

### 4. Async handling: poll GetServiceStatus after CreateService
**Decision**: `CreateService` is asynchronous (API doc: "返回时资源可能仍在发货，不代表服务已可用"). After `CreateService` returns success, the Invoke handler SHALL poll `GetServiceStatus` until `ServiceEnabled == true`. The polling is bounded by `tccommon.ReadRetryTimeout` (reusing `resource.Retry(ReadRetryTimeout, ...)`): on each poll call the retry-wrapped `GetServiceStatus`; if `ServiceEnabled != nil && *ServiceEnabled == true` → done; otherwise return `resource.RetryableError(...)` so the loop continues; if `ReadRetryTimeout` is exhausted the action fails with a diagnostic.

**Rationale**:
- User requirement: for asynchronous interfaces, poll the Read interface until it takes effect.
- `GetServiceStatus.ServiceEnabled` is the authoritative "is the service open" boolean.
- Reusing `resource.Retry(ReadRetryTimeout, ...)` gives a single, project-standard bounded wait with `tccommon.RetryError` wrapping for non-retryable errors and a `RetryableError` to keep polling while not-yet-enabled.

**Alternatives Considered**:
- Return immediately after `CreateService` and skip polling: rejected — the API doc explicitly warns the returned resource may still be provisioning; users would see a "success" while the service is not yet usable.
- Add a `timeouts` block (like mongodb restore): rejected for now — `CreateService` activation is fast and a single `ReadRetryTimeout`-bounded poll is sufficient; the timeouts block adds complexity not needed here. (Documented as an open option only if polling proves too short in practice.)

### 5. Service Layer Method Signatures
**Decision**: Add to `service_tencentcloud_ssm.go`:

```go
// CreateService opens the SSM service (PostPay). Returns the CreateServiceResponse
// so the Invoke handler can log ResourceId/OrderId.
func (me *SsmService) CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)

// GetServiceStatus returns the current SSM service status (retry-wrapped), used by
// the action's activation polling loop.
func (me *SsmService) GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)
```

**Rationale**:
- Both use `resource.Retry` + `tccommon.RetryError(e)` + `ratelimit.Check(request.GetAction())` and the `log.Printf("[DEBUG]...")` / deferred `log.Printf("[CRITAL]...")` pattern used by every existing SSM service method.
- `CreateService` uses `tccommon.WriteRetryTimeout`; `GetServiceStatus` uses `tccommon.ReadRetryTimeout`.
- Returns raw responses so the Invoke handler / polling loop can inspect `ServiceEnabled` and log fields. `GetServiceStatus` returns `*GetServiceStatusResponseParams` (the `Response` inner struct) for parity with the existing `DescribeSsmServiceStatusByFilter` and to simplify nil-safe field access.
- We deliberately do NOT add a separate `pay_mode`/`client_token` per-field param map — direct typed args are clearer, matching `RestoreDBInstance(ctx, instanceId, restoreTime, restoreDatabases)`.

### 6. Action Factory & Registration
**Decision**: Factory `NewSsmServiceAction() action.Action` returns `&SsmServiceAction{}`. Register `ssm.NewSsmServiceAction` in `tencentcloud/framework/registry.go`'s `actionFactories`. The `ssm` services package is ALREADY imported in `registry.go` (for the framework resource + ephemeral), so NO new import line is needed.

**Rationale**:
- Mirrors the exact registration pattern of the reference actions.
- Framework actions are NOT registered in SDKv2 `provider.go`.

### 7. Naming
**Decision**:
- Type name: `SsmServiceAction`
- Factory: `NewSsmServiceAction`
- Model type: `SsmServiceActionModel`
- Action type name (Metadata): `tencentcloud_ssm_service`
- Files: `action_tc_ssm_service.go`, `action_tc_ssm_service_test.go`, `action_tc_ssm_service.md`

**Rationale**:
- Follows the `NewDlcInitializeTCLake` / `NewAntiddosUnblockResources` naming convention (action files prefixed with `action_tc_`).
- The action type name matches the requested resource name `tencentcloud_ssm_service`.

### 8. Test Strategy
**Decision**: Use gomonkey mocks (no Terraform acceptance test suite), mirroring `action_tc_antiddos_unblock_resources_test.go`. Mock `UseSsmClient` (returns `*ssm.Client`), `CreateServiceWithContext`, and `GetServiceStatusWithContext`. Tests cover: success (CreateService ok + GetServiceStatus returns ServiceEnabled==true on the second poll), missing pay_mode, missing client_token, CreateService error, polling-never-enabled (bounded, asserts error diagnostic), client-not-configured, and metadata/schema validation.

## Risks / Trade-offs

- **[Risk] Users expect `resource_id`/`order_id`/`service_enabled` as outputs** → Mitigation: Document clearly in the `.md` and design that framework actions do not expose output attributes; the API responses are logged for observability. The action's value is triggering the one-time activation, not reading back state.
- **[Risk] `CreateService` may fail if the SSM service is already opened** (`ResourceInUse.ServiceAlreadyEnabled`) → Mitigation: The `tccommon.RetryError` attention maps transient errors as retryable; a non-retryable "already enabled" surfaces as a diagnostic. Since the goal is to ensure the service is enabled, an "already enabled" error may be treated as success (the desired end state is reached). Decision: surface it as a diagnostic (do not special-case) to keep behavior explicit; tracked as an Open Question if user-friendliness requires suppressing it later.
- **[Risk] Polling `GetServiceStatus` may exceed `ReadRetryTimeout` on slow provisioning** → Mitigation: `ReadRetryTimeout` bounds the wait and surfaces a clear failure; the user can re-run the action (it is idempotent-ish). Documented as a future `timeouts` option if needed.
- **[Risk] The `ssm` package import already exists in `registry.go`** → Mitigation: Tasks explicitly note no new import is needed, only the `actionFactories` append, avoiding a duplicate-import compile error.

## Migration Plan

Purely additive; no migration. New action registered alongside existing framework actions. Rollback = revert the additive files and the single `actionFactories` line.