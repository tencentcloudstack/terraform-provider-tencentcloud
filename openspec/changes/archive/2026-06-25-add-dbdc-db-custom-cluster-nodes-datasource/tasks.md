## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomClusterNodesByFilter` 方法，支持 cluster_id、filters 参数转换，内部自动分页（Limit=100），返回 `[]*dbdcv20201029.DBCustomClusterNode` 和 `totalCount`
- [x] 1.2 在 service 方法 retry 块内检查空响应（response == nil || Response == nil || NodeSet == nil），返回 NonRetryableError；失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`

## 2. Datasource Schema & Read

- [x] 2.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.go`，定义 schema：`cluster_id`（Required String）、`filters`（Optional TypeList，含 name + values）、`result_output_file`（Optional String）、`node_set`（Computed TypeList，含 node_id/node_name/lan_ip/ssh_endpoint/status/zone/node_type）、`total_count`（Computed TypeInt）
- [x] 2.2 实现 `dataSourceTencentCloudDbdcDbCustomClusterNodesRead` 函数：组装 paramMap（cluster_id + filters），调用 service 方法（retry 包装），扁平化 node_set 结果，设置 total_count，设置 `d.SetId(helper.BuildToken())`，处理 result_output_file

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_dbdc_db_custom_cluster_nodes": dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodes()`
- [x] 3.2 在 `tencentcloud/provider.md` 的 DBDC Data Source 部分添加 `tencentcloud_dbdc_db_custom_cluster_nodes`

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.md`，包含一句话描述（Use this data source to query...）、Example Usage（按 cluster_id 查询、按 filters 查询的示例）

## 5. Unit Tests

- [x] 5.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes_test.go`，使用 gomonkey mock 模式实现 `TestDbdcDbCustomClusterNodesDS_ReadBasic`（验证 node_set 字段映射）、`TestDbdcDbCustomClusterNodesDS_Schema`（验证 schema 结构）、`TestDbdcDbCustomClusterNodesDS_ReadWithEmptyResponse`（验证空响应返回错误）
- [x] 5.2 使用 `go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClusterNodesDS" -v -count=1 -gcflags="all=-l"` 验证测试通过
