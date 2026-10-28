## Why

The CVM `CreateDisasterRecoverGroup` API supports `Strategy` and `PartitionCount` for partition-type placement groups (`Strategy=PARTITION`), but the Terraform resource `tencentcloud_placement_group` does not expose these parameters. Users who want to create partition placement groups cannot specify the strategy or partition count through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `strategy` (Optional, Computed, immutable after creation) parameter to `tencentcloud_placement_group` resource to support specifying the placement group strategy (valid values: `SPREAD`, `PARTITION`). Default is `SPREAD`. Changes to `strategy` after creation are rejected via immutable args check in Update.
- Add `partition_count` (Optional, Computed, immutable after creation) parameter to `tencentcloud_placement_group` resource to support specifying the number of partitions when creating a partition placement group (valid range: 2-30). Only valid when `strategy` is `PARTITION`. Changes to `partition_count` after creation are rejected via immutable args check in Update.
- Add validation in Create: if `partition_count` is set but `strategy` is not `PARTITION`, return an error.
- Use `immutableArgs` array pattern in Update function (consistent with ckafka resource pattern) for `type`, `strategy`, `affinity`, and `partition_count`.
- Read function reads `Strategy` and `PartitionCount` from `DescribeDisasterRecoverGroups` API response (requires updated SDK), so imported resources can have these fields populated.
- `CreatePlacementGroup` service function returns the full `*cvm.CreateDisasterRecoverGroupResponse` for callers to extract needed fields (id, etc.), rather than returning individual values.

## Capabilities

### New Capabilities
- `cvm-placement-group-partition-count`: Enable the `strategy` and `partition_count` parameters on the `tencentcloud_placement_group` resource to allow users to specify strategy and partition count when creating placement groups.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cvm/resource_tc_placement_group.go` — add `strategy` and `partition_count` schema fields, add validation, wire through Create flow, add immutable args to Update, add Read support
  - `tencentcloud/services/cvm/service_tencentcloud_cvm.go` — update `CreatePlacementGroup` to accept and pass `Strategy` and `PartitionCount`, return full response
  - `tencentcloud/services/cvm/extension_cvm.go` — add `CVM_PLACEMENT_GROUP_STRATEGY` constants
  - `tencentcloud/services/cvm/resource_tc_placement_group.md` — update documentation
- **SDK dependency:** Updated `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312` — `DisasterRecoverGroup` struct now includes `Strategy` and `PartitionCount` fields, enabling Read to refresh these values
- **Backward compatibility:** fully backward compatible — the new parameters are Optional and default to not being set
- **API constraints:** `PartitionCount` and `Strategy` are only accepted by `CreateDisasterRecoverGroup` (not by `ModifyDisasterRecoverGroupAttribute`), so these parameters are immutable after creation. The `DescribeDisasterRecoverGroups` response now includes `Strategy` and `PartitionCount` (after SDK update), so Read can refresh these values.