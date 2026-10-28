## 1. 数据源核心实现

- [x] 1.1 新增 `tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status.go`，实现数据源 Schema 定义（入参 `nat_cluster_region_status_query_list`，出参 `total`、`region_fw_status`）和 Read 函数，调用 `DescribeNatFwClusterRegionStatus` 接口，使用 `tccommon.ReadRetryTimeout` 超时和 retry 处理
- [x] 1.2 在 Read 函数中正确处理 nil 判断，云 API 返回空时在 retry 块内返回 `NonRetryableError`，避免清空 state

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加数据源引用

## 3. 文档

- [x] 3.1 新增 `tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status.md`，包含一句话描述、Example Usage（使用 jsonencode 处理 JSON 字段）

## 4. 单元测试

- [x] 4.1 新增 `tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status_test.go`，使用 gomonkey mock 云 API，覆盖主要业务逻辑，使用 `go test -gcflags=all=-l` 跑通测试
