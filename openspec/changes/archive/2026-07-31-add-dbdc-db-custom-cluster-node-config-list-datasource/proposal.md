## Why

The `dbdc` service currently exposes data sources for listing DB Custom clusters,
nodes, images, and cluster node membership, but there is no way to discover the
Kubernetes scheduling configuration (labels and taints) attached to each node
inside a DB Custom cluster. Operators need this information to reason about pod
placement and to drive targeted scheduling decisions declaratively from
Terraform. The cloud API `DescribeDBCustomClusterNodeConfig` already returns
this data, so a new read-only data source closes the gap.

## What Changes

- Add a new data source `tencentcloud_dbdc_db_custom_cluster_node_config_list`
  that queries node configuration (labels / taints) for a given DB Custom
  cluster.
- The data source wraps the cloud API
  `DescribeDBCustomClusterNodeConfig` (`dbdc` SDK `v20201029`), which accepts a
  `ClusterId` and an optional list of `NodeIds` and returns a `NodeSet` of
  `DBCustomClusterNodeConfig` items.
- Expose `cluster_id` (Required) and `node_ids` (Optional) arguments.
- Expose a computed `node_set` list. Per the provider convention for
  Describe-style data sources, the list is **flattened**: each element exposes
  `node_id`, `labels` (list of `key`/`value`), and `taints` (list of
  `key`/`value`/`effect`) directly — no extra nesting layer.
- Add `result_output_file` (Optional) for writing results to a file, consistent
  with the sibling `tencentcloud_dbdc_db_custom_cluster_nodes` data source.
- Register the data source in `provider.go` DataSourcesMap and add a provider
  doc entry.
- Add a `data_source_tc_dbdc_db_custom_cluster_node_config_list.md` example
  file (generated into `website/docs/d/` via `make doc`).
- Add unit tests using `gomonkey` mocks (no Terraform acceptance suite).

## Capabilities

### New Capabilities
- `dbdc-db-custom-cluster-node-config-list-datasource`: Query DB Custom cluster
  node configuration (labels and taints) via the
  `DescribeDBCustomClusterNodeConfig` cloud API, exposing a flattened
  `node_set` of per-node scheduling config.

### Modified Capabilities
<!-- None. This is a purely additive data source; no existing spec requirements change. -->

## Impact

- **New files**:
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_config_list.go`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_config_list_test.go`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_config_list.md`
- **Modified files**:
  - `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — add
    `DescribeDBCustomClusterNodeConfigByFilter` service helper.
  - `tencentcloud/provider.go` — register
    `tencentcloud_dbdc_db_custom_cluster_node_config_list` in DataSourcesMap.
  - `tencentcloud/provider.md` — add provider doc entry (generated via
    `make doc`).
- **APIs**: `DescribeDBCustomClusterNodeConfig` (synchronous read, no async
  polling needed).
- **Dependencies**: Reuses existing vendored SDK
  `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`;
  no SDK upgrade required. `UseDbdcV20201029Client()` already exists.
- **Backward compatibility**: Purely additive; no existing schema or state
  changes.
