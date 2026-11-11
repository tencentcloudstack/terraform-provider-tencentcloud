## ADDED Requirements

### Requirement: User can add nodes to a MongoDB sharding instance
The `tencentcloud_mongodb_sharding_instance` resource SHALL expose an `add_node_list` parameter that allows users to add nodes (SECONDARY, READONLY, or MONGOS) to an existing sharding instance via the `ModifyDBInstanceSpec` API.

#### Scenario: Add a read-only node
- **WHEN** user specifies `add_node_list` with `role = "READONLY"` and `zone = "ap-guangzhou-2"`
- **THEN** the provider calls `ModifyDBInstanceSpec` with `AddNodeList` containing the specified role and zone, and waits for the operation to complete

#### Scenario: Add multiple node types
- **WHEN** user specifies `add_node_list` with multiple entries (e.g., one SECONDARY and one READONLY)
- **THEN** the provider calls `ModifyDBInstanceSpec` with `AddNodeList` containing all specified nodes

### Requirement: User can remove nodes from a MongoDB sharding instance
The `tencentcloud_mongodb_sharding_instance` resource SHALL expose a `remove_node_list` parameter that allows users to remove nodes from an existing sharding instance via the `ModifyDBInstanceSpec` API.

#### Scenario: Remove a read-only node by node name
- **WHEN** user specifies `remove_node_list` with `role = "READONLY"`, `node_name = "cmgo-xxxx_0-node-readonly0"`, and `zone = "ap-guangzhou-2"`
- **THEN** the provider calls `ModifyDBInstanceSpec` with `RemoveNodeList` containing the specified role, node name, and zone, and waits for the operation to complete

### Requirement: Backward compatibility
The `add_node_list` and `remove_node_list` parameters SHALL both be Optional. Existing Terraform configurations that do not use these parameters SHALL continue to work without any changes.

#### Scenario: Existing configuration without new parameters
- **WHEN** user applies an existing `tencentcloud_mongodb_sharding_instance` configuration that does not include `add_node_list` or `remove_node_list`
- **THEN** the provider behaves identically to before this change, with no diff or modification