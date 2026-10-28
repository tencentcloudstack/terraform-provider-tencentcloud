## Context

The SSM (Secrets Manager) cloud API exposes two account-level operations in the vendored SDK `ssm/v20190923` (both verified present — present in `vendor/.../ssm/v20190923/`, no SDK upgrade needed):

- `CreateService` (`CreateServiceWithContext`) — activates (opens) the SSM service for the calling account. Request fields: `PayMode *string` (PostPay only, case-sensitive, required), `ClientToken *string` (client-generated idempotent token, required, 1–64 visible ASCII chars, no spaces/Chinese/control chars; recommended `ssm-<UUID>`). Response fields: `ResourceId *string` (SSM service resource id, e.g. `ssm_<AppId>`), `OrderId *string` (billing order id). Per the API doc the returned resource "可能仍在发货，不代表服务已可用" — the call is **asynchronous**; activation is not guaranteed to have taken effect when the call returns.
- `GetServiceStatus` (`GetServiceStatusWithContext`) — returns the SSM service activation status. Request: no parameters. Response fields: `ServiceEnabled *bool` (true = opened), `InvalidType *int64`, `AccessKeyEscrowEnabled *bool`, `ExpireTime *string`, `QPSLimit *int64`, `SecretLimit *int64`, `PayModel *string`, `RenewFlag *int64`, `ResourceId *string`, `TotalCount *int64`, `ResourceRegion *int64`.

This is a **one-time operation** (RESOURCE_KIND_OPERATION): it triggers service activation; `Read`/`Delete` are no-ops. The implementation uses the SDKv2 operation resource pattern following `resource_tc_dnspod_modify_record_group_operation.go` and `resource_tc_sts_assume_role_operation.go`.

The user requirement states `client_token` shall NOT be exposed as a schema field — it is generated internally via `helper.BuildToken()`.

### Vendored cloud API struct definitions (excerpted for implementation reference)

From `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923/models.go`:

`CreateServiceRequestParams` / `CreateServiceRequest`:
```go
type CreateServiceRequestParams struct {
    PayMode     *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`
    ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}
