## Why

当前 RabbitMQ VIP 实例资源 `tencentcloud_tdmq_rabbitmq_vip_instance` 的 update 逻辑存在严重限制，大量关键字段被标记为不可修改，包括节点规格（`node_spec`）、节点数量（`node_num`）、存储大小（`storage_size`）、带宽（`band_width`）和公网访问开关（`enable_public_access`）等核心配置参数。这些参数在实际使用中是需要动态调整的，当前的限制导致用户无法通过 Terraform 灵活地调整实例配置，降低了资源管理的效率和灵活性。

## What Changes

- 优化 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数的实现逻辑
- 移除或放宽不必要的参数限制，使核心配置参数支持 update 操作
- 为无法通过 ModifyRabbitMQVipInstance API 修改的参数，探索替代的更新方式（如通过其他 API 或重建实例的策略）
- 改进错误处理和重试逻辑，确保 update 操作的稳定性
- 增强对 update 操作的日志记录，便于问题排查和审计

## Capabilities

### New Capabilities
- `rabbitmq-instance-dynamic-update`: 支持 RabbitMQ VIP 实例核心配置参数的动态更新能力，包括节点规格、节点数量、存储大小、带宽和公网访问等参数的修改支持

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 扩展现有资源的 update 能力，允许修改更多实例配置参数

## Impact

**受影响的代码**：
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 主要修改文件，更新 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` - 需要更新测试用例以覆盖新的 update 场景

**受影响的 API**：
- `ModifyRabbitMQVipInstance` API - 可能需要探索其他 API 来实现参数更新
- 可能涉及其他 TDMQ 相关 API，如修改实例规格、网络配置等

**依赖和系统**：
- 依赖于 TDMQ API 的发展和更新
- 需要考虑向后兼容性，确保不会破坏现有用户配置
- 可能需要更新 Terraform provider 文档和示例
