## Context

The `tencentcloud_cdwdoris_instance` resource (RESOURCE_KIND_GENERAL) manages the full lifecycle of a CDWDoris cluster via the `CreateInstanceNew`, `DescribeInstance`, `ModifyInstance`, and `DestroyInstance` APIs of package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228`.

**Current state:**
- Resource file: `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance.go`
- Service layer: `tencentcloud/services/cdwdoris/service_tencentcloud_cdwdoris.go` (`DescribeCdwdorisInstanceById` returns `*cdwdorisv20211228.InstanceInfo`)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228` (vendored; the struct `CreateInstanceNewRequest` and `NodesSummary` already expose the needed fields `DiskEncrypt` and `Encrypt`).

The resource currently does not expose disk encryption. The `CreateInstanceNew` request now supports a `DiskEncrypt` field, and `DescribeInstance` returns the encryption status through `InstanceInfo.MasterSummary.Encrypt` / `InstanceInfo.CoreSummary.Encrypt` (both `NodesSummary` carry the same `Encrypt` value).

### Cloud API field definitions (from vendor models.go)

**`CreateInstanceNewRequest.DiskEncrypt`** — `models.go:1028` (inside `CreateInstanceNewRequestParams`, mirrored in `CreateInstanceNewRequest` at line 1095):
```go
// 磁盘加密
DiskEncrypt *int64 `json:"DiskEncrypt,omitnil,omitempty" name:"DiskEncrypt"`
```
The request struct `CreateInstanceNewRequest` (line 1031) embeds `*tchttp.BaseRequest` plus the same fields; `DiskEncrypt` is a direct top-level field of the request (not nested inside `FeSpec`/`BeSpec`).

