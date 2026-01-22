# Change: Add TKE Cluster Admin Role Data Source

## Why

用户需要通过 Terraform 管理 TKE（容器服务）集群的 RBAC 管理员权限授予操作。当前 Provider 缺少对 `AcquireClusterAdminRole` 接口的支持，该接口用于给子账户授予集群的 `tke:admin` ClusterRole 权限。

虽然此接口是操作型接口而非查询接口，但用户明确需要通过 Data Source 的方式来触发此授权操作，以符合其基础设施即代码的使用场景。

## What Changes

- 新增 Data Source: `tencentcloud_kubernetes_cluster_admin_role`
- 实现对 TKE API `AcquireClusterAdminRole` 接口的调用
- 支持通过 `cluster_id` 参数触发集群管理员角色授予操作
- 返回操作的请求 ID 用于追溯

## Impact

- **新增能力**: TKE 集群管理员角色授予
- **受影响的服务**: TKE (tencentcloud/services/tke)
- **新增文件**:
  - `tencentcloud/services/tke/data_source_tc_kubernetes_cluster_admin_role.go`
  - `tencentcloud/services/tke/data_source_tc_kubernetes_cluster_admin_role_test.go`
  - Provider 注册代码需要添加此 data source
- **API 依赖**: 
  - TKE API v20180525: `AcquireClusterAdminRole`
  - 文档: https://cloud.tencent.com/document/product/457/49014
