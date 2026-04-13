## Why

当前 RabbitMQ VIP 实例资源在 Update 操作中存在两个关键问题：
1. 云 API 支持修改的字段（Remark、EnableDeletionProtection、EnableRiskWarning）在 Terraform Provider 中未被实现，导致用户无法通过 Terraform 更新这些配置
2. 虽然云 API 的 ModifyRabbitMQVipInstance 接口支持这些字段的修改，但当前的 Schema 定义和 Update 逻辑未能充分利用这些能力

## What Changes

- **新增 Schema 字段支持**：
  - `remark`: 集群说明信息（可选，Computed）
  - `enable_deletion_protection`: 是否开启删除保护（可选，Computed，默认 false）
  - `enable_risk_warning`: 是否开启集群风险提示（可选，Computed，默认 false）

- **优化 Update 逻辑**：
  - 支持通过 ModifyRabbitMQVipInstance 接口更新 remark、enable_deletion_protection、enable_risk_warning 字段
  - 完善 resource_tags 的更新逻辑，正确处理 RemoveAllTags 标志

- **优化 Read 逻辑**：
  - 从 DescribeRabbitMQVipInstance 和 DescribeRabbitMQVipInstances 响应中读取新增字段并设置到 state

## Capabilities

### New Capabilities
- `tdmq-rabbitmq-vip-instance-update-fields`: 为 tencentcloud_tdmq_rabbitmq_vip_instance 资源添加 Remark、EnableDeletionProtection、EnableRiskWarning 字段的读取和更新支持

### Modified Capabilities
- 无需求级别的变更，仅实现细节优化

## Impact

- 受影响的代码文件：
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - Schema 定义和 CRUD 函数
- 受影响的 API：
  - 使用 ModifyRabbitMQVipInstance API 更新新增字段
  - 从 DescribeRabbitMQVipInstance 和 DescribeRabbitMQVipInstances API 响应中读取新增字段
- 向后兼容性：
  - 新增字段均为 Optional 和 Computed，不会破坏现有的 Terraform 配置和 state
  - 不修改现有字段的定义或行为
