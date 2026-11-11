## Context

Tencent Cloud Config supports 账号组 (Aggregate / Account Group) to manage resources across member accounts. Each account group has its own 投递设置 (delivery settings) that controls how Config audit logs are delivered to COS or CLS. The cloud API provides two relevant endpoints:

- `DescribeAggregateConfigDeliver` — reads the delivery settings of an account group (takes `AccountGroupId`).
- `UpdateAggregateConfigDeliver` — edits the delivery settings of an account group (takes `AccountGroupId` + delivery fields).

There is no Create or Delete API — the delivery config exists per account group and is managed via Read + Update only. This maps to RESOURCE_KIND_CONFIG: the resource exists as long as the account group exists, and Terraform manages reading and updating the configuration.

The existing sibling resource `tencentcloud_config_deliver_config` (non-aggregate, global singleton) already implements the same pattern against `DescribeConfigDeliver`/`UpdateConfigDeliver`. The new aggregate resource follows the same lifecycle but is keyed by `account_group_id` instead of a synthetic token.

## Goals / Non-Goals

**Goals:**
- Provide Terraform RU lifecycle management for the aggregate delivery config.
- Register the resource in provider.go and provider.md.
- Deliver design.md with enough vendored API reference so the implementer need not re-read vendor models.go.

**Non-Goals:**
- No account group lifecycle management (Create/Delete of account groups themselves).
- No support for non-aggregate delivery config (already exists separately).
- No async polling — both APIs are synchronous.

## Decisions

### 1. Resource ID = account_group_id

There is exactly one delivery config per account group, and `DescribeAggregateConfigDeliver` / `UpdateAggregateConfigDeliver` both take `AccountGroupId` as the sole key. Therefore `account_group_id` is the natural resource ID. This differs from the non-aggregate singleton resource which uses `helper.BuildToken()`. Using the real account group id as `d.Id()` allows the resource to be imported and re-read.

### 2. Create = Update (RU pattern, no Create API)

There is no Create API. `terraform apply` on a new resource calls `UpdateAggregateConfigDeliver` to set the delivery config, then `d.SetId(account_group_id)`. This follows the existing `tencentcloud_config_deliver_config` pattern where Create delegates to Update. Note: for the aggregate resource we set the id to the account_group_id (not BuildToken) because the account group id is the real key and is available before the call.

### 3. Delete = no-op

No Delete API exists. Removing the resource from state should not delete the account group's delivery config. Delete handler returns nil. (If the user wants to disable delivery they set `status = 0`.)

### 4. Schema field types match cloud API types

| Schema field | Type in SDK | Terraform schema type |
|---|---|---|
| account_group_id | `*string` | TypeString |
| status | `*uint64` | TypeInt |
| deliver_name | `*string` | TypeString |
| target_arn | `*string` | TypeString |
| deliver_prefix | `*string` | TypeString |
| deliver_type | `*string` | TypeString |
| deliver_uin | `*int64` | TypeInt |
| deliver_content_type | `*uint64` | TypeInt |
| create_time | `*string` (response only) | TypeString (Computed) |

### 5. Read handler must guard nil before set

Per project rules, before each `d.Set(...)` the handler checks the response field is non-nil. If the response (or its inner Response) is nil, log `[CRUD]` with the id then `d.SetId("")`.

### 6. Update handler retries via tccommon.WriteRetryTimeout

