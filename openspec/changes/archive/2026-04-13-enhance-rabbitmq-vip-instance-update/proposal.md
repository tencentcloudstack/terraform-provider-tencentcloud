## Why

当前 RabbitMQ VIP 实例资源的 update 逻辑不够完善，没有充分利用腾讯云云 API ModifyRabbitMQVipInstance 支持的所有可更新字段。这导致用户在更新某些实例配置时（如删除保护、集群风险提示、备注等）需要手动通过其他方式操作，无法通过 Terraform 统一管理。为了提升用户体验和配置管理的一致性，需要扩展 update 逻辑以支持更多可更新的字段。

## What Changes

- 新增 `remark` 字段支持，允许用户更新实例备注信息
- 新增 `enable_deletion_protection` 字段支持，允许用户控制实例删除保护
- 新增 `enable_risk_warning` 字段支持，允许用户控制集群风险提示
- 保持现有 `cluster_name` 和 `resource_tags` 的更新能力
- 更新文档以反映新增的可更新字段

## Capabilities

### New Capabilities

None (this change enhances existing capability, no new spec file needed)

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 扩展现有 RabbitMQ VIP 实例资源的 update 能力，添加 remark、enable_deletion_protection、enable_risk_warning 三个可更新字段

## Impact

- 修改文件: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- 新增/修改字段: Schema 中新增三个 Optional 字段，update 函数中添加对应的更新逻辑
- 向后兼容性: 完全兼容，仅新增 Optional 字段，不修改现有字段的行为
- 测试文件: 需要更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 以覆盖新增字段的测试场景
