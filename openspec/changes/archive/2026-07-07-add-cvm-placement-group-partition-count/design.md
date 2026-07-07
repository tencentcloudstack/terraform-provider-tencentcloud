## Context

The `tencentcloud_placement_group` resource currently supports creating placement groups with `name`, `type`, and `affinity` parameters. The CVM `CreateDisasterRecoverGroup` API also accepts `PartitionCount` for partition-type placement groups (when `Strategy` is `PARTITION`), but the Terraform resource does not expose this parameter.

**Current state:**
- Resource file: `tencentcloud/services/cvm/resource_tc_placement_group.go`
- Service layer: `tencentcloud/services/cvm/service_tencentcloud_cvm.go` (`CreatePlacementGroup` function)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312`

**API behavior analysis (from vendor SDK):**

| API | PartitionCount in Request | PartitionCount in Response |
|-----|--------------------------|---------------------------|
| `CreateDisasterRecoverGroup` | Yes (`PartitionCount *int64`, range 2-30) | Yes (`PartitionCount *int64`) |
| `DescribeDisasterRecoverGroups` | N/A (no request field) | No (`DisasterRecoverGroup` struct lacks `PartitionCount`) |
| `ModifyDisasterRecoverGroupAttribute` | No (only `Name`, `Affinity`) | N/A |
| `DeleteDisasterRecoverGroups` | No (only `DisasterRecoverGroupIds`) | N/A |

**Key constraint:** `PartitionCount` is only available in the Create request/response. The Describe API's `DisasterRecoverGroup` struct does not include `PartitionCount`, so the value cannot be refreshed on Read.

## Goals / Non-Goals

**Goals:**
- Add `partition_count` (Optional, ForceNew, Computed) parameter to `tencentcloud_placement_group`
- Pass `PartitionCount` to `CreateDisasterRecoverGroup` API when specified by user
- Capture the `PartitionCount` from the Create API response and persist it in Terraform state
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Adding `Strategy` parameter (out of scope for this change)
- Making `PartitionCount` updatable (API does not support it; must be ForceNew)
- Modifying `DescribeDisasterRecoverGroups` SDK model to add `PartitionCount` (SDK limitation)
- Adding `partition_count` to the datasource `tencentcloud_placement_groups` (out of scope)

## Decisions

### Decision 1: `partition_count` is ForceNew

**Rationale:** The `ModifyDisasterRecoverGroupAttribute` API does not accept `PartitionCount`. There is no way to update the partition count of an existing placement group. Therefore, any change to `partition_count` must trigger resource recreation.

### Decision 2: Skip setting `partition_count` on Read

**Rationale:** The `DisasterRecoverGroup` struct returned by `DescribeDisasterRecoverGroups` does not include `PartitionCount`. We cannot refresh the value. Since `partition_count` is ForceNew, Terraform Plugin SDK v2 will retain the Create-time value in state across subsequent reads. This means:
- On initial create: value is set from user config + API response
- On subsequent reads: value is NOT reset (remains from Create), avoiding spurious diffs
- On import: `partition_count` will not be available (documented limitation)

### Decision 3: Use `int` type with validation range 2-30

**Rationale:** Consistent with existing `affinity` field (also `TypeInt` with `ValidateIntegerInRange`). The API documentation states the valid range is 2-30.

### Decision 4: Update `CreatePlacementGroup` service function signature

**Rationale:** The service function `CreatePlacementGroup` currently accepts `(ctx, name, type, affinity, tags)`. We add `partitionCount int` as a new parameter. This is a minor signature change to an internal function, not a breaking change to the public Terraform API.

## Risks / Trade-offs

- **[Risk] Imported resources lose `partition_count`**: When importing an existing placement group, `partition_count` cannot be read from the Describe API. The field will not appear in state.
  - **Mitigation:** Document this limitation in the resource documentation. Users can manually add the value to their config after import.

- **[Risk] No SDK model for `PartitionCount` in Describe response**: The `DisasterRecoverGroup` struct lacks `PartitionCount`.
  - **Mitigation:** This is acceptable since `PartitionCount` is ForceNew. The value is set once at creation and retained in Terraform state.