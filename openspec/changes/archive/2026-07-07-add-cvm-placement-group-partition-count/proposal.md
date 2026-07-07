## Why

The CVM `CreateDisasterRecoverGroup` API supports `PartitionCount` for partition-type placement groups (`Strategy=PARTITION`), but the Terraform resource `tencentcloud_placement_group` does not expose this parameter. Users who want to create partition placement groups cannot specify the partition count through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `partition_count` (Optional, ForceNew, Computed) parameter to `tencentcloud_placement_group` resource to support specifying the number of partitions when creating a partition placement group (valid range: 2-30).

## Capabilities

### New Capabilities
- `cvm-placement-group-partition-count`: Enable the `partition_count` parameter on the `tencentcloud_placement_group` resource to allow users to specify partition count when creating partition-type placement groups.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cvm/resource_tc_placement_group.go` — add `partition_count` schema field and wire it through Create flow
  - `tencentcloud/services/cvm/service_tencentcloud_cvm.go` — update `CreatePlacementGroup` to accept and pass `PartitionCount`
  - `tencentcloud/services/cvm/resource_tc_placement_group.md` — update documentation
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312` — `CreateDisasterRecoverGroupRequest.ParameterCount` already exists in vendor; no SDK update needed
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and defaults to not being set
- **API constraints:** `PartitionCount` is only accepted by `CreateDisasterRecoverGroup` (not by `ModifyDisasterRecoverGroupAttribute`), so the parameter must be `ForceNew`. The `DescribeDisasterRecoverGroups` response (`DisasterRecoverGroup` struct) does not currently include `PartitionCount`, so the value will be set from the Create API response and persisted in state without being refreshed on Read.