## ADDED Requirements

### Requirement: Data source tencentcloud_cdwdoris_instance_nodes provides node list query

The system SHALL provide a Terraform data source `tencentcloud_cdwdoris_instance_nodes` (RESOURCE_KIND_DATASOURCE) that queries CDWDoris cluster node information by calling the cloud API `DescribeInstanceNodes`. The data source SHALL be read-only with a single Read operation and SHALL NOT create, update, or delete any cloud resources.

#### Scenario: Query nodes by required instance_id

- **WHEN** the user declares `data "tencentcloud_cdwdoris_instance_nodes"` with `instance_id` set to a valid CDWDoris cluster instance ID
- **THEN** the provider SHALL call `DescribeInstanceNodes` with `request.InstanceId` and populate `instance_nodes_list` with each node's `ip`, `spec`, `core`, `memory`, `disk_type`, `disk_size`, `role`, `status`, `rip`, `fe_role`, `uuid`, `zone`, `virtual_zone`, `create_time`, and `compute_group_id`

#### Scenario: Query nodes with optional node_role filter

- **WHEN** the user sets `node_role` (e.g. "data") in addition to `instance_id`
- **THEN** the provider SHALL pass `request.NodeRole` to `DescribeInstanceNodes` so only nodes of that role are returned

#### Scenario: Query nodes with optional display_policy filter

- **WHEN** the user sets `display_policy` (e.g. "All") in addition to `instance_id`
- **THEN** the provider SHALL pass `request.DisplayPolicy` to `DescribeInstanceNodes` to control the display policy of the returned nodes

#### Scenario: instance_id is required

- **WHEN** the user declares the data source without `instance_id`
- **THEN** the provider SHALL reject the configuration because `instance_id` is a Required argument

### Requirement: Output instance_nodes_list element fields map to InstanceNode struct

The system SHALL map `response.Response.InstanceNodesList` to the `instance_nodes_list` Computed list. Each list element SHALL expose the following fields from the `InstanceNode` struct, setting each only when the corresponding response field is non-nil: `ip`, `spec`, `core`, `memory`, `disk_type`, `disk_size`, `role`, `status`, `rip`, `fe_role`, `uuid`, `zone`, `virtual_zone`, `create_time`, `compute_group_id`.

#### Scenario: All node fields populated

- **WHEN** `DescribeInstanceNodes` returns an `InstanceNode` with all fields present
- **THEN** the corresponding `instance_nodes_list` element SHALL contain all 15 field values

#### Scenario: Nullable fields omitted gracefully

- **WHEN** a returned `InstanceNode` has some nil fields
- **THEN** the provider SHALL skip setting those fields for that element without error

### Requirement: Output node_roles maps to response NodeRoles

The system SHALL map `response.Response.NodeRoles` to the `node_roles` Computed list of strings, representing the cluster-supported node role types.

#### Scenario: node_roles populated

- **WHEN** `DescribeInstanceNodes` returns a non-empty `NodeRoles` array
- **THEN** the provider SHALL set `node_roles` to the list of role strings

### Requirement: Service layer paginates DescribeInstanceNodes automatically

The service layer method `DescribeCdwdorisInstanceNodesByFilter` SHALL internally paginate `DescribeInstanceNodes` using `Offset`/`Limit` (page size 10) and aggregate all nodes, without exposing pagination parameters to the Terraform schema.

#### Scenario: Many nodes fetched across pages

- **WHEN** the cluster has more nodes than a single page (Limit=10)
- **THEN** the service layer SHALL issue subsequent requests with incremented `Offset` until all nodes are collected

### Requirement: Read handles empty response without clearing state

The data source Read function SHALL wrap the API call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with errors wrapped via `tccommon.RetryError`. When the API returns an empty response (`response == nil`, `response.Response == nil`, or `len(InstanceNodesList) == 0`), the Read function SHALL return a `NonRetryableError` instead of calling `d.SetId("")`, to avoid clearing state on transient API fluctuations. A `log.Printf("[DATASOURCE] read empty, skip SetId")` notice SHALL be retained on the failure path.

#### Scenario: Empty response returns NonRetryableError

- **WHEN** `DescribeInstanceNodes` returns no nodes (empty list)
- **THEN** the Read function SHALL return a NonRetryableError and SHALL NOT clear the data source id

### Requirement: Provider registration and documentation

The provider SHALL register `tencentcloud_cdwdoris_instance_nodes` in `tencentcloud/provider.go` dataSources map, and a documentation file `data_source_tc_cdwdoris_instance_nodes.md` SHALL be created with a one-line description mentioning CDWDoris product and Example Usage. The provider.md entry SHALL be generated via `make doc` during finalize.

#### Scenario: Data source registered

- **WHEN** the provider is built after the change
- **THEN** `tencentcloud_cdwdoris_instance_nodes` SHALL be available as a valid data source in Terraform