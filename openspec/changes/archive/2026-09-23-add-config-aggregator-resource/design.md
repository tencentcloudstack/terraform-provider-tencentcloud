## Context

This change adds a new Terraform resource `tencentcloud_config_aggregator` to manage Tencent Cloud Config (配置审计) "账号组" (Aggregator). The Aggregator groups member accounts for unified compliance evaluation. The config SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802` is already vendored and exposes the four CRUD APIs used here. All four APIs are synchronous (no flow/task id to poll).

Reference style file: `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` (composite-id GENERAL resource).

## Goals / Non-Goals

**Goals:**
- Full CRUD for the Aggregator via Terraform.
- Composite resource ID `account_group_id#owner_uin`.
- gomonkey-based unit tests (no terraform test suite, no real API calls).

**Non-Goals:**
- Not managing the member accounts' resources themselves (only the aggregator membership list).
- No async polling (APIs are synchronous).

## Cloud API Reference (from vendor)

Package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`
File: `vendor/.../config/v20220802/models.go`

### AggregatorAccount struct (models.go:684)
```go
type AggregatorAccount struct {
    MemberUin   *uint64 `json:"MemberUin,omitnil,omitempty"`   // 成员ID
    MemberName  *string `json:"MemberName,omitnil,omitempty"`  // 成员名称 (may be null)
}
```

### CreateAggregator
- Request (models.go:1241) `CreateAggregatorRequest`:
  - `Name *string` — 账号组名称 (required)
  - `Description *string` — 账号组描述 (required)
  - `Type *string` — 账号组类型, 枚举 `RD`(全局账号组) / `CUSTOM`(自定义账号组) (required)
  - `AggregatorAccounts []*AggregatorAccount` — 成员信息列表, 最多100个 (optional)
- Response (models.go:1280) `CreateAggregatorResponseParams`:
  - `AccountGroupId *string` — 账号组Id
- Client method: `(*Client).CreateAggregatorWithContext(ctx, request)` (client.go:529)

### DescribeAggregator
- Request (models.go:2138) `DescribeAggregatorRequest`:
  - `AccountGroupId *string` — 账号组ID (required)
  - `OwnerUin *uint64` — 账号组创建者ID (required)
- Response (models.go:2169) `DescribeAggregatorResponseParams`:
  - `Name *string`
  - `Description *string`
  - `Type *string`
  - `AggregatorAccounts []*AggregatorAccount` (may be null)
  - `AggregatorStatus *uint64` — 创建状态
- Client method: `(*Client).DescribeAggregatorWithContext(ctx, request)` (client.go:1257)

### UpdateAggregator
- Request (models.go:5449) `UpdateAggregatorRequest`:
  - `Name *string` (required)
  - `Description *string` (required)
  - `AccountGroupId *string` (required)
  - `OwnerUin *uint64` (required)
  - `AggregatorAccounts []*AggregatorAccount` (optional)
- Response: only `RequestId`.
- Client method: `(*Client).UpdateAggregatorWithContext(ctx, request)` (client.go:3281)

### DeleteAggregators
- Request (models.go:1520) `DeleteAggregatorsRequest`:
  - `AccountGroupId *string` (required)
  - `OwnerUin *uint64` (required)
- Response: only `RequestId`.
- Client method: `(*Client).DeleteAggregatorsWithContext(ctx, request)` (client.go:763)

### Notes on OwnerUin
`OwnerUin` is `*uint64`. Create does NOT return it, but Describe/Update/Delete all require it as an input. Therefore `owner_uin` MUST be a user-provided schema field (TypeString in TF; converted via `helper.StrToUint64`). It is stored in the composite ID and re-read from the ID in Read/Update/Delete (DescribeAggregator does not return OwnerUin).

## Schema

### Required
| Field | TF Type | ForceNew | Description |
|---|---|---|---|
| `name` | String | No | 账号组名称 |
| `description` | String | No | 账号组描述 |
| `type` | String | Yes | 账号组类型: `RD` / `CUSTOM` |
| `owner_uin` | String | Yes | 账号组创建者 UIN (needed by Describe/Update/Delete; Create does not return it) |

### Optional
| Field | TF Type | Description |
|---|---|---|
| `aggregator_accounts` | List(object) | 成员信息列表, 最多100个 |

#### aggregator_accounts element schema
| Field | TF Type | Required | Description |
|---|---|---|---|
| `member_uin` | Int | Yes | 成员ID (uint64) |
| `member_name` | String | Yes | 成员名称 |

### Computed
| Field | TF Type | Description |
|---|---|---|
| `account_group_id` | String | 账号组ID (cloud-assigned; part of composite ID) |
| `aggregator_status` | Int | 创建状态 (uint64) |

`type` is ForceNew because UpdateAggregator accepts it as a path param but it represents the aggregator kind (the create API documents it as the aggregator type enum); changing type on an existing aggregator is not meaningful. `owner_uin` is ForceNew because it identifies the creator account context.

## Resource ID

Composite ID = `strings.Join([]string{account_group_id, owner_uin}, tccommon.FILED_SP)`.

Import support: Yes (RESOURCE_KIND_GENERAL). Import uses composite id `account_group_id#owner_uin`.

## CRUD Design

