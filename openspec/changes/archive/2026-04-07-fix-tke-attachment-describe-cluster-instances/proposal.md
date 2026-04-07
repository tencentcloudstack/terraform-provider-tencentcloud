# 变更提案：重构 TKE ClusterAttachment 的 DescribeClusterInstances 调用方式

## 变更类型

**Bug 修复 / 性能优化** — 针对 `tencentcloud_kubernetes_cluster_attachment` 资源的两处 `DescribeClusterInstances` 调用进行重构：
1. 全量查询接口 `DescribeClusterInstances` 改为 for 循环 + retry 方式
2. 等待节点就绪的轮询从全量查询改为按 instanceId 精确查询

## Why

### 问题描述

**问题 1：`DescribeClusterInstances` 全量查询未做 retry**

`resourceTencentCloudKubernetesClusterAttachmentCreatePostFillRequest0` 通过 `nodeHasAttachedToCluster` 调用 `DescribeClusterInstances`，该函数内部使用 `goto getMoreData` 方式分页，但整个分页循环没有 retry 保护——网络抖动时整个函数会立即失败。同时，`nodeHasAttachedToCluster` 中的重试逻辑是先调用一次，失败后才 retry，模式不够统一。

**问题 2：等待节点初始化使用全量查询，效率低**

`resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0` 中等待 TKE 节点 init 成功时，使用 `tkeService.DescribeClusterInstances(ctx, clusterId)` 拉取集群所有节点，然后遍历匹配 `instanceId`。当集群节点数量多时，这会：
- 产生大量不必要的数据传输
- 分页查询多次，增加 API 调用次数
- 每轮 Retry 都做全量拉取

### 改进方案

1. **新增 `DescribeAllClusterInstances` 方法**（`service_tencentcloud_tke.go`）：将原 `DescribeClusterInstances` 的分页逻辑改造为 for 循环 + `resource.Retry` 包裹每次分页请求，参考 `DescribeClusterInstancesByRole` 的 for 循环风格，避免 `goto` 语句

2. **新增 `DescribeClusterInstanceById` 方法**（`service_tencentcloud_tke.go`）：通过 `InstanceIds` 字段传入单个 instanceId 进行精确查询，返回单个 `*InstanceInfo`，替代全量查询后遍历匹配的方式

3. **更新 `nodeHasAttachedToCluster`**：改为调用 `DescribeAllClusterInstances`

4. **更新 `resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0`**：等待节点就绪的 Retry 逻辑中改为调用 `DescribeClusterInstanceById`，判断逻辑（failed/running/retrying）保持不变

## What Changes

### 新增 service 方法

| 方法 | 说明 |
|------|------|
| `DescribeAllClusterInstances(ctx, id)` | 全量分页查询，for 循环 + 每次分页内 retry |
| `DescribeClusterInstanceById(ctx, clusterId, instanceId)` | 按 instanceId 精确查询，返回 `(*InstanceInfo, error)` |

### 修改位置

| 文件 | 修改内容 |
|------|---------|
| `service_tencentcloud_tke.go` | 新增上述两个方法 |
| `resource_tc_kubernetes_cluster_attachment_extension.go` | `nodeHasAttachedToCluster` 改调 `DescribeAllClusterInstances`；`CreatePostHandleResponse0` 改调 `DescribeClusterInstanceById` |

### 向后兼容性

✅ 完全向后兼容：
- 原 `DescribeClusterInstances` 保持不变，其他调用方不受影响
- 判断节点 init 成功/失败的逻辑完全不变，仅数据来源从全量列表改为精确查询
