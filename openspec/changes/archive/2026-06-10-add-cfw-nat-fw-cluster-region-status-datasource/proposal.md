## Why

用户需要通过 Terraform 查询 CFW（云防火墙）NAT 防火墙引流集群的地域状态，以便了解各地域的集群部署情况和引流网络配置，支持基础设施即代码的管理模式。

## What Changes

- 新增数据源 `tencentcloud_cfw_nat_fw_cluster_region_status`，调用 `DescribeNatFwClusterRegionStatus` 接口查询 NAT 防火墙引流集群地域状态
- 支持通过 `nat_cluster_region_status_query_list` 参数过滤查询，返回 `total`（地域数量）和 `region_fw_status`（地域防火墙集群状态列表）

## Capabilities

### New Capabilities

- `cfw-nat-fw-cluster-region-status-datasource`: 新增 CFW NAT 防火墙引流集群地域状态数据源，支持按 CCN ID、NAT 网关 ID、资产类型、引流路由方法等条件查询各地域的集群部署状态和引流网络 CIDR

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status.go`
- 新增文件：`tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status_test.go`
- 新增文件：`tencentcloud/services/cfw/data_source_tc_cfw_nat_fw_cluster_region_status.md`
- 修改文件：`tencentcloud/provider.go`（注册数据源）
- 修改文件：`tencentcloud/provider.md`（添加数据源文档引用）
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904`
