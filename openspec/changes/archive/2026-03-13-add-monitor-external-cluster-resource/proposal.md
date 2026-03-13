# Change: 新增 Monitor External Cluster 资源

## Why
目前 terraform-provider-tencentcloud 缺少对 Prometheus 外部集群（External Cluster）的管理能力。用户需要能够通过 Terraform 将外部 Kubernetes 集群注册到腾讯云 TMP（Tencent Managed Prometheus）实例中，实现外部集群的监控管理。

## What Changes
- 新增资源: `tencentcloud_monitor_external_cluster`
- 支持将外部 K8s 集群注册到腾讯云 TMP 实例
- 支持查询已注册的外部集群信息
- 支持解除外部集群与 TMP 实例的关联

## Impact
- 新增文件: `tencentcloud/services/monitor/resource_tc_monitor_external_cluster.go`
- 修改文件: `tencentcloud/provider.go` (注册新资源)
- 新增测试文件: `tencentcloud/services/monitor/resource_tc_monitor_external_cluster_test.go`
- 新增文档: `website/docs/r/monitor_external_cluster.html.markdown`
- 受影响的包: `tencentcloud/services/monitor`
- 依赖服务: Monitor SDK v20180724 (现有依赖，无需新增)

## API Details

### 创建接口 (CreateExternalCluster)
- **接口名称**: CreateExternalCluster
- **接口文档**: https://cloud.tencent.com/document/api/248/118983
- **主要参数**:
  - `InstanceId` (必选): TMP 实例 ID
  - `ClusterRegion` (必选): 集群所在地域
  - `ClusterName` (可选): 集群名称
  - `ClusterId` (可选): 集群 ID
  - `ExternalLabels` (可选): 外部标签数组
  - `OpenDefaultRecord` (可选): 是否开启预聚合规则
  - `EnableExternal` (可选): 是否开启公网
- **返回值**: `ClusterId` (集群 ID)

### 查询接口 (DescribePrometheusClusterAgents)
- **接口名称**: DescribePrometheusClusterAgents
- **接口文档**: https://cloud.tencent.com/document/api/248/86040
- **主要参数**:
  - `InstanceId` (必选): TMP 实例 ID
  - `ClusterIds` (可选): 集群 ID 列表（字符串数组）
  - `Offset` (可选): 偏移量
  - `Limit` (可选): 返回数量
- **返回字段**: 
  - `ClusterId`: 集群 ID
  - `ClusterType`: 集群类型 (需要作为 computed 字段)
  - `ClusterName`: 集群名称
  - `Status`: 状态
  - 其他集群信息

### 删除接口 (DeletePrometheusClusterAgent)
- **接口名称**: DeletePrometheusClusterAgent
- **接口文档**: https://cloud.tencent.com/document/api/248/86041
- **主要参数**:
  - `InstanceId` (必选): TMP 实例 ID
  - `Agents` (必选): Agent 列表数组
    - `ClusterId`: 集群 ID
    - `ClusterType`: 集群类型

## Resource Design

### Resource ID 格式
- 格式: `{instanceId}#{clusterId}`
- 示例: `prom-abcd#ecls-1234`
- ID 由创建接口返回的 `InstanceId` 和 `ClusterId` 组成

### Schema 设计
```hcl
resource "tencentcloud_monitor_external_cluster" "example" {
  instance_id      = "prom-abcd"           # Required, ForceNew
  cluster_region   = "ap-shanghai"          # Required
  cluster_name     = "test-cluster"         # Optional
  cluster_id       = "ecls-abcd"           # Optional
  open_default_record = false               # Optional
  enable_external  = false                  # Optional
  
  external_labels {
    name  = "cluster_name"                  # Required
    value = "ai"                            # Optional
  }
  
  # Computed fields
  cluster_type = "external"                 # Computed
}
```

## Implementation Notes

### 重要注意事项
1. **资源唯一 ID**: CreateExternalCluster 接口返回 `ClusterId`，因此资源的唯一标识为 `instanceId#clusterId`
   
2. **查询接口映射**: DescribePrometheusClusterAgents 接口
   - `InstanceId` 参数: 从资源 ID 中解析的 `instanceId`
   - `ClusterIds` 参数: 从资源 ID 中解析的 `clusterId` (注意：是字符串列表类型，需要传递 `[clusterId]` 格式)
   - `ClusterType` 字段: 接口返回值中包含此字段，需要作为 `Computed` 字段保存到 schema 中

3. **删除接口参数构造**: DeletePrometheusClusterAgent 接口
   - `InstanceId`: 从资源 ID 解析获取
   - `Agents` 对象数组:
     - `Agents.ClusterId`: 从资源 ID 解析获取
     - `Agents.ClusterType`: 通过 `d.GetOk("cluster_type")` 从状态中获取（该字段在 Read 时已存储）

4. **参考实现**: 代码格式和结构严格参考 `tencentcloud_igtm_strategy` 资源
   - 文件位置: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
   - 遵循相同的代码组织模式、错误处理、重试逻辑和日志记录

5. **ForceNew 字段**: `instance_id` 应设置为 `ForceNew: true`，因为集群不能在实例之间迁移

## References
- 参考资源实现: `resource_tc_igtm_strategy.go`
- 类似资源: `resource_tc_monitor_tmp_tke_cluster_agent.go` (TKE 集群 Agent 管理)
- Monitor 服务层: `service_tencentcloud_monitor.go`
