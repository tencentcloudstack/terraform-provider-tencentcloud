# Design: 优化 RabbitMQ 实例的 update 逻辑

## Context

当前 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中的 update 函数将多个可以通过 ModifyRabbitMQVipInstance API 修改的参数标记为不可修改，包括：
- `enable_deletion_protection`（删除保护）
- `remark`（备注）
- `enable_risk_warning`（集群风险提示）
- `enable_public_access`（公网访问）
- `band_width`（公网带宽）

实际上，根据腾讯云 TDMQ API 文档，`ModifyRabbitMQVipInstance` API 支持以下参数的修改：
- `ClusterName`（集群名称）
- `Remark`（备注）
- `EnableDeletionProtection`（删除保护）
- `RemoveAllTags`（删除所有标签）
- `Tags`（标签信息，全量替换）
- `EnableRiskWarning`（集群风险提示）

对于公网访问相关参数，可能需要调用其他 API（如 `ModifyPublicNetworkAccess` 或类似接口），需要进一步验证。

## Goals / Non-Goals

**Goals:**
- 扩展 RabbitMQ VIP 实例的可修改参数范围
- 支持修改 `enable_deletion_protection`、`remark`、`enable_risk_warning` 参数
- 保持向后兼容性，不破坏现有配置和状态
- 提供清晰的错误提示，区分真正不可修改的基础设施参数

**Non-Goals:**
- 不修改 `zone_ids`、`vpc_id`、`subnet_id`、`node_spec`、`node_num`、`storage_size`、`cluster_version`、`auto_renew_flag`、`time_span`、`pay_mode` 等基础设施级别参数（这些参数确实不支持修改或需要重建实例）
- 不实现公网访问相关的修改（需要单独的 API 调用，不在本次优化范围内）
- 不修改资源的 schema（只新增 Optional 字段）

## Decisions

### 决策 1：新增三个可修改参数

**选择：** 将 `enable_deletion_protection`、`remark`、`enable_risk_warning` 从不可修改列表中移除，并在 update 函数中添加对这些参数的处理

**理由：**
- 这些参数在 ModifyRabbitMQVipInstance API 中明确支持
- 这些配置的修改不会影响实例的基础设施架构，符合"配置"而非"规格"的变更
- 用户可以通过控制台修改这些参数，Terraform 也应该支持

**考虑的替代方案：**
- 不做任何修改（当前状态） - ❌ 用户体验差，需要手动操作
- 将所有参数都标记为可修改 - ❌ 部分参数确实不支持修改，会导致 API 调用失败

### 决策 2：保持基础设施参数不可修改

**选择：** 继续将 `zone_ids`、`vpc_id`、`subnet_id`、`node_spec`、`node_num`、`storage_size`、`cluster_version`、`auto_renew_flag`、`time_span`、`pay_mode` 等参数标记为不可修改

**理由：**
- 这些参数涉及实例的基础设施规格，修改通常需要重建实例
- ModifyRabbitMQVipInstance API 不支持这些参数的修改
- 保持与腾讯云控制台行为一致

### 决策 3：参数映射逻辑

**选择：** 在 update 函数中，针对每个可修改参数单独检查 `d.HasChange()`，只有参数确实发生变化时才将其添加到 API 请求中

**理由：**
- 避免不必要的 API 调用
- 符合 Terraform 插件的最佳实践
- 减少 API 调用失败的风险

### 决策 4：新增 Optional schema 字段

**选择：** 在资源 schema 中新增 `remark` 和 `enable_risk_warning` 两个 Optional 字段（enable_deletion_protection 已在 schema 中）

**理由：**
- 这两个参数在 Create API 中存在，但在 schema 中可能未定义
- 需要确保 schema 包含所有需要管理的参数
- Optional 字段不会破坏现有配置

## Risks / Trade-offs

### 风险 1：API 参数支持不一致

**风险：** 可能存在某些环境或版本下 Modify API 不支持某些参数的情况

**缓解措施：**
- 在代码中添加适当的错误处理
- 如果 API 返回不支持该参数的错误，提供清晰的错误提示
- 在验收测试中验证参数修改功能

### 风险 2：资源 tags 修改逻辑

**风险：** 当前的 tags 修改逻辑可能存在边界情况处理不当的问题（如清空所有标签）

**缓解措施：**
- 保持当前的 `RemoveAllTags` 逻辑
- 确保当 `resource_tags` 被移除时正确设置 `RemoveAllTags` 为 true

### 风险 3：向后兼容性

**风险：** 新增 Optional 字段可能影响某些依赖旧版本 Provider 的用户

**缓解措施：**
- Optional 字段不会强制要求用户配置
- 现有配置继续正常工作
- Read 函数会自动读取并填充这些字段的值

## Migration Plan

### 部署步骤

1. 修改 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中的 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
2. 在资源 schema 中添加 `remark` 和 `enable_risk_warning` 字段（如果不存在）
3. 修改 Update 函数中的 `immutableArgs` 列表，移除新增的可修改参数
4. 添加对 `remark`、`enable_risk_warning`、`enable_deletion_protection` 参数的处理逻辑
5. 更新验收测试，覆盖新参数的修改场景
6. 更新资源文档，说明可修改的参数范围

### 回滚策略

如果出现问题，可以通过以下方式回滚：
1. 将新增的可修改参数重新添加到 `immutableArgs` 列表中
2. 移除对这些参数的处理逻辑
3. 保持 schema 中的字段（Optional 字段不会破坏现有配置）

## Open Questions

1. **问题：** 公网访问（`enable_public_access` 和 `band_width`）的修改是否需要调用其他 API？

   **状态：** 需要进一步验证腾讯云 API 文档

2. **问题：** `enable_deletion_protection` 参数是否已经在 schema 中存在？

   **状态：** 需要确认当前 schema 的完整定义

3. **问题：** 修改这些参数是否需要等待实例状态稳定？

   **状态：** 根据 Modify API 的异步特性，可能需要在修改后添加状态等待逻辑
