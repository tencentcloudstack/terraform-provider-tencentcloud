## Context

Tencent Cloud Config (配置审计) allows delivering audit logs to COS buckets or CLS topics. The cloud exposes two relevant APIs:
- `DescribeConfigDeliver` — returns the global singleton delivery settings (no request params).
- `UpdateConfigDeliver` — edits the delivery settings (idempotent; used as both Create and Update).

There is no dedicated Create or Delete API. The existing `tencentcloud_config_deliver_config` resource already implements this same cloud behavior as a singleton config; this change adds a parallel RESOURCE_KIND_CONFIG resource named `tencentcloud_config_update_config_deliver` whose file naming follows the project convention `resource_tc_<product>_<name>_config.go` (here `<name>` = `update_config_deliver`). It uses the same Read/Update-only lifecycle.

Reference code style: `tencentcloud_igtm_strategy` (CRUD skeleton, retry blocks, helper usage) and `resource_tc_config_recorder_config.go` / `resource_tc_config_deliver_config.go` (singleton CONFIG pattern with `helper.BuildToken()` ID and no-op Delete).

## Goals / Non-Goals

**Goals:**
- Provide a RESOURCE_KIND_CONFIG resource `tencentcloud_config_update_config_deliver` that performs Read (via `DescribeConfigDeliver`) and Update (via `UpdateConfigDeliver`).
- Reuse the existing `ConfigService.DescribeConfigDeliver` service-layer method for Read.
- Register the resource in the provider and add documentation + unit tests.

**Non-Goals:**
- Do not modify or remove the existing `tencentcloud_config_deliver_config` resource.
- Do not add a true remote Delete (no such API exists); Delete remains a no-op.
- Do not introduce async polling — `UpdateConfigDeliver` is synchronous.

## Decisions

### Decision 1: Singleton ID via `helper.BuildToken()`
The delivery config is a global singleton with no natural unique key. The resource ID SHALL be `helper.BuildToken()`. This matches the existing `tencentcloud_config_deliver_config` and `tencentcloud_config_recorder_config` patterns.

### Decision 2: Create = Update
Create SHALL set `d.SetId(helper.BuildToken())` then delegate to the Update handler, which calls `UpdateConfigDeliver`. This is the established CONFIG-resource pattern (see `resource_tc_config_recorder_config.go` lines 68-75 and `resource_tc_config_deliver_config.go` lines 71-78).

### Decision 3: Reuse existing service-layer DescribeConfigDeliver
`ConfigService.DescribeConfigDeliver(ctx)` already exists at `service_tencentcloud_config.go:579` and returns `*configv20220802.DescribeConfigDeliverResponseParams`. The Read handler SHALL call it directly rather than adding a new service method.

### Decision 4: Delete is a no-op
No Delete API exists. Delete SHALL return nil. The remote delivery setting is not removed. Users who want to disable delivery can set `status = 0`.

### Decision 5: Update fields handling
Per project rules, the Update handler uses `GetOkExists` for int fields (`status`, `deliver_content_type`) so that `0` is honored, and `GetOk` for string fields. All fields are sent on every Update (the API is idempotent). The retry block uses `tccommon.WriteRetryTimeout` and `tccommon.RetryError`.

## Risks / Trade-offs

- [Risk] Singleton token ID cannot be used to re-import a specific remote config → Mitigation: Importer uses `schema.ImportStatePassthrough`, matching existing CONFIG resources; import reads the global singleton.
- [Risk] Two resources (`tencentcloud_config_deliver_config` and `tencentcloud_config_update_config_deliver`) manage the same backend, which could conflict if both are applied → Mitigation: This is a user-side concern; the resources are independent TF resources and users should choose one. Documented behavior is identical.
- [Trade-off] No-op Delete means `terraform destroy` leaves the remote config in place → Mitigation: Documented; matches the cloud API reality (no delete endpoint).

## Reference: Cloud API struct definitions (from vendor)

