## Why

Terraform Provider for TencentCloud currently lacks a datasource for querying DB Custom node list in the dbdc (Database Dedicated Cluster) product. Users need to query DB Custom nodes by node IDs, filters, or tags to reference existing node attributes in their Terraform configurations. Adding this datasource enables users to discover and reference DB Custom nodes without hardcoding values.

## What Changes

- Add a new RESOURCE_KIND_DATASOURCE `tencentcloud_dbdc_db_custom_nodes` to the Terraform provider
- Implement the Read operation using the `DescribeDBCustomNodes` cloud API (dbdc v20201029)
- Support query parameters: `node_ids`, `filters`, `tags` as optional input filters
- Expose `node_set` as computed output containing the list of DB Custom nodes with all their attributes
- Register the new datasource in `tencentcloud/provider.go`
- Add corresponding documentation in `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.md`

## Capabilities

### New Capabilities
- `dbdc-db-custom-nodes-datasource`: A datasource that queries DB Custom node list via the DescribeDBCustomNodes API, supporting filtering by node_ids, filters, and tags, and exposing the full node set attributes.

### Modified Capabilities
(No existing capabilities are modified)

## Impact

- New files: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.go`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes_test.go`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.md`
- Modified files: `tencentcloud/provider.go` (datasource registration), `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (service layer method for DescribeDBCustomNodes)
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` (DescribeDBCustomNodes)
