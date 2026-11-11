## Context

The `tencentcloud_cdwdoris_instance` resource currently creates CDW Doris instances without the `IsSSC` (storage-compute separation) parameter. The SDK's `CreateInstanceNewRequest` struct includes `IsSSC *bool` (line 1075), but the Terraform resource does not expose it. The `DescribeInstance` API returns `InstanceInfo` which does NOT include `IsSSC`, and `ModifyInstance` API only supports `InstanceId` and `InstanceName` — making this parameter write-only and immutable after creation.

## Goals / Non-Goals

**Goals:**
- Add `is_ssc` (TypeBool, Optional, ForceNew) to the resource schema
- Pass the value to `CreateInstanceNewRequest.IsSSC` during resource creation
- Mark the parameter as immutable (ForceNew + added to immutableArgs)

**Non-Goals:**
- Reading `is_ssc` from the API (the `DescribeInstance` response does not include this field)
- Modifying `is_ssc` after creation (the `ModifyInstance` API does not support it)
- Adding related SSC parameters (`SSCCU`, `CacheDiskSize`, `CacheDataDiskSize`) — these are out of scope

## Decisions

1. **Parameter type: `TypeBool`** — matches the SDK's `*bool` type for `IsSSC`.
2. **ForceNew: true** — since `ModifyInstance` does not support this field, changing it would require destroy-and-recreate.
3. **Write-only semantics** — the parameter is set during Create but not read back in Read. The Description will note this is a write-only field. In Read, since `InstanceInfo` has no `IsSSC`, no `d.Set("is_ssc", ...)` call is needed.
4. **Add to immutableArgs** — added to the existing `immutableArgs` list in the Update function to prevent modification attempts with a clear error message.
5. **Use `d.GetOkExists`** — consistent with the existing boolean parameter pattern in the Create function (e.g., `ha_flag`, `enable_multi_zones`).

## Risks / Trade-offs

- **[Risk] Terraform state drift**: Since `is_ssc` is write-only, Terraform cannot detect if the value was changed externally. → **Mitigation**: The parameter is `ForceNew` and the API does not support modifying it, so external changes are not possible.
- **[Risk] Import scenario**: When importing an existing instance, the `is_ssc` value will default to `false`/unset since it cannot be read from the API. → **Mitigation**: Document this limitation. Users importing existing storage-compute separation instances may need to manually set the value in their config.