## Why

当前 tencentcloud_tdmq_rabbitmq_vip_instance 资源的 update 逻辑过于严格，将许多可以通过 ModifyRabbitMQVipInstance API 修改的参数标记为不可变，导致用户无法通过 Terraform 更新这些参数。腾讯云 TDMQ RabbitMQ API 实际上支持更多参数的更新，包括 remark（备注）、enable_deletion_protection（删除保护）和 enable_risk_warning（风险提示）。通过优化 update 逻辑，可以充分利用云 API 的能力，提供更好的用户体验。

## What Changes

- 在 resource_tc_tdmq_rabbitmq_vip_instance schema 中新增以下可选参数：
  - `remark`：实例备注信息
  - `enable_deletion_protection`：是否开启删除保护
  - `enable_risk_warning`：是否开启集群风险提示
- 在 Create 方法中支持这些新参数
- 在 Read 方法中读取这些参数
- 在 Update 方法中移除对这些参数的不可变限制，并实现它们的更新逻辑
- 保持向后兼容性，所有新增参数都是可选的
- 更新相关单元测试以覆盖新增功能

## Capabilities

### New Capabilities

此变更不引入新的独立功能，而是扩展现有 RabbitMQ VIP 实例资源的更新能力。

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 扩展现有资源的 update 能力，支持 remark、enable_deletion_protection、enable_risk_warning 参数的更新。这不是需求级别的变更，而是实现细节的优化，使得资源能够正确映射云 API 的更新能力。

## Impact

- 受影响代码：`tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- 受影响测试：`tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- 受影响文档：`website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
- 依赖：tencentcloud-sdk-go 中的 TDMQ API v20200217（已包含必要的 API 支持）
- 系统影响：无，纯功能增强，不破坏现有功能