**`DescribeInstanceResponseParams.InstanceInfo`** — `models.go:3095`:
```go
type DescribeInstanceResponseParams struct {
    // 实例描述信息
    InstanceInfo *InstanceInfo `json:"InstanceInfo,omitnil,omitempty" name:"InstanceInfo"`
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

**`InstanceInfo.MasterSummary` / `InstanceInfo.CoreSummary`** — `models.go:4388` (both are `*NodesSummary`):
```go
// 数据节点描述信息
MasterSummary *NodesSummary `json:"MasterSummary,omitnil,omitempty" name:"MasterSummary"`  // line 4427
// zookeeper节点描述信息
CoreSummary *NodesSummary `json:"CoreSummary,omitnil,omitempty" name:"CoreSummary"`            // line 4430
```

**`NodesSummary.Encrypt`** — `models.go:5602`:
```go
// 是否加密
Encrypt *int64 `json:"Encrypt,omitnil,omitempty" name:"Encrypt"`
```

### Existing schema structure relevant to the new parameters

The resource schema (`ResourceTencentCloudCdwdorisInstance` in `resource_tc_cdwdoris_instance.go`) currently has no encryption field. Relevant adjacent fields:
- `fe_spec` (TypeList, MaxItems 1) — read from `respData.MasterSummary` (Spec/NodeSize/Disk).
- `be_spec` (TypeList, MaxItems 1) — read from `respData.CoreSummary` (Spec/NodeSize/Disk).

The Read function already retrieves `respData` (a `*InstanceInfo`) via `service.DescribeCdwdorisInstanceById(ctx, instanceId)`, so `respData.MasterSummary.Encrypt` and `respData.CoreSummary.Encrypt` are directly accessible in the existing Read flow.

The Create function builds a `CreateInstanceNewRequest` and currently sets fields like `Zone`, `FeSpec`, `BeSpec`, `HaFlag`, etc., then calls the API inside a `resource.Retry(tccommon.WriteRetryTimeout, ...)` block. The new `request.DiskEncrypt` assignment must be placed **outside** the retry block (alongside the other request field assignments, before the retry), consistent with existing pattern.

The Update function already contains an `immutableArgs` array (line 667):
```go
immutableArgs := []string{"zone", "fe_spec", "be_spec", "ha_flag", "user_vpc_id", "user_subnet_id", "product_version", "charge_type", "charge_properties", "doris_user_pwd", "tags", "case_sensitive", "enable_multi_zones", "user_multi_zone_infos", "is_ssc"}
```
`disk_encrypt` must be added to this array because there is no update API that accepts `DiskEncrypt`.

### API behavior analysis

| API | DiskEncrypt in Request | Encrypt in Response |
|-----|------------------------|---------------------|
| `CreateInstanceNew` | Yes (`DiskEncrypt *int64` on request) | No |
| `DescribeInstance` | N/A | Yes (`InstanceInfo.MasterSummary.Encrypt` and `InstanceInfo.CoreSummary.Encrypt`, both `*int64`) |
| `ModifyInstance` | No | N/A |
| `DestroyInstance` | No | N/A |

**Key constraint:** `DiskEncrypt` is only available in the Create request. `Encrypt` is read from the Describe response, so it can be refreshed on Read and supports imported resources.

## Goals / Non-Goals

**Goals:**
- Add `disk_encrypt` (Optional, TypeInt, immutable after creation) parameter to `tencentcloud_cdwdoris_instance` that is passed to the `CreateInstanceNew` API as `DiskEncrypt`.
- Add `encrypt` (Computed, TypeInt) parameter to `tencentcloud_cdwdoris_instance` that is refreshed from the `DescribeInstance` API response (`MasterSummary.Encrypt`, falling back to `CoreSummary.Encrypt`).
- Reject changes to `disk_encrypt` after creation via the existing `immutableArgs` array in the Update function.
- Maintain full backward compatibility — existing configurations continue to work unchanged.

**Non-Goals:**
- Making `disk_encrypt` updatable (no API supports changing encryption post-creation).
- Adding `disk_encrypt`/`encrypt` to the `tencentcloud_cdwdoris_instances` datasource (out of scope).

## Decisions

### Decision 1: `disk_encrypt` is immutable via `immutableArgs` (not ForceNew)

**Rationale:** There is no update API that accepts `DiskEncrypt`. Instead of using `ForceNew: true` (which silently destroys and recreates the cluster), we add `disk_encrypt` to the existing `immutableArgs` array in the Update function. This returns a clear error when the value changes, which is safer and more explicit than silently destroying a production cluster. This is consistent with how the resource already handles other creation-only parameters (`zone`, `fe_spec`, `be_spec`, etc.).

### Decision 2: Use TypeInt for `disk_encrypt` and `encrypt`

**Rationale:** Both the request field `DiskEncrypt` and the response field `Encrypt` are `*int64` in the SDK. Using `TypeInt` in the schema keeps a direct type mapping without unnecessary string conversion. The API treats the value as a numeric flag (e.g. `1` = enable encryption).

### Decision 3: `encrypt` is Computed, read from `MasterSummary.Encrypt` / `CoreSummary.Encrypt`

**Rationale:** The `DescribeInstance` response exposes the encryption status on both `MasterSummary` and `CoreSummary` (both are `NodesSummary`, carrying the same `Encrypt` value per the requirement mapping). We set the single `encrypt` computed field from `MasterSummary.Encrypt` when present; if `MasterSummary` or its `Encrypt` is nil, fall back to `CoreSummary.Encrypt`. This follows the resource Read convention of nil-checking before `d.Set`.

### Decision 4: `disk_encrypt` assignment placed before the retry block in Create

**Rationale:** Per project convention, request field assignments must be placed outside the retry block; the retry block should only contain the API call and response nil-checks. `request.DiskEncrypt` is set right after the other request field assignments (e.g. after `is_ssc`), before `err := resource.Retry(...)`.

## Risks / Trade-offs

- **[Risk] Users may import existing clusters and have `encrypt` populated**: `DescribeInstance` returns `Encrypt`, so imported resources will have the correct value. No special handling needed.
- **[Risk] Changing `disk_encrypt` after creation**: Users receive a clear error `argument 'disk_encrypt' cannot be changed` rather than silent destruction.
  - **Mitigation:** `immutableArgs` check returns an explicit error.
- **[Risk] `MasterSummary.Encrypt` and `CoreSummary.Encrypt` may differ in edge cases**: Per the requirement mapping both map to the same `Encrypt` schema field. We prefer `MasterSummary.Encrypt` and fall back to `CoreSummary.Encrypt` to always populate the computed value when the API returns it.