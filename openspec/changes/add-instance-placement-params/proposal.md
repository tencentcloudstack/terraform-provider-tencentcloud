## Why

The `tencentcloud_instance` resource currently does not support the `DedicatedResourcePackTenancy` and `DedicatedResourcePackIds` parameters in the `Placement` structure of the `RunInstances` API. These parameters are essential for users who want to create CVM instances using dedicated resource pool packs (资源池包). Without these parameters, users cannot leverage the dedicated resource pool functionality when creating instances through Terraform.

Additionally, the `disaster_recover_group_ids` parameter now supports batch setting of placement groups (up to 3), and the newer `DescribeInstances` API also returns `DisasterRecoverGroupIds` in the response. When `disaster_recover_group_ids` is set, it should take priority over the singular `placement_group_id` for CRUD operations.

## What Changes

### Dedicated Resource Pack Placement Parameters
- Add two new optional schema fields to `tencentcloud_instance` resource:
  - `dedicated_resource_pack_tenancy`: String field to specify the dedicated resource pack tenancy strategy (e.g., "ResourcePool")
  - `dedicated_resource_pack_ids`: List of strings to specify the dedicated resource pack IDs to use for instance creation
- Map these fields to the `Placement.DedicatedResourcePackTenancy` and `Placement.DedicatedResourcePackIds` parameters in the `RunInstances` API call
- Add validation to ensure both fields are specified together (if one is set, the other must also be set)
- Mark both fields as `ForceNew: true` since changing these parameters requires recreating the instance

### Disaster Recover Group IDs Priority
- `disaster_recover_group_ids` (TypeSet, MaxItems=3) now takes priority over `placement_group_id`:
  - **Create**: If `disaster_recover_group_ids` is set, it is used and `placement_group_id` is ignored (for both the initial `RunInstances` request and the `rpgFlag` post-create call)
  - **Read**: When `disaster_recover_group_ids` is in state, `placement_group_id` is NOT read back from API to avoid plan diffs
  - **Read**: `disaster_recover_group_ids` is populated from `instance.DisasterRecoverGroupIds` response field
  - **Update**: If `disaster_recover_group_ids` is set, changes to `placement_group_id` are rejected with an error
- Removed `ConflictsWith` between `disaster_recover_group_ids` and `placement_group_id`
- Added `MaxItems: 3` to `disaster_recover_group_ids` (API supports up to 3 group IDs)
- Added `Computed: true` to `disaster_recover_group_ids`

## Capabilities

### New Capabilities

None - this is an enhancement to an existing resource.

### Modified Capabilities

- `instance-resource`: Add support for dedicated resource pack placement parameters and disaster_recover_group_ids priority logic to the existing `tencentcloud_instance` resource

## Impact

**Affected Code:**
- `tencentcloud/services/cvm/resource_tc_instance.go`: Schema definition, Create function, Read function, Update function
- `tencentcloud/services/cvm/resource_tc_instance.md`: Documentation with usage examples

**User Impact:**
- **Non-breaking**: Existing configurations continue to work unchanged
- Users can now specify dedicated resource pack parameters when creating instances
- Supports the resource pool pack feature introduced in `tencentcloud_cvm_resource_pool_packs`
- `disaster_recover_group_ids` now supports batch setting of up to 3 placement groups
- When `disaster_recover_group_ids` and `placement_group_id` are both configured, `disaster_recover_group_ids` takes priority

**Dependencies:**
- Requires no changes to the SDK (parameters already exist in `cvm.Placement` struct and `cvm.Instance` struct)
- Works with the existing `RunInstances` and `DescribeInstances` APIs
