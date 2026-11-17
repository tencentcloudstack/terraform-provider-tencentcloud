# dbdc-db-custom-node-disaster-recover-group Specification

## Purpose
TBD - created by syncing change add-dbdc-db-custom-node-disaster-recover-group. Update Purpose after sync.
## Requirements
### Requirement: Disaster recover group ids on node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `disaster_recover_group_ids` parameter (TypeList of TypeString, MaxItems 1, ForceNew) that is passed to the `CreateDBCustomNodes` API as `request.DisasterRecoverGroupIds`. The API supports specifying only one placement group ID. Changes to `disaster_recover_group_ids` after creation SHALL force replacement of the resource.

#### Scenario: Create node with a disaster recover group
- **WHEN** a user specifies `disaster_recover_group_ids = ["dbrg-xxxxxxxx"]` in the `tencentcloud_dbdc_db_custom_node` resource configuration
- **THEN** the provider SHALL pass `DisasterRecoverGroupIds=["dbrg-xxxxxxxx"]` in the `CreateDBCustomNodes` API request

#### Scenario: Create node without a disaster recover group
- **WHEN** a user does NOT specify `disaster_recover_group_ids` in the `tencentcloud_dbdc_db_custom_node` resource configuration
- **THEN** the provider SHALL NOT set `DisasterRecoverGroupIds` in the `CreateDBCustomNodes` API request

#### Scenario: More than one disaster recover group id is rejected
- **WHEN** a user specifies more than one element in `disaster_recover_group_ids`
- **THEN** the provider SHALL reject the configuration via `MaxItems: 1`

#### Scenario: Changing disaster recover group ids forces replacement
- **WHEN** a user changes `disaster_recover_group_ids` after creation
- **THEN** the provider SHALL force replacement (recreate) of the `tencentcloud_dbdc_db_custom_node` resource, because no update API can modify the placement group

### Requirement: Disaster recover group id computed output on Read
The `tencentcloud_dbdc_db_custom_node` resource SHALL expose a computed `disaster_recover_group_id` parameter (TypeString, Computed) that is refreshed from the `DescribeDBCustomNodes` API response (`response.NodeSet[].DisasterRecoverGroupId`).

#### Scenario: Read node bound to a disaster recover group
- **WHEN** the provider reads an existing `tencentcloud_dbdc_db_custom_node` whose node is bound to a placement group
- **THEN** `disaster_recover_group_id` SHALL be refreshed from `DBCustomNode.DisasterRecoverGroupId` returned by `DescribeDBCustomNodes`

#### Scenario: Read node not bound to a disaster recover group
- **WHEN** the provider reads an existing `tencentcloud_dbdc_db_custom_node` whose node is NOT bound to a placement group
- **THEN** the provider SHALL skip setting `disaster_recover_group_id` when `DBCustomNode.DisasterRecoverGroupId` is nil, consistent with the nil-guarded Read pattern for other computed fields
