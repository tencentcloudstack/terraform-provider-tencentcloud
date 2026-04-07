## Why

The `tencentcloud_instance` resource currently does not support the `DedicatedResourcePackTenancy` and `DedicatedResourcePackIds` parameters in the `Placement` structure of the `RunInstances` API. These parameters are essential for users who want to create CVM instances using dedicated resource pool packs (资源池包). Without these parameters, users cannot leverage the dedicated resource pool functionality when creating instances through Terraform.

## What Changes

- Add two new optional schema fields to `tencentcloud_instance` resource:
  - `dedicated_resource_pack_tenancy`: String field to specify the dedicated resource pack tenancy strategy (e.g., "ResourcePool")
  - `dedicated_resource_pack_ids`: List of strings to specify the dedicated resource pack IDs to use for instance creation
- Map these fields to the `Placement.DedicatedResourcePackTenancy` and `Placement.DedicatedResourcePackIds` parameters in the `RunInstances` API call
- Add validation to ensure both fields are specified together (if one is set, the other must also be set)
- Mark both fields as `ForceNew: true` since changing these parameters requires recreating the instance

## Capabilities

### New Capabilities

None - this is an enhancement to an existing resource.

### Modified Capabilities

- `instance-resource`: Add support for dedicated resource pack placement parameters to the existing `tencentcloud_instance` resource

## Impact

**Affected Code:**
- `tencentcloud/services/cvm/resource_tc_instance.go`: Schema definition and Create function
- `tencentcloud/services/cvm/resource_tc_instance.md`: Documentation with usage examples

**User Impact:**
- **Non-breaking**: Existing configurations continue to work unchanged
- Users can now specify dedicated resource pack parameters when creating instances
- Supports the resource pool pack feature introduced in `tencentcloud_cvm_resource_pool_packs`

**Dependencies:**
- Requires no changes to the SDK (parameters already exist in `cvm.Placement` struct)
- Works with the existing `RunInstances` API
