## Why

Terraform Provider for TencentCloud currently has no datasource for DB Custom Cluster (dbdc product). Users need to query DB Custom cluster list information in their Terraform configurations to reference existing clusters, filter clusters by name/status, and integrate cluster data into other resource definitions. Adding `tencentcloud_dbdc_db_custom_clusters` datasource enables users to read and query DB Custom cluster attributes within Terraform workflows.

## What Changes

- Add new datasource `tencentcloud_dbdc_db_custom_clusters` (RESOURCE_KIND_DATASOURCE) that calls `DescribeDBCustomClusters` API to query DB Custom cluster list
- Add datasource schema with query parameters: `cluster_ids`, `filters`, `tags`
- Add computed output fields: `cluster_set` (list of cluster details), `total_count`
- Add service layer for dbdc product (`service_tencentcloud_dbdc.go`) with `DescribeDBCustomClustersByFilter` method
- Register the datasource in `provider.go` and `provider.md`
- Add documentation file `data_source_tc_dbdc_db_custom_clusters.md`
- Add unit test file `data_source_tc_dbdc_db_custom_clusters_test.go`

## Capabilities

### New Capabilities
- `dbdc-db-custom-clusters-datasource`: Datasource to query DB Custom cluster list from dbdc product, supporting filtering by cluster IDs, filter conditions (cluster-name, cluster-status), and tags. Returns cluster details including cluster_id, cluster_name, region, cluster_level, cluster_status, cluster_version, cluster_node_num, cluster_description, created_time, and tags.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- New files: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters.go`, `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters.md`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters_test.go`
- Modified files: `tencentcloud/provider.go` (register datasource), `tencentcloud/provider.md` (add datasource entry)
- API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029.DescribeDBCustomClusters`
