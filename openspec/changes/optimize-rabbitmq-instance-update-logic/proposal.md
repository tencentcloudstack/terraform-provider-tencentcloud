## Why

RabbitMQ 实例资源的 update 逻辑当前只支持更新 `cluster_name` 和 `resource_tags` 两个参数，但腾讯云 API (`ModifyRabbitMQVipInstance`) 实际支持更多的可更新参数，包括 `remark`（实例备注）、`enable_deletion_protection`（删除保护）、`enable_risk_warning`（风险提示）。这导致用户无法通过 Terraform 更新这些可配置的实例属性，限制了资源管理的灵活性。

## What Changes

- 在 RabbitMQ VIP 实例资源中新增 `remark` 参数，支持更新实例备注信息
- 在 RabbitMQ VIP 实例资源中新增 `enable_deletion_protection` 参数，支持开启/关闭删除保护功能
- 在 RabbitMQ VIP 实例资源中新增 `enable_risk_warning` 参数，支持开启/关闭集群风险提示
- 更新 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数，支持上述三个新参数的更新
- 更新 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数，读取这三个新参数的值并设置到资源状态中

## Capabilities

### New Capabilities
- `update-rabbitmq-instance-extra-parameters`: Add support for updating RabbitMQ instance parameters that are supported by cloud API but currently not available in the provider, including remark, enable_deletion_protection, and enable_risk_warning.

### Modified Capabilities
(Leave empty if no requirement changes)

## Impact

- Affected code: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- Affected tests: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- API dependencies: `ModifyRabbitMQVipInstance` API (already in use)
- Backward compatibility: All changes are additive (new optional parameters), no breaking changes to existing configurations
