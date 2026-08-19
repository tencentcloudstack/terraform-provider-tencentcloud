## Why

The `dbdc` service exposes node-level resource usage (CPU / Memory / Pods for
Capacity, Allocatable, Requests, Limits, Available) through the cloud API
`DescribeDBCustomClusterNodeResources`, but the Terraform provider has no data
source to query it. Users cannot discover a DB Custom node's current resource
allocation and scheduling headroom declaratively, which limits capacity
planning and node selection within a cluster. Adding a datasource closes this
gap.

## What Changes

- Add a new Terraform data source
  `tencentcloud_dbdc_db_custom_cluster_node_resources` (RESOURCE_KIND_DATASOURCE,
  read-only query).
- It wraps the synchronous cloud API
  `DescribeDBCustomClusterNodeResources` (`dbdc/v20201029`) to query node
  resource information within a DB Custom cluster.
- Input arguments: `cluster_id` (required), `node_ids` (optional, list of node
  IDs, max 50 per request).
- Output: `node_set` list of `DBCustomClusterNodeResource` items, each with
  `node_id` plus five nested `MetaResource` blocks (`capacity`, `allocatable`,
  `requests`, `limits`, `available`) containing `cpu`, `memory`, `pods`.
- Register the data source in `tencentcloud/provider.go` DataSourcesMap and add
  an entry to `tencentcloud/provider.md`.
- Add service-layer helper `DescribeDBCustomClusterNodeResourcesByFilter` in
  `service_tencentcloud_dbdc.go` (synchronous, no pagination — the API filters
  via `NodeIds`).
- Add a gomonkey-mock unit test and a `.md` example doc.

## Capabilities

### New Capabilities
- `dbdc-db-custom-cluster-node-resources-datasource`: Query DB Custom cluster node resource information (Capacity / Allocatable / Requests / Limits / Available) by `cluster_id` and optional `node_ids`.

### Modified Capabilities
<!-- None. This is a new, independent data source. -->

## Impact

- **New files** (all under `tencentcloud/services/dbdc/`):
  - `data_source_tc_dbdc_db_custom_cluster_node_resources.go` — schema + Read
  - `data_source_tc_dbdc_db_custom_cluster_node_resources_test.go` — gomonkey mock unit tests
  - `data_source_tc_dbdc_db_custom_cluster_node_resources.md` — example usage doc
- **Modified files**:
  - `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — new helper
    `DescribeDBCustomClusterNodeResourcesByFilter`
  - `tencentcloud/provider.go` — register data source in DataSourcesMap
  - `tencentcloud/provider.md` — add data source entry
- **APIs / dependencies**:
  - Uses vendored SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`.
  - `DescribeDBCustomClusterNodeResources` is already present in the vendored
    SDK; `UseDbdcV20201029Client()` already exists. No SDK upgrade required.
- **No breaking changes** — purely additive.
