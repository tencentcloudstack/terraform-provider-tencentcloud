# Technical Design Document

## Context
腾讯云 Prometheus 监控服务（TMP）支持管理外部 Kubernetes 集群的监控。当前 terraform-provider-tencentcloud 已经支持 TKE 集群的监控管理（通过 `tencentcloud_monitor_tmp_tke_cluster_agent`），但缺少对外部集群的支持。本次新增 `tencentcloud_monitor_external_cluster` 资源以填补这一功能空白。

## Goals / Non-Goals

### Goals
- 实现 `tencentcloud_monitor_external_cluster` 资源的完整 CRUD 操作
- 正确处理 API 参数格式差异（特别是 `ClusterIds` 字符串数组和 `Agents` 对象数组）
- 遵循项目代码规范和 `resource_tc_igtm_strategy.go` 的实现模式
- 提供完整的测试覆盖和文档

### Non-Goals
- 不实现外部集群的自动发现功能
- 不支持批量管理多个外部集群（每个资源管理一个集群）
- 不实现 TKE 集群的管理（已有专门资源）
- 不修改现有 Monitor 服务的其他资源

## Decisions

### Decision 1: 复合资源 ID 格式
**选择**: 使用 `{instanceId}#{clusterId}` 格式

**原因**:
- `CreateExternalCluster` API 返回 `ClusterId`，而 `InstanceId` 是必需的输入参数
- 两者组合才能唯一标识一个外部集群资源
- 符合项目中复合 ID 的通用模式（如 `resource_tc_igtm_strategy.go` 使用 `instanceId#strategyId`）
- 分隔符 `#` 已在项目中广泛使用（`tccommon.FILED_SP` 常量）

**替代方案考虑**:
- 仅使用 `ClusterId`: 不可行，因为同一个 `ClusterId` 可能在不同 Instance 中重复
- 使用其他分隔符: 不必要，`#` 已是项目标准

### Decision 2: ClusterType 作为 Computed 字段
**选择**: `cluster_type` 字段设置为 `Computed: true`

**原因**:
- `CreateExternalCluster` API 不接受 `ClusterType` 参数
- `DescribePrometheusClusterAgents` API 返回 `ClusterType` 字段
- Delete 操作需要 `ClusterType` 值
- 设置为 Computed 字段可以在 Read 时自动填充，Delete 时从状态获取

**实现细节**:
```go
"cluster_type": {
    Type:        schema.TypeString,
    Computed:    true,
    Description: "Cluster type, returned by API.",
}
```

### Decision 3: DescribePrometheusClusterAgents 查询策略
**选择**: 使用 `ClusterIds` 参数精确查询

**原因**:
- API 支持通过 `ClusterIds` 数组过滤结果，提高查询效率
- 避免分页遍历所有集群后再过滤
- 注意 `ClusterIds` 是字符串数组类型，需要传递 `[clusterId]` 格式

**实现示例**:
```go
request.InstanceId = &instanceId
clusterIds := []*string{&clusterId}
request.ClusterIds = clusterIds
```

### Decision 4: 不支持 Update 操作
**选择**: Update 函数返回错误或调用 Read（取决于是否有 Modify API）

**原因**:
- 需要检查腾讯云是否提供 `ModifyExternalCluster` 或类似 API
- 如果没有 Modify API，则所有字段应设置为 `ForceNew` 或在 Update 中返回不支持错误
- 参考 `resource_tc_monitor_tmp_tke_cluster_agent.go`，只支持更新 `external_labels`

**备选实现**:
- 如果发现有 Modify API，实现完整的 Update 逻辑
- 如果没有，Update 函数直接返回 Read 结果

### Decision 5: 参考 igtm_strategy 实现模式
**选择**: 代码结构严格遵循 `resource_tc_igtm_strategy.go`

**原因**:
- 用户明确要求参考此文件
- 该文件展示了标准的资源实现模式
- 包含正确的错误处理、重试逻辑、日志记录
- 复合 ID 的解析和使用模式一致

**关键代码模式**:
```go
// Create: 构造复合 ID
d.SetId(strings.Join([]string{instanceId, clusterId}, tccommon.FILED_SP))

// Read/Update/Delete: 解析复合 ID
idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
if len(idSplit) != 2 {
    return fmt.Errorf("id is broken,%s", d.Id())
}
instanceId := idSplit[0]
clusterId := idSplit[1]
```

## Risks / Trade-offs

### Risk 1: API 字段不一致
**风险**: CreateExternalCluster 的参数名可能与文档略有差异

**缓解措施**:
- 参考 tencentcloud-sdk-go 的类型定义
- 在测试中验证所有参数
- 添加详细的错误日志记录

### Risk 2: ClusterType 值的变化
**风险**: `ClusterType` 值可能随 API 版本变化

