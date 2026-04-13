## Why

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 方法存在两个问题：
1. **功能不完整**：ModifyRabbitMQVipInstance API 支持的 `remark`、`enable_risk_warning` 等字段未在资源中实现，导致用户无法通过 Terraform 管理这些属性
2. **不可变字段过多**：当前将 `enable_deletion_protection` 等字段标记为不可变，但云 API 实际上支持这些字段的修改，限制了用户对实例的灵活管理

## What Changes

- 新增 `remark` 字段到 Schema，支持修改实例备注
- 新增 `enable_risk_warning` 字段到 Schema，支持开启/关闭集群风险提示
- 将 `enable_deletion_protection` 字段从不可变列表中移除，支持修改删除保护状态
- 优化 update 方法的不可变字段列表，使其与云 API 的实际能力保持一致

## Capabilities

### New Capabilities

无新增能力，仅扩展现有 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 功能

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 扩展 update 能力，新增 `remark`、`enable_risk_warning` 字段支持，并允许修改 `enable_deletion_protection` 字段

## Impact

### 受影响的文件
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 修改 Schema 和 Update 方法
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` - 扩展测试用例
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` - 更新资源文档

### 向后兼容性
- ✅ 完全向后兼容，仅新增 Optional 字段
- ✅ 不修改现有字段的行为和类型

### 云 API 依赖
- 依赖 `ModifyRabbitMQVipInstance` API，支持修改实例备注、删除保护和风险提示
