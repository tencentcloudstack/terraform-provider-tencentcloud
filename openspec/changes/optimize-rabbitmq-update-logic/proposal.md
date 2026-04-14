## Why

当前 RabbitMQ VIP 实例的 update 逻辑存在以下问题：
1. 云 API `ModifyRabbitMQVipInstance` 支持的参数未完全实现，包括 `remark`（备注）、`enable_deletion_protection`（是否开启删除保护）、`enable_risk_warning`（是否开启集群风险提示）等参数
2. 用户无法通过 Terraform 管理实例的重要配置（如删除保护、备注等）
3. 不完整的 update 能力限制了基础设施即代码（IaC）的完整性

通过优化 update 逻辑，可以实现更完整的实例配置管理，提高 Terraform Provider 的功能完整性。

## What Changes

- 在 `resource_tc_tdmq_rabbitmq_vip_instance` schema 中新增以下可更新参数：
  - `remark`: 备注信息（可选，字符串）
  - `enable_deletion_protection`: 是否开启删除保护（可选，布尔值）
  - `enable_risk_warning`: 是否开启集群风险提示（可选，布尔值）
- 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中实现对新参数的更新逻辑
- 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中实现对新参数的读取逻辑
- 保持现有不可变参数的验证逻辑不变，确保向后兼容

## Capabilities

### New Capabilities

- `rabbitmq-instance-update-enhancement`: 增强 RabbitMQ VIP 实例的更新能力，支持备注、删除保护、风险提示等参数的更新

### Modified Capabilities

- 无

## Impact

- 受影响的代码：`tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- 影响的云 API：`ModifyRabbitMQVipInstance`（腾讯云 TDMQ 服务）
- 依赖：无新增外部依赖
- 兼容性：完全向后兼容，仅新增可选参数，不影响现有配置
