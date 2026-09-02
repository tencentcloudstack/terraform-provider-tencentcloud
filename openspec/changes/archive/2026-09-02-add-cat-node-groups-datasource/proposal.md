## Why

The Terraform Provider for TencentCloud currently exposes `tencentcloud_cat_node`, `tencentcloud_cat_probe_data`, and `tencentcloud_cat_metric_data` data sources for the CAT (Cloud Automated Testing / 拨测) product, but there is no way to query the available **probe node groups** (拨测点组). Probe node groups — including availability probe groups, advanced probe groups, and user-defined "my probe groups" — are the top-level organizing unit for selecting dial-test nodes when creating probe tasks. Without a data source for node groups, users must hard-code node group IDs into their Terraform configurations, which is fragile and not portable across accounts or regions. Adding a `tencentcloud_cat_node_groups` data source lets users dynamically discover probe node groups (filtered by node type, task category, IP type, region, district, ISP, etc.) and feed them into probe task resources.

## What Changes

- Add a new data source `tencentcloud_cat_node_groups` backed by the cloud API `DescribeNodeGroups` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409`).
- Expose 10 optional filter arguments: `node_type`, `task_category`, `ip_type`, `name`, `region_id`, `district_id`, `net_service_id`, `node_group_type`, `task_type`, `probe_type`.
- Expose three computed output collections:
  - `node_list` — the tree of probe node groups (two-level: `NodeTree` → `Children` (`NodeLeaf`) → `Children` (`NodeInfoBase`)), each level with `id`, `content`, and nested `children`.
  - `district_list` — province/country list, each with `id` and `name`.
  - `net_service_list` — ISP list, each with `id` and `name`.
- Register the new data source in `tencentcloud/provider.go` (`DataSourcesMap`) and update the provider doc-comment index for `make doc`.
- Add the hand-written example file `tencentcloud/services/cat/data_source_tc_cat_node_groups.md`.
- Add the service-layer method `DescribeCatNodeGroupsByFilter` to `tencentcloud/services/cat/service_tencentcloud_cat.go` (single API call, no pagination — `DescribeNodeGroups` returns the full tree).
- Add the data source file `tencentcloud/services/cat/data_source_tc_cat_node_groups.go`.
- Add gomonkey-based unit tests in `tencentcloud/services/cat/data_source_tc_cat_node_groups_test.go` (mock the SDK client; no Terraform acceptance suite).

## Capabilities

### New Capabilities

- `cat-node-groups-datasource`: A data source that queries CAT probe node groups via `DescribeNodeGroups`, supports filter arguments (node type, task category, IP type, name, region, district, ISP, node group type, task type, probe type), and flattens the returned `node_list` tree plus `district_list` and `net_service_list` into Terraform state.

### Modified Capabilities

<!-- None — this is a purely additive data source. No existing spec-level behavior changes. -->

## Impact

- **New files**:
  - `tencentcloud/services/cat/data_source_tc_cat_node_groups.go`
  - `tencentcloud/services/cat/data_source_tc_cat_node_groups_test.go`
  - `tencentcloud/services/cat/data_source_tc_cat_node_groups.md`
- **Modified files**:
  - `tencentcloud/services/cat/service_tencentcloud_cat.go` — add `DescribeCatNodeGroupsByFilter` method.
  - `tencentcloud/provider.go` — register `tencentcloud_cat_node_groups` in `DataSourcesMap` and update the doc-comment index.
- **Cloud API dependency**: `DescribeNodeGroups` in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409` (already vendored).
- **Documentation**: `make doc` will generate `website/docs/d/cat_node_groups.html.markdown` and update `website/tencentcloud.erb`.
- **Backward compatibility**: Fully additive; no existing resources or data sources are modified. No breaking changes.
