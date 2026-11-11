## Why

The `dbdc` `CreateDBCustomNodes` API supports specifying a disaster-recover (placement) group via the `DisasterRecoverGroupIds` request parameter, and the `DescribeDBCustomNodes` API returns the bound `DisasterRecoverGroupId` on each node, but the Terraform resource `tencentcloud_dbdc_db_custom_node` does not expose this parameter. Users who want to place a DB Custom node into a placement group at creation time cannot do so through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `disaster_recover_group_ids` (Optional, `ForceNew`, `TypeList` of `TypeString`) parameter to the `tencentcloud_dbdc_db_custom_node` resource. It maps to `request.DisasterRecoverGroupIds` in the `CreateDBCustomNodes` API. The API documents that only one placement group ID is supported ("仅支持指定一个"), so `MaxItems: 1` is enforced.
- Add `disaster_recover_group_id` (Computed, `TypeString`) parameter to the `tencentcloud_dbdc_db_custom_node` resource. It maps to `response.NodeSet[].DisasterRecoverGroupId` in the `DescribeDBCustomNodes` API and is refreshed during Read.

## Capabilities

### New Capabilities
- `dbdc-db-custom-node-disaster-recover-group`: Enable the `disaster_recover_group_ids` input parameter and the `disaster_recover_group_id` computed parameter on the `tencentcloud_dbdc_db_custom_node` resource, so a node can be placed into a DB Custom disaster-recover (placement) group at creation time and the bound group can be refreshed on Read.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go` — add `disaster_recover_group_ids` and `disaster_recover_group_id` schema fields, wire `disaster_recover_group_ids` into the Create flow, add Read support for `disaster_recover_group_id`
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node_test.go` — add unit test coverage for the new parameters
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.md` — update documentation example
- **SDK dependency:** The vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` already includes `DisasterRecoverGroupIds` in `CreateDBCustomNodesRequest` and `DisasterRecoverGroupId` in the `DBCustomNode` response struct. No SDK upgrade is required.
- **Backward compatibility:** fully backward compatible — the new input parameter is Optional and defaults to not being set; the new computed parameter is read-only.
- **API constraints:** `DisasterRecoverGroupIds` is only accepted by `CreateDBCustomNodes` (not by any update API), so `disaster_recover_group_ids` is `ForceNew`. The `DescribeDBCustomNodes` response includes `DisasterRecoverGroupId`, so Read can refresh `disaster_recover_group_id`.