### DescribeConfigDeliver (Read)
File: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802/models.go`

Request (line 2288) — no params:
```go
type DescribeConfigDeliverRequest struct {
	*tchttp.BaseRequest
}
```

ResponseParams (line 2313):
```go
type DescribeConfigDeliverResponseParams struct {
	// 投递名称
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`
	// 资源六段式
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`
	// 投递状态 DeliverStatus：0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`
	// 创建时间
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`
	// 日志前缀
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`
	// 投递类型
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`
	// 1：配置变更   2： 资源列表 3：全部
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *DescribeConfigDeliverResponseParams `json:"Response"`
}
```

Client method (client.go:1350): `func (c *Client) DescribeConfigDeliver(request *DescribeConfigDeliverRequest) (response *DescribeConfigDeliverResponse, err error)`.
Service-layer wrapper exists: `ConfigService.DescribeConfigDeliver(ctx)` at `service_tencentcloud_config.go:579`, returns `*DescribeConfigDeliverResponseParams`; uses `ReadRetryTimeout` retry and checks `result == nil || result.Response == nil` → `NonRetryableError`.

### UpdateConfigDeliver (Create/Update)
File: `vendor/.../config/v20220802/models.go`

Request (line 5790):
```go
type UpdateConfigDeliverRequest struct {
	*tchttp.BaseRequest
	// 0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`
	// 投递服务名称
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`
	// 资源六段式
	// COS：qcs::cos:$region:$account:prefix/$appid/$BucketName
	// CLS: qcs::cls:$region:$account:cls/topicId
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`
	// clonfig_fix
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`
	// 投递类型
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`
	// 1：配置变更 2： 资源列表 3：全选
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`
}
```

ResponseParams (line 5840) — only RequestId:
```go
type UpdateConfigDeliverResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

Client method (client.go:3489): `func (c *Client) UpdateConfigDeliver(request *UpdateConfigDeliverRequest) (response *UpdateConfigDeliverResponse, err error)`. Called via `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateConfigDeliverWithContext(ctx, request)`.

**Field type mapping**: `Status` and `DeliverContentType` are `*uint64` (use `helper.IntUint64`); the rest are `*string` (use `helper.String`).

## Reference: Code style (tencentcloud_igtm_strategy + config singleton pattern)

Key structures the implementation must mirror:

1. **Resource definition** — `func ResourceTencentCloudConfigUpdateConfigDeliver() *schema.Resource` returning Create/Read/Update/Delete + Importer (passthrough) + Schema map.

2. **Create handler** (`resource_tc_config_deliver_config.go:71` pattern):
```go
func resourceTencentCloudConfigUpdateConfigDeliverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_update_config_deliver.create")()
	defer tccommon.InconsistentCheck(d, meta)()
	d.SetId(helper.BuildToken())
	return resourceTencentCloudConfigUpdateConfigDeliverUpdate(d, meta)
}
```

3. **Read handler** — use `ConfigService` + `DescribeConfigDeliver`, nil-check each field before `d.Set`. On nil response: `log.Printf("[WARN]%s resource tencentcloud_config_update_config_deliver [%s] not found...", logId, d.Id())` then `d.SetId("")`.

4. **Update handler** — build `UpdateConfigDeliverRequest`, use `GetOkExists` for int fields (`status`, `deliver_content_type`) and `GetOk` for strings, wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)`, then call Read. Log via `log.Printf("[CRITAL]%s update config update_config_deliver failed, reason:%+v", logId, reqErr)`.

5. **Delete handler** — return nil.

6. **Unit tests** — for a NEW resource, use gomock (gomonkey) per project rules, NOT the TF acceptance test suite. Reference mock patterns from other config resources' test files.

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/resource_tc_config_update_config_deliver.go` | New — resource |
| `tencentcloud/services/config/resource_tc_config_update_config_deliver_test.go` | New — unit tests (mock) |
| `tencentcloud/services/config/resource_tc_config_update_config_deliver.md` | New — doc |
| `tencentcloud/services/config/service_tencentcloud_config.go` | No change — reuse existing `DescribeConfigDeliver` |
| `tencentcloud/provider.go` | Modified — register `tencentcloud_config_update_config_deliver` + comment index |