The Update call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` on error. After the retry block completes, read is called to refresh state. `account_group_id` is read from `d.Id()`.

### 7. No importer needed for RU CONFIG resource

Wait — actually `account_group_id` is a meaningful key, so an importer is reasonable (matching the existing config_deliver_config which had Importer). However the resource file naming requested is `resource_tc_config_update_aggregate_config_deliver.go`. The import would pass the account_group_id as the id. We include Importer support (ImportStatePassthrough).

## Risks / Trade-offs

- **[No Create/Delete API]** → Resource uses RU pattern; Create delegates to Update, Delete is no-op. State removal does not clear the cloud delivery config. Documented in the .md.
- **[Read returns null fields on new configs]** → Response fields are marked `omitnil` and may return null on cloud side before the first Update. Read handler guards each field with nil checks before `d.Set`.
- **[account_group_id ForceNew]** → `account_group_id` is the resource key and MUST be ForceNew (changing it means a different account group's delivery config). All other top-level fields are mutable via Update.

## Cloud API Reference (vendored from vendor/.../config/v20220802/models.go)

### DescribeAggregateConfigDeliver

**Request** (`models.go:1878`):
```go
type DescribeAggregateConfigDeliverRequest struct {
	*tchttp.BaseRequest
	// 账号组ID
	AccountGroupId *string `json:"AccountGroupId,omitnil,omitempty" name:"AccountGroupId"`
}
```

**Response params** (`models.go:1905` → `DescribeAggregateConfigDeliverResponseParams`, embedded in `DescribeAggregateConfigDeliverResponse.Response`):
```go
type DescribeAggregateConfigDeliverResponseParams struct {
	DeliverName       *string `json:"DeliverName,omitnil,omitempty"`       // 投递名称
	TargetArn         *string `json:"TargetArn,omitnil,omitempty"`         // 资源六段式
	Status            *uint64 `json:"Status,omitnil,omitempty"`            // 0 关闭  1 开启
	CreateTime        *string `json:"CreateTime,omitnil,omitempty"`        // 创建时间
	DeliverPrefix     *string `json:"DeliverPrefix,omitnil,omitempty"`    // 日志前缀
	DeliverType       *string `json:"DeliverType,omitnil,omitempty"`      // 投递类型
	DeliverUin        *int64  `json:"DeliverUin,omitnil,omitempty"`        // 跨账号投递成员uin, 默认0
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty"`// 1:配置变更 2:资源列表 3:全部
	RequestId         *string `json:"RequestId,omitnil,omitempty"`
}
type DescribeAggregateConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAggregateConfigDeliverResponseParams `json:"Response"`
}
```

Client method (client.go:1089): `func (c *Client) DescribeAggregateConfigDeliverWithContext(ctx, request) (*DescribeAggregateConfigDeliverResponse, error)`; convenience `DescribeAggregateConfigDeliver(request)`.

### UpdateAggregateConfigDeliver

**Request** (`models.go:5229`):
```go
type UpdateAggregateConfigDeliverRequest struct {
	*tchttp.BaseRequest
	Status             *uint64 `json:"Status,omitnil,omitempty"`             // 0 关闭  1 开启
	AccountGroupId     *string `json:"AccountGroupId,omitnil,omitempty"`     // 账号组ID
	DeliverName        *string `json:"DeliverName,omitnil,omitempty"`        // 投递服务名称
	TargetArn          *string `json:"TargetArn,omitnil,omitempty"`          // 资源六段式 (COS/CLS)
	DeliverPrefix      *string `json:"DeliverPrefix,omitnil,omitempty"`      // 资源前缀
	DeliverType        *string `json:"DeliverType,omitnil,omitempty"`        // 投递类型 COS CLS
	DeliverUin         *int64  `json:"DeliverUin,omitnil,omitempty"`         // 跨账号投递成员uin, 默认0
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty"` // 1:配置变更 2:资源列表 3:全选
}
```

**Response** (`models.go:5287`): only `RequestId` (no data).
```go
type UpdateAggregateConfigDeliverResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty"`
}
type UpdateAggregateConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *UpdateAggregateConfigDeliverResponseParams `json:"Response"`
}
```

Client method (client.go:3163): `func (c *Client) UpdateAggregateConfigDeliverWithContext(ctx, request) (*UpdateAggregateConfigDeliverResponse, error)`; convenience `UpdateAggregateConfigDeliver(request)`.

SDK client accessor: `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client()`.

## Code Style Reference (sibling resource resource_tc_config_deliver_config.go)

The existing non-aggregate resource `tencentcloud/services/config/resource_tc_config_deliver_config.go` is the closest sibling. Key structural points to reuse for the aggregate resource:

- **Schema builder** `ResourceTencentCloudConfigDeliverConfig()` returns `*schema.Resource` with Create/Read/Update/Delete handlers and an `Importer` (ImportStatePassthrough).
- **Create handler**: `d.SetId(helper.BuildToken())` then `return resourceTencentCloudConfigDeliverConfigUpdate(d, meta)`.
  - **For the aggregate resource**: Create handler should set `d.SetId(account_group_id)` (read from `d.Get("account_group_id")`) then delegate to Update. After Update success, call Read to refresh. Since `account_group_id` is the real key and ForceNew, it is available in Create.
- **Read handler**: builds service = `ConfigService{client: ...}`, calls service method, guards `respData == nil` → log + `d.SetId("")`; then for each field checks `!= nil` before `d.Set`.
- **Update handler**: builds `configv20220802.NewUpdateConfigDeliverRequest()`, populates from `d.GetOk`/`d.GetOkExists`, wraps call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)`, after retry returns nil calls Read.
- **Delete handler**: no-op, returns nil.
- **Service method** `DescribeConfigDeliver(ctx)` at `service_tencentcloud_config.go:579`: builds request, wraps in `resource.Retry(tccommon.ReadRetryTimeout, ...)` checking `result == nil || result.Response == nil` → `NonRetryableError`, returns `response.Response` (the params struct).
  - **For aggregate**: method signature `DescribeAggregateConfigDeliver(ctx, accountGroupId string) (*configv20220802.DescribeAggregateConfigDeliverResponseParams, error)` — sets `request.AccountGroupId`.
- **provider.go registration**: line 2675 `"tencentcloud_config_deliver_config": config.ResourceTencentCloudConfigDeliverConfig(),` — add the new resource below it.
- **provider.md**: line 2779 lists `tencentcloud_config_deliver_config` — add the new resource name below it.

### Notes on deliver_uin

`deliver_uin` is `*int64` in the SDK. The existing non-aggregate resource did NOT expose `deliver_uin` in its schema. The aggregate resource MUST expose `deliver_uin` per the requirement (it is in the Update request and Describe response). Use `helper.Int64()` for int64 and `helper.IntUint64()` for uint64.

### Import

Since `account_group_id` is the real key and is ForceNew, the resource supports import. The .md import example should state: `terraform import tencentcloud_config_update_aggregate_config_deliver.example <account_group_id>`.