type CreateServiceRequest struct {
    *tchttp.BaseRequest
    PayMode     *string `json:"PayMode,omitnil,omitempty" name:"PayMode"`
    ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`
}
```

`CreateServiceResponseParams` / `CreateServiceResponse`:
```go
type CreateServiceResponseParams struct {
    ResourceId *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`
    OrderId    *string `json:"OrderId,omitnil,omitempty" name:"OrderId"`
    RequestId  *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateServiceResponse struct {
    *tchttp.BaseResponse
    Response *CreateServiceResponseParams `json:"Response"`
}
```

`GetServiceStatusRequestParams` / `GetServiceStatusRequest`: **empty** (no fields).

`GetServiceStatusResponseParams` / `GetServiceStatusResponse`:
```go
type GetServiceStatusResponseParams struct {
    ServiceEnabled         *bool   `json:"ServiceEnabled,omitnil,omitempty" name:"ServiceEnabled"`
    InvalidType            *int64  `json:"InvalidType,omitnil,omitempty" name:"InvalidType"`
    AccessKeyEscrowEnabled *bool   `json:"AccessKeyEscrowEnabled,omitnil,omitempty" name:"AccessKeyEscrowEnabled"`
    ExpireTime             *string `json:"ExpireTime,omitnil,omitempty" name:"ExpireTime"`
    QPSLimit               *int64  `json:"QPSLimit,omitnil,omitempty" name:"QPSLimit"`
    SecretLimit            *int64  `json:"SecretLimit,omitnil,omitempty" name:"SecretLimit"`
    PayModel               *string `json:"PayModel,omitnil,omitempty" name:"PayModel"`
    RenewFlag              *int64  `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`
    ResourceId             *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`
    TotalCount            *int64  `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    ResourceRegion         *int64  `json:"ResourceRegion,omitnil,omitempty" name:"ResourceRegion"`
    RequestId              *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

Client methods from `vendor/.../ssm/v20190923/client.go`:
```go
func (c *Client) CreateServiceWithContext(ctx context.Context, request *CreateServiceRequest) (response *CreateServiceResponse, err error)
func (c *Client) GetServiceStatusWithContext(ctx context.Context, request *GetServiceStatusRequest) (response *GetServiceStatusResponse, err error)
```

### Current SSM service layer state

`tencentcloud/services/ssm/service_tencentcloud_ssm.go`:
- `SsmService` struct holds `client *connectivity.TencentCloudClient`; client accessor `me.client.UseSsmClient()` returns `*ssm.Client`.
- `NewSsmService(client *connectivity.TencentCloudClient) SsmService` is the constructor used by SDKv2 resources.
- Existing `DescribeSsmServiceStatusByFilter(ctx)` wraps `GetServiceStatus` WITHOUT a retry block (returns `*ssm.GetServiceStatusResponseParams`). We add a NEW retry-wrapped method `GetServiceStatus` so the operation resource's polling loop benefits from `tccommon.ReadRetryTimeout`-bounded retry and `tccommon.RetryError` wrapping.

### Reference operation resource structure (key structure, from `resource_tc_dnspod_modify_record_group_operation.go` & `resource_tc_sts_assume_role_operation.go`)

```go
package <product>

func ResourceTencentCloudXxxYyyOperation() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudXxxYyyOperationCreate,
        Read:   resourceTencentCloudXxxYyyOperationRead,
        Delete: resourceTencentCloudXxxYyyOperationDelete,
        Schema: map[string]*schema.Schema{
            "input_foo": { Required: true, ForceNew: true, Type: schema.TypeString, Description: "..." },
            // ...
        },
    }
}

func resourceTencentCloudXxxYyyOperationCreate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_xxx_yyy_operation.create")()
    defer tccommon.InconsistentCheck(d, meta)()
    logId := tccommon.GetLogId(tccommon.ContextNil)
    ctx := context.Background()
    // ... read inputs, build service, call API, set id, return Read ...
}

func resourceTencentCloudXxxYyyOperationRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_xxx_yyy_operation.read")()
    defer tccommon.InconsistentCheck(d, meta)()
    return nil  // no-op
}

func resourceTencentCloudXxxYyyOperationDelete(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_xxx_yyy_operation.delete")()
    defer tccommon.InconsistentCheck(d, meta)()
    return nil  // no-op
}
```

## Goals / Non-Goals

**Goals:**
- Add an SDKv2 operation resource `tencentcloud_ssm_service_operation` whose `Create` reads `pay_mode`, generates `client_token` internally, calls SSM `CreateService` to open the service, sets the id, then polls `GetServiceStatus` until `ServiceEnabled == true` to confirm the asynchronous activation took effect.
- Follow the SDKv2 operation resource references precisely.
- Wrap `CreateService` with `tccommon.WriteRetryTimeout` retry and `GetServiceStatus` with `tccommon.ReadRetryTimeout` retry, both with `resource.Retry` + `tccommon.RetryError(e)` and `ratelimit.Check(request.GetAction())`, in the service layer.
- Poll `GetServiceStatus` in a bounded loop until `ServiceEnabled == true`; cap the loop by `tccommon.ReadRetryTimeout`.
- Register the resource in `tencentcloud/provider.go`'s `ResourcesMap`; remove the prior framework action from `registry.go`.
- Provide gomonkey-based unit tests covering: successful create, polling across not-yet-enabled, CreateService error, Read/Delete no-op, schema validation.
- Provide a `.md` documentation file following the reference operation resource docs format.

**Non-Goals:**
- Do NOT expose `client_token` as a schema field (generated internally).
- Do NOT implement a Read/Update/Delete that calls cloud APIs (one-time operation; Read/Delete are no-ops).
- Do NOT modify any existing SSM resource, data source, or schema.

## Decisions

### 1. Implementation Pattern: SDKv2 Operation Resource
**Decision**: Implement `tencentcloud_ssm_service_operation` as an SDKv2 operation resource (`RESOURCE_KIND_OPERATION`), following `resource_tc_dnspod_modify_record_group_operation.go` / `resource_tc_sts_assume_role_operation.go`.

**Rationale**:
- The resource type is RESOURCE_KIND_OPERATION (one-time operation, Create-only, RUD no-ops).
- Matches the user requirement for this resource kind.

### 2. client_token Handling
**Decision**: Generate `clientToken` internally via `helper.BuildToken()`; do NOT expose it as a schema field.

**Rationale**:
- User requirement: "不暴露clientToken字段，直接使用clientToken := helper.BuildToken()即可".

### 3. Schema Definition
**Decision**: The schema defines exactly one required `StringAttribute`: `pay_mode` (ForceNew). It maps to `CreateServiceRequest.PayMode`.

### 4. Async handling: poll GetServiceStatus after CreateService; set id before polling
**Decision**: `CreateService` is asynchronous. After `CreateService` returns success, set `d.SetId(helper.BuildToken())` BEFORE polling, then poll `GetServiceStatus` until `ServiceEnabled == true`, bounded by `tccommon.ReadRetryTimeout`.

**Rationale**:
- Rule: for async resources, set id before polling so state retains it even if polling fails.
- `GetServiceStatus.ServiceEnabled` is the authoritative "is the service open" boolean.

### 5. Service Layer Method Signatures
**Decision**: Add to `service_tencentcloud_ssm.go`:
```go
func (me *SsmService) CreateService(ctx context.Context, payMode, clientToken string) (response *ssm.CreateServiceResponse, errRet error)
func (me *SsmService) GetServiceStatus(ctx context.Context) (response *ssm.GetServiceStatusResponseParams, errRet error)
```

### 6. Resource Registration
**Decision**: Register in `tencentcloud/provider.go` `ResourcesMap`. Remove `ssm.NewSsmServiceAction` from `registry.go` `actionFactories`.

### 7. Naming
- Type name: `ResourceTencentCloudSsmServiceOperation` (factory)
- Resource name: `tencentcloud_ssm_service_operation`
- Files: `resource_tc_ssm_service_operation.go`, `resource_tc_ssm_service_operation_test.go`, `resource_tc_ssm_service_operation.md`
- Log messages use `ssm_service_operation` (snake_case)

## Risks / Trade-offs

- **[Risk] Users expect `resource_id`/`order_id`/`service_enabled` as outputs** → This is a one-time operation resource; the API responses are logged for observability. No output attributes are exposed.
- **[Risk] `CreateService` may fail if SSM already opened** → Surfaces as a create error.
- **[Risk] Polling `GetServiceStatus` may exceed `ReadRetryTimeout`** → Bounded failure; user can re-run.

## Migration Plan

Purely additive; no migration. Rollback = revert the additive files and the `ResourcesMap`/`actionFactories` line changes.