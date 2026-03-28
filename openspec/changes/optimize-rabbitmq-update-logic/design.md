## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 实现将大量字段标记为不可变，包括 `node_spec`、`node_num`、`storage_size`、`band_width`、`enable_public_access`、`cluster_version` 等规格相关字段。这些字段实际上可能通过腾讯云 API 的 `ModifyRabbitMQVipInstance` 接口进行修改。

当前实现的问题：
1. 用户无法在不重建实例的情况下调整节点规格、数量或存储容量，这在云资源管理中是常见需求
2. 代码中的 `immutableArgs` 列表过于保守，可能没有充分测试 API 的实际能力
3. 不支持规格更新会导致用户资源管理效率低下，成本控制困难

相关代码位置：
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`: 第 450-523 行的 update 函数
- 不可变字段列表在第 460-465 行定义

硬约束：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 不能修改已有资源的 schema（除非只新增 Optional 字段）
- 需要通过 tencentcloud-sdk-go 调用云 API
- 需要正确处理异步操作的最终一致性

## Goals / Non-Goals

**Goals:**
- 支持通过 `ModifyRabbitMQVipInstance` API 更新 RabbitMQ 实例的规格参数
- 明确区分可更新字段和不可更新字段，提供清晰的用户反馈
- 确保更新操作的正确性和一致性，包括错误处理和重试逻辑
- 保持向后兼容性，不影响现有资源的正常运行

**Non-Goals:**
- 修改实例的基础架构属性（可用区、VPC、子网）- 这些字段确实需要重建
- 修改付费模式或购买时长 - 这些字段在实例创建后不应变更
- 修改 schema 定义或添加新字段 - 仅更新现有字段的可更新性
- 添加新的 API 调用或服务层逻辑 - 仅使用现有的 `ModifyRabbitMQVipInstance` API

## Decisions

### 决策 1: 将规格相关字段从不可变列表中移除

**选择:** 从 `immutableArgs` 列表中移除 `node_spec`、`node_num`、`storage_size`、`band_width`、`enable_public_access`、`cluster_version`、`enable_create_default_ha_mirror_queue` 字段，在 update 函数中支持这些字段的更新。

**理由:**
- 这些字段是实例规格参数，在云资源管理中通常支持动态调整
- 腾讯云的 `ModifyRabbitMQVipInstance` API 应该支持这些字段的修改（需要通过实际测试验证）
- 允许用户在不停机的情况下调整实例规格，提高资源管理效率

**替代方案:**
- 保持现状：用户必须销毁并重建实例才能调整规格，这会导致数据丢失和服务中断
- 仅支持部分字段更新：选择性地支持某些字段，但需要明确判断哪些字段 API 确实支持，增加实现复杂度

### 决策 2: 保持基础架构和付费字段为不可变

**选择:** 保持 `zone_ids`、`vpc_id`、`subnet_id`、`pay_mode`、`time_span`、`auto_renew_flag` 字段在 `immutableArgs` 列表中，更新这些字段时返回错误提示。

**理由:**
- 可用区、VPC、子网是实例的基础架构属性，创建后无法变更
- 付费模式和购买时长属于计费相关属性，通常不允许修改（需通过续费或购买新实例实现）
- 保留这些字段为不可变可以避免 API 调用失败或不一致的状态

**替代方案:**
- 尝试通过 API 更新这些字段：但 API 可能不支持，导致更新失败和用户体验混乱
- 完全移除不可变检查：让 API 返回错误，但用户会收到不清晰的错误消息

### 决策 3: 在 update 函数中逐步构建 API 请求

**选择:** 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，对每个可更新字段使用 `d.HasChange()` 检查变更，并在 `needUpdate` 标志控制下调用 `ModifyRabbitMQVipInstance` API。

**理由:**
- 仅在有字段变更时才调用 API，避免不必要的 API 请求
- 符合现有代码的模式（已用于 `cluster_name` 和 `resource_tags` 字段）
- 清晰地将变更检测、请求构建和 API 调用分离

**替代方案:**
- 无论是否有变更都调用 API：增加 API 调用次数，可能触发不必要的云资源操作
- 为每个字段单独调用 API：增加实现复杂度，且 API 可能不支持单字段更新

### 决策 4: 保持现有的重试和错误处理逻辑

**选择:** 保持使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行 API 调用，并记录详细的日志。

**理由:**
- 符合现有代码模式，保持一致性
- 云 API 调用可能因临时性错误失败，重试机制可以提高稳定性
- 详细的日志有助于问题排查

**替代方案:**
- 移除重试机制：可能导致临时性错误直接失败
- 使用不同的重试策略：增加维护成本，且现有策略已被验证有效

## Risks / Trade-offs

**风险 1: API 可能不支持某些字段的更新**
- **描述:** 某些字段（如 `cluster_version` 或 `enable_create_default_ha_mirror_queue`）可能实际上无法通过 `ModifyRabbitMQVipInstance` API 修改，导致更新失败
- **缓解措施:** 在实现前通过 API 文档或实际测试验证每个字段的可更新性；如果 API 返回不支持的字段，将该字段重新加入不可变列表
- **影响:** 中等 - 如果假设错误，用户更新时会收到 API 错误

**风险 2: 更新操作可能导致实例状态不稳定**
- **描述:** 修改节点规格、数量或存储容量可能触发实例的重启或状态变更，导致在更新过程中实例暂时不可用
- **缓解措施:** 在 API 调用后增加状态检查，确保实例更新完成后再返回；在文档中明确说明更新可能的影响
- **影响:** 高 - 可能影响用户的业务连续性

**风险 3: 向后兼容性问题**
- **描述:** 现有资源的 state 中可能包含某些字段的旧值，如果 API 的行为发生变化，可能导致 state 不一致
- **缓解措施:** 确保更新函数仅在有实际变更时才调用 API；在 read 函数中正确读取所有字段的最新值；不修改 schema 定义
- **影响:** 低 - 仅更新行为不会破坏现有 state

**风险 4: 文档未及时更新导致用户困惑**
- **描述:** 如果文档没有明确说明哪些字段可以更新、哪些需要重建，用户可能会尝试更新不可变字段并收到错误
- **缓解措施:** 在实现后同步更新相关文档；在错误消息中明确提示需要重建的字段
- **影响:** 中等 - 影响用户体验

**权衡: 更灵活 vs 更安全**
- 当前实现过于保守（所有规格字段不可变），导致用户无法灵活管理实例
- 新实现增加了灵活性，但也增加了 API 调用的风险
- 决策优先考虑灵活性，但通过错误处理和文档确保安全性

## Migration Plan

**部署步骤:**
1. 修改 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中的 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
2. 更新 `immutableArgs` 列表，移除可更新字段
3. 为每个可更新字段添加变更检测和 API 请求构建逻辑
4. 运行单元测试和验收测试，验证更新操作的正确性
5. 更新相关文档，明确说明可更新字段和不可更新字段
6. 提交代码并通过 code review
7. 发布新版本的 provider

**回滚策略:**
- 如果发现严重的兼容性问题或 bug，可以通过恢复代码回滚到之前的版本
- 由于不涉及 schema 修改，回滚不会破坏现有 state
- 如果某些字段实际上无法更新，可以将它们重新加入不可变列表并发布修复版本

**测试策略:**
- 使用 `TF_ACC=1` 运行验收测试，覆盖所有可更新字段的更新场景
- 测试边界情况：空值、零值、无效值等
- 测试错误场景：API 失败、网络错误、权限错误等
- 测试并发场景：同时更新多个字段

## Open Questions

1. **API 支持验证:** 需要验证 `ModifyRabbitMQVipInstance` API 是否确实支持所有计划更新的字段（特别是 `cluster_version` 和 `enable_create_default_ha_mirror_queue`），如果不支持应该如何处理？

2. **异步操作处理:** 修改节点规格、数量或存储容量可能触发异步操作，是否需要在 update 函数中等待操作完成？如何检测操作完成？

3. **状态检查机制:** 更新后如何确保实例处于正常运行状态？是否需要增加状态轮询逻辑？

4. **错误消息优化:** 当用户尝试更新不可变字段时，错误消息是否可以更加友好，明确提示需要重建实例？
