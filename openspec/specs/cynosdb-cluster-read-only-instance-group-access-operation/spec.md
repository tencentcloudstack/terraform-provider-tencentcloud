## ADDED Requirements

### Requirement: Open read-only instance group access operation resource

The system SHALL provide a Terraform resource `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` that opens read-only instance group access for a CynosDB cluster by calling the `OpenClusterReadOnlyInstanceGroupAccess` API.

The resource SHALL accept the following input parameters:
- `cluster_id` (Required, ForceNew, String): The CynosDB cluster ID.
- `port` (Optional, ForceNew, String): The port for the read-only instance group access.
- `security_group_ids` (Optional, ForceNew, List of String): Security group IDs to associate.

The resource SHALL expose the following computed attributes:
- `flow_id` (Computed, Int): The flow ID returned by the async operation.

The resource SHALL be a one-shot operation (RESOURCE_KIND_OPERATION):
- Create: Calls `OpenClusterReadOnlyInstanceGroupAccess` API, then polls `DescribeFlow` until completion.
- Read: Returns nil (no state to read back).
- Delete: Returns nil (nothing to destroy).

#### Scenario: Successful open read-only instance group access

- **WHEN** user applies a Terraform configuration with `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` specifying a valid `cluster_id`
- **THEN** the resource SHALL call `OpenClusterReadOnlyInstanceGroupAccess` with the provided parameters, poll `DescribeFlow` until the flow completes successfully, store the `flow_id` in state, and set a generated token as the resource ID

#### Scenario: API returns nil response

- **WHEN** the `OpenClusterReadOnlyInstanceGroupAccess` API returns a nil response or nil FlowId
- **THEN** the resource SHALL return a `NonRetryableError` indicating the operation failed

#### Scenario: Async flow polling timeout

- **WHEN** the `DescribeFlow` polling exceeds the configured Create timeout
- **THEN** the resource SHALL return an error indicating the operation timed out

#### Scenario: Async flow reports failure

- **WHEN** `DescribeFlow` returns a failure status (status == 2)
- **THEN** the resource SHALL return an error indicating the flow execution failed

### Requirement: Resource registration in provider

The resource `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider

- **WHEN** a user references `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` in their Terraform configuration
- **THEN** the provider SHALL recognize the resource type and route CRUD operations to the correct handler functions

### Requirement: Resource documentation

The resource SHALL have a documentation file at `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation.md` with:
- A one-line description mentioning CynosDB
- An Example Usage section showing a basic configuration
- No Import section (operation resources are not importable)

#### Scenario: Documentation provides valid example

- **WHEN** a user reads the resource documentation
- **THEN** the example SHALL show a valid Terraform configuration using `cluster_id` as a required parameter