### Create (`resourceTencentCloudConfigAggregatorCreate`)
1. Build `CreateAggregatorRequest`: `Name`, `Description`, `Type` from schema; `AggregatorAccounts` (iterate list, set `MemberUin` via `helper.IntUint64`, `MemberName`).
2. `resource.Retry(tccommon.WriteRetryTimeout, ...)`: call `CreateAggregatorWithContext`. On error → `tccommon.RetryError(e)`. Check `result==nil || result.Response==nil` → `NonRetryableError`.
3. After retry: check `response.Response.AccountGroupId == nil` / `""` → `NonRetryableError`. Print `logId` + `d.Id()`.
4. `d.SetId(strings.Join([]string{accountGroupId, ownerUin}, tccommon.FILED_SP))`.
5. `return resourceTencentCloudConfigAggregatorRead(d, meta)`.

### Read (`resourceTencentCloudConfigAggregatorRead`)
- Parse id into `accountGroupId`, `ownerUin` (must be 2 parts).
- Optional service helper `DescribeConfigAggregatorById(ctx, accountGroupId, ownerUin)` wrapping `DescribeAggregatorWithContext` with `tccommon.ReadRetryTimeout` retry.
- If response/instance empty: `log.Printf("[CRUD] tencentcloud_config_aggregator id=%s", d.Id())` then `d.SetId("")`.
- Set `name`, `description`, `type`, `aggregator_status` (nil-check before each set). Set `account_group_id`.
- Set `aggregator_accounts` list from `AggregatorAccounts` (map `MemberUin`→`member_uin`, `MemberName`→`member_name`).

### Update (`resourceTencentCloudConfigAggregatorUpdate`)
- Parse id. `accountGroupId`, `ownerUin` from id.
- `needChange` over mutable args `["name", "description", "aggregator_accounts"]`.
- If changed: build `UpdateAggregatorRequest` with `Name`, `Description`, `AccountGroupId`, `OwnerUin` (from id), `AggregatorAccounts`. `resource.Retry(tccommon.WriteRetryTimeout, ...)` → `UpdateAggregatorWithContext`. Check NilResponse → `NonRetryableError`.
- Note: `type` and `owner_uin` are ForceNew so not in mutable args.
- `return resourceTencentCloudConfigAggregatorRead(d, meta)`.

### Delete (`resourceTencentCloudConfigAggregatorDelete`)
- Parse id. Build `DeleteAggregatorsRequest`: `AccountGroupId`, `OwnerUin`.
- `resource.Retry(tccommon.WriteRetryTimeout, ...)` → `DeleteAggregatorsWithContext`. Check NilResponse → `NonRetryableError`.

## Code Style Reference (tencentcloud_igtm_strategy)

Key structural points to follow (file: `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`):
- `ResourceTencentCloudIgtmStrategy()` returns `*schema.Resource` with `Create/Read/Update/Delete` + `Importer{State: schema.ImportStatePassthrough}`.
- Each CRUD fn: `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)`; `logId := tccommon.GetLogId(tccommon.ContextNil)`; `ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`.
- Composite id parsed via `strings.Split(d.Id(), tccommon.FILED_SP)` with length check (==2) → error if broken.
- Retry blocks use `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError{...})` for write ops; nil-response guard returning `resource.NonRetryableError`.
- Update uses `mutableArgs` slice + `d.HasChange` to decide `needChange`.
- Client access: `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client()`.
- Nested list build/flatten pattern: iterate `[]interface{}`, build SDK struct, `append`; flatten returns `[]map[string]interface{}`.

## Test Strategy (gomonkey)

Follow `tencentcloud/services/mongodb/resource_tc_mongodb_instance_test.go` gomonkey pattern:
- `mockMeta` struct implementing `tccommon.ProviderMeta` returning `&connectivity.TencentCloudClient{Region: "ap-guangzhou"}`.
- `patches.ApplyMethodReturn(mockMeta.client, "UseConfigV20220802Client", configClient)`.
- `patches.ApplyMethodFunc(configClient, "CreateAggregatorWithContext", func(...){return mockCreateResponse, nil})` and similarly for Describe/Update/Delete.
- Test Create→Read→Update→Delete business logic (no terraform test suite, no real API).
- Run note: `go test ./tencentcloud/services/config/ -run "TestConfigAggregator" -v -count=1 -gcflags="all=-l"` (do NOT actually run; for documentation only).

## Decisions

1. **owner_uin is user-provided & ForceNew**: Create doesn't return it; Describe/Update/Delete require it. Stored in composite ID. User must know their creator UIN (the account running the config service). ForceNew because it is the identity context.
2. **Composite ID**: `account_group_id#owner_uin` via `tccommon.FILED_SP`, matching igtm_strategy pattern. enables Read/Update/Delete to recover both keys without extra state.
3. **type is ForceNew**: Although UpdateAggregator has an `AccountGroupId` field whose comment mentions type enum, `Type` is not a real update param of UpdateAggregator (only Name/Description/AccountGroupId/OwnerUin/AggregatorAccounts). Changing aggregator type is not supported; treat as recreate.
4. **aggregator_accounts optional**: For `RD` (global) type there may be no explicit member list.
5. **aggregator_status computed only**: uint64 status from Describe; not settable.

## Risks / Trade-offs

- [Risk] OwnerUin must be supplied by user and could drift from actual creator UIN → Mitigation: document clearly in schema description; it is the responsibility of the operator to provide the correct creator UIN.
- [Risk] DescribeAggregator does not return OwnerUin, so the field is never refreshed from cloud → Mitigation: OwnerUin stored in composite ID; Read sets it from the ID, not the (absent) response field. After import the user-imported composite id carries it.
- [Trade-off] Using gomonkey tests only (no ACC tests) per workflow rule → faster, no credentials needed; relies on mock fidelity.