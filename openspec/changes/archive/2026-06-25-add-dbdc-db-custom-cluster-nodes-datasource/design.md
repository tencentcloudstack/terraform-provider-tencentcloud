## Context

The Terraform Provider for TencentCloud already includes a datasource `tencentcloud_dbdc_db_custom_clusters` for querying DB Custom cluster list, implemented in `tencentcloud/services/dbdc/`. The service layer (`service_tencentcloud_dbdc.go`) contains `DbdcService` struct with `DescribeDBCustomClustersByFilter` method that handles paginated API calls.

Users now need to query nodes within a specific DB Custom cluster. The `DescribeDBCustomClusterNodes` API in `dbdc/v20201029` SDK provides this capability, requiring a `ClusterId` (required) and optional `Filters` (supporting `node-name` filter), with pagination via `Offset`/`Limit` (max Limit=100).

The response returns `NodeSet` (array of `DBCustomClusterNode` structs with 7 fields: NodeId, NodeName, LanIP, SSHEndpoint, Status, Zone, NodeType) and `TotalCount`.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_dbdc_db_custom_cluster_nodes` datasource following the existing `tencentcloud_dbdc_db_custom_clusters` pattern
- Support querying nodes by `cluster_id` (Required) and optional `filters`
- Return complete node details in `node_set` computed field
- Implement paginated data retrieval in service layer (internal, not exposed to users)
- Add gomonkey-based unit tests for the datasource
- Register the datasource in provider.go and provider.md

**Non-Goals:**
- Creating a full resource (CRUD) for DB Custom cluster nodes - only a datasource (Read)
- Exposing `offset`/`limit` parameters to Terraform users - handled internally in service layer
- Modifying the existing `tencentcloud_dbdc_db_custom_clusters` datasource behavior

## Decisions

### Decision 1: Schema Design - Follow existing dbdc datasource pattern
**Choice**: Mirror the schema pattern from `data_source_tc_dbdc_db_custom_clusters.go`
**Rationale**: The existing dbdc datasource already establishes the pattern for this product. Key differences:
- `cluster_id` is Required (not Optional) since the API requires it to identify which cluster's nodes to query
- `filters` follows the same structure (name + values) as the existing datasource, but supports only `node-name` filter
- `node_set` computed field contains flattened node attributes (NodeId, NodeName, LanIP, SSHEndpoint, Status, Zone, NodeType)
- `total_count` computed field from API response
**Alternatives**: Could make cluster_id Optional and query all nodes across all clusters, but the API requires ClusterId - no alternative exists.

### Decision 2: Service Layer - Add method to existing DbdcService
**Choice**: Add `DescribeDBCustomClusterNodesByFilter` method to the existing `DbdcService` struct in `service_tencentcloud_dbdc.go`
**Rationale**: Consistent with how `DescribeDBCustomClustersByFilter` was added. Uses `paramMap` pattern for flexible parameter passing. Pagination handled internally with Limit=100.
**Alternatives**: Create a separate service file - rejected as it fragments the service layer for the same product.

### Decision 3: Filter handling - Use dbdc Filter type
**Choice**: Use `dbdcv20201029.Filter` struct (with Name and Values fields) for the filters parameter
**Rationale**: The `DescribeDBCustomClusterNodes` API accepts `[]*Filter` type, same as the clusters API. This is the dbdc product's standard filter format. The supported filter key is `node-name` for DB Custom node name.
**Alternatives**: None - this matches the SDK type.

### Decision 4: Test approach - gomonkey mock
**Choice**: Use gomonkey-based mock tests (not Terraform acceptance tests)
**Rationale**: Following the project requirement for new resources to use mock-based unit tests. The existing `data_source_tc_dbdc_db_custom_clusters_test.go` already uses this pattern.
**Alternatives**: Terraform acceptance tests - rejected per project rules for new resources.

## Risks / Trade-offs

- [API requires ClusterId] → cluster_id is marked as Required in schema, making the datasource dependent on having a known cluster ID first. Users must first discover clusters via `tencentcloud_dbdc_db_custom_clusters` then query nodes.
- [Filter only supports node-name] → Limited filtering capability compared to the clusters datasource which supports cluster-name, cluster-status, and tags. Documentation should clearly state supported filter names.
- [No nested structures in DBCustomClusterNode] → All 7 fields are flat string types, simplifying the implementation but limiting detail (no sub-structs for disks, etc. unlike DBCustomNode which is a different API).
