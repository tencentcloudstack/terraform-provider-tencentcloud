## Why

RabbitMQ VIP 实例的 Update 逻辑存在以下问题：
1. `ModifyRabbitMQVipInstance` API 支持的多个字段未在 update 逻辑中使用（如 `EnableDeletionProtection`、`EnableRiskWarning`、`Remark`）
2. 字段被错误地标记为不可修改（immutable），限制了用户更新实例配置的能力
3. 缺少对变更操作的等待和状态检查机制，可能导致状态不一致
4. 错误处理和重试逻辑不够完善

本次变更旨在优化 Update 逻辑，充分利用 API 能力，提供更好的用户体验和状态一致性。

## What Changes

- **新增 Update 支持字段**：支持更新 `enable_deletion_protection`、`enable_risk_warning`、`remark` 字段
- **修正不可修改字段列表**：移除 API 实际支持修改但被错误标记为不可修改的字段
- **增强 Update 状态等待**：在 Update 操作后添加状态等待机制，确保实例状态一致
- **优化错误处理**：改进错误消息，提供更清晰的变更失败原因
- **完善标签更新逻辑**：优化 `resource_tags` 的增量更新处理

## Capabilities

### New Capabilities

None (enhancing existing capability)

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 扩展 Update 操作的支持字段，优化更新逻辑

## Impact

- **影响代码**：`tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- **影响测试**：需要更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 的 update 测试用例
- **向后兼容性**：保持完全向后兼容，仅新增功能和优化现有逻辑
- **API 调用**：增强 `ModifyRabbitMQVipInstance` API 的使用，添加 `DescribeRabbitMQVipInstance` 的状态检查
