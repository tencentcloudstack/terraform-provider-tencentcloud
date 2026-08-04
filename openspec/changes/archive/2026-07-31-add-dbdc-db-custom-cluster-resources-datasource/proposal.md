## Why

用户需要通过 Terraform 查询 DB Custom 集群的资源信息（节点总数、CPU/内存/Pods 的容量、可分配量、已申请量、上限和可用余量）。当前 `dbdc` 服务已提供集群、节点、镜像等数据源，但缺少查询集群资源汇总的数据源，用户无法在 Terraform 配置中动态获取集群的资源使用情况用于容量规划和调度决策。

## What Changes

- 新增 Data Source: `tencentcloud_dbdc_db_custom_cluster_resources`
- 实现对 dbdc API `DescribeDBCustomClusterResources` 接口的调用，通过 `cluster_id` 查询指定集群的资源信息
- 返回集群资源汇总字段：
  - `node_count`: 参与汇总的工作节点总数（不含控制面节点）
  - `capacity`: 集群所有节点的资源物理总容量之和（嵌套对象，含 `cpu`、`memory`、`pods`）
  - `allocatable`: 集群所有节点的可分配容量之和（嵌套对象，含 `cpu`、`memory`、`pods`）
  - `requests`: 集群所有非终态 Pod 的 requests 申请量之和（嵌套对象，含 `cpu`、`memory`、`pods`）
  - `limits`: 集群所有非终态 Pod 的 limits 上限之和（嵌套对象，含 `cpu`、`memory`、`pods`）
  - `available`: 集群可再调度余量（嵌套对象，含 `cpu`、`memory`、`pods`）
- 在 `provider.go` 中注册该数据源
- 在 `provider.md` 中添加该数据源的声明

## Capabilities

### New Capabilities
- `dbdc-db-custom-cluster-resources-datasource`: 查询 DB Custom 集群资源信息的数据源，支持通过 `cluster_id` 获取集群的资源容量、可分配量、申请量、上限和可用余量等汇总信息

### Modified Capabilities

（无）

## Impact

- **新增能力**: DB Custom 集群资源信息查询
- **受影响的服务**: dbdc (tencentcloud/services/dbdc)
- **新增文件**:
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources.go`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources.md`
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources_test.go`
- **修改文件**:
  - `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — 新增 `DescribeDBCustomClusterResources` service 层方法
  - `tencentcloud/provider.go` — 注册数据源
  - `tencentcloud/provider.md` — 添加数据源声明
- **API 依赖**:
  - dbdc API v20201029: `DescribeDBCustomClusterResources`
- **兼容性**: 无破坏性变更，纯新增功能