**缓解措施**:
- 不在代码中硬编码 `ClusterType` 值
- 完全依赖 API 返回值
- 在文档中说明此字段由 API 自动填充

### Risk 3: Delete 操作失败处理
**风险**: 如果 `cluster_type` 在状态中缺失，Delete 操作可能失败

**缓解措施**:
- 在 Read 操作中确保 `cluster_type` 正确保存
- 在 Delete 前检查字段是否存在
- 提供清晰的错误信息指导用户

**实现**:
```go
clusterType, ok := d.GetOk("cluster_type")
if !ok {
    return fmt.Errorf("cluster_type not found in state, cannot delete cluster")
}
```

### Trade-off: 简单实现 vs 完整功能
**选择**: 首先实现核心 CRUD 功能，扩展功能在后续迭代

**权衡**:
- ✓ 快速交付核心功能
- ✓ 降低初始实现复杂度
- ✗ 可能需要后续增加字段

**计划**: 初始版本只实现 API 文档中明确的字段，后续根据用户需求扩展

## Migration Plan

### Phase 1: 实现和测试 (Week 1)
1. 实现资源代码
2. 编写单元测试
3. 在测试环境验证

### Phase 2: 文档和审查 (Week 1)
1. 编写使用文档
2. 代码审查
3. 修复反馈问题

### Phase 3: 发布 (Week 2)
1. 合并到主分支
2. 发布新版本
3. 更新 CHANGELOG

### Rollback Plan
如果发现严重问题:
1. 从 Provider 注册中移除资源
2. 发布 hotfix 版本
3. 在文档中标记资源为废弃（如果已有用户使用）

## Open Questions

1. **Q**: 是否存在 `ModifyExternalCluster` API？
   **A**: 需要查阅完整 API 文档或 SDK 确认，如果不存在，所有字段应为 ForceNew

2. **Q**: `ClusterIds` 参数是否支持空数组？
   **A**: 需要测试，可能需要处理不传递此参数的情况（查询所有集群后过滤）

3. **Q**: Delete 操作是否需要等待集群状态变化？
   **A**: 参考 `resource_tc_monitor_tmp_tke_cluster_agent.go`，可能需要轮询等待状态变为删除完成

4. **Q**: 外部标签是否支持更新？
   **A**: 如果有 `ModifyPrometheusAgentExternalLabels` API（类似 TKE Agent），可以支持

## Technical Details

### API 调用流程

#### Create 流程:
```
1. 用户定义 Terraform 配置
2. Provider 验证必需字段
3. 调用 CreateExternalCluster API
4. API 返回 ClusterId
5. 构造 ID: instanceId#clusterId
6. 调用 Read 填充完整状态
```

#### Read 流程:
```
1. 解析 ID 获取 instanceId 和 clusterId
2. 调用 DescribePrometheusClusterAgents API
   - 参数: InstanceId, ClusterIds=[clusterId]
3. 在返回的 Agents 列表中查找匹配的 clusterId
4. 提取 ClusterType 等字段
5. 更新 Terraform 状态
```

#### Delete 流程:
```
1. 解析 ID 获取 instanceId 和 clusterId
2. 从状态获取 cluster_type
3. 构造 Agents 参数:
   Agents = [{
     ClusterId: clusterId,
     ClusterType: clusterType
   }]
4. 调用 DeletePrometheusClusterAgent API
5. 可选: 轮询验证删除完成
```

### 错误处理策略

| 错误场景 | 处理方式 |
|---------|---------|
| API 限流 | 使用 `resource.Retry` 和 `tccommon.WriteRetryTimeout` 重试 |
| 参数错误 | 立即返回，不重试 |
| 资源不存在 (Read) | 设置 ID 为空，标记为已删除 |
| 网络临时错误 | 重试，记录详细日志 |
| 资源冲突 | 根据 API 错误码决定是否重试 |

### 测试策略

#### 单元测试:
- ID 解析和构造逻辑
- Schema 验证
- 参数映射正确性

#### 集成测试 (Acceptance Tests):
```go
func TestAccTencentCloudMonitorExternalCluster_basic(t *testing.T) {
    // 测试创建、读取、删除基本流程
}

func TestAccTencentCloudMonitorExternalCluster_full(t *testing.T) {
    // 测试所有参数（包含可选字段）
}

func TestAccTencentCloudMonitorExternalCluster_import(t *testing.T) {
    // 测试 terraform import 功能
}
```

## References
- [CreateExternalCluster API 文档](https://cloud.tencent.com/document/api/248/118983)
- [DescribePrometheusClusterAgents API 文档](https://cloud.tencent.com/document/api/248/86040)
- [DeletePrometheusClusterAgent API 文档](https://cloud.tencent.com/document/api/248/86041)
- 参考实现: `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
- 类似资源: `tencentcloud/services/tmp/resource_tc_monitor_tmp_tke_cluster_agent.go`
