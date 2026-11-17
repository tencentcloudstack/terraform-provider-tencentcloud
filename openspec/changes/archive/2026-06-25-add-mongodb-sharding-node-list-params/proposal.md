## Why

The `tencentcloud_mongodb_sharding_instance` resource currently uses `ModifyDBInstanceSpec` API for scaling memory/volume, but does not expose `AddNodeList` and `RemoveNodeList` parameters that allow users to dynamically add or remove nodes (Mongod nodes, read-only nodes, Mongos nodes) from a sharding instance. This prevents Terraform users from managing node topology changes declaratively.

## What Changes

- Add `add_node_list` parameter (TypeList) to `tencentcloud_mongodb_sharding_instance` resource, supporting the following sub-fields:
  - `role` (TypeString, Required): Node role to add - `SECONDARY`, `READONLY`, or `MONGOS`
  - `zone` (TypeString, Required): Availability zone for the new node
- Add `remove_node_list` parameter (TypeList) to `tencentcloud_mongodb_sharding_instance` resource, supporting the following sub-fields:
  - `role` (TypeString, Required): Node role to remove - `SECONDARY`, `READONLY`, or `MONGOS`
  - `node_name` (TypeString, Required): Node ID to remove (e.g., `cmgo-xxxx_0-node-readonly0`)
- Wire these parameters through the existing `ModifyDBInstanceSpec` API call in the update path

## Capabilities

### New Capabilities
- `mongodb-sharding-node-list-management`: Ability to add and remove nodes (Mongod, read-only, Mongos) from a MongoDB sharding instance via `add_node_list` and `remove_node_list` parameters in the `ModifyDBInstanceSpec` API.

### Modified Capabilities
<!-- None - this is purely additive, no existing capability requirements change -->

## Impact

- **Code**: `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go` - Add schema definitions and wire parameters in the update function
- **Code**: `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance_test.go` - Add test cases for new parameters
- **Docs**: `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.md` - Add example and parameter documentation
- **API**: Uses existing `ModifyDBInstanceSpec` API (already available in vendor SDK `v20190725`)
- **Backward Compatibility**: Fully backward compatible - both new fields are Optional