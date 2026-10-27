## Why

Terraform Provider for TencentCloud 目前缺少查询 CDWDoris（腾讯云数据仓库 Doris）集群节点信息的能力。用户在编排 Doris 集群时无法通过 terraform data source 获取集群下各节点（FE/BE）的 IP、规格、磁盘、可用区等明细，不便于做容量规划、监控配置与节点级编排。云 API 已提供 `DescribeInstanceNodes` 接口获取集群节点信息列表，需要为其封装一个数据源资源，使这批信息可在 Terraform 中被直接消费。

## What Changes

- 新增数据源资源 `tencentcloud_cdwdoris_instance_nodes`（RESOURCE_KIND_DATASOURCE），文件位于 `tencentcloud/services/cdwdoris/data_source_tc_cdwdoris_instance_nodes.go`，仅实现 Read 操作，调用云 API `DescribeInstanceNodes` 按 `instance_id` 查询集群节点列表。
- 新增 service 层方法 `DescribeCdwdorisInstanceNodesByFilter`（位于 `tencentcloud/services/cdwdoris/service_tencentcloud_cdwdoris.go`），负责构造请求、调用 `DescribeInstanceNodes` 接口并做分页拉取。
- 新增数据源测试文件 `data_source_tc_cdwdoris_instance_nodes_test.go`，使用 gomonkey mock 云 API 进行业务逻辑单元测试。
- 新增文档文件 `data_source_tc_cdwdoris_instance_nodes.md`，供 `make doc` 生成 website 文档。
- 在 `tencentcloud/provider.go` 中注册数据源 `tencentcloud_cdwdoris_instance_nodes`，并在 `tencentcloud/provider.md` 对应资源列表中补充该数据源条目。

## Capabilities

### New Capabilities
- `cdwdoris-instance-nodes-datasource`: 通过云 API `DescribeInstanceNodes` 查询指定 CDWDoris 集群实例下的节点信息列表，支持按 `node_role` 和 `display_policy` 过滤，输出每个节点的 IP、规格、CPU、内存、磁盘、角色、状态、可用区等明细字段，以及集群支持的节点角色列表 `node_roles`。

### Modified Capabilities
<!-- 无现有 capability 的 spec 级要求变更 -->

## Impact

- 新增代码文件：`data_source_tc_cdwdoris_instance_nodes.go`、`data_source_tc_cdwdoris_instance_nodes_test.go`、`data_source_tc_cdwdoris_instance_nodes.md`。
- 修改文件：`tencentcloud/services/cdwdoris/service_tencentcloud_cdwdoris.go`（新增 service 方法）、`tencentcloud/provider.go`（注册数据源）、`tencentcloud/provider.md`（补充数据源文档条目）。
- 依赖：`DescribeInstanceNodes` 接口已在 vendor 中存在（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228`），request/response struct 及 `InstanceNode` struct 均已就绪，无需额外升级 vendor。
- 向后兼容：仅新增数据源，不影响任何现有资源与 state，完全向后兼容。