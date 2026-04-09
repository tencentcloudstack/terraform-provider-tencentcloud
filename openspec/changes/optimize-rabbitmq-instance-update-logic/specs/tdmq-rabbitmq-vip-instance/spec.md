# Spec: tdmq-rabbitmq-vip-instance (Modified)

## ADDED Requirements

### Requirement: RabbitMQ 实例支持备注字段
The RabbitMQ instance resource SHALL support a `remark` field to allow users to configure instance remarks via Terraform, and SHALL allow updating this remark when needed.

#### Scenario: 创建实例时设置备注
- **WHEN** 用户在创建 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源时设置了 `remark` 字段
- **THEN** 创建的 RabbitMQ 实例应该包含指定的备注信息
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回设置的备注信息

#### Scenario: 更新实例备注
- **WHEN** 用户修改现有 RabbitMQ 实例的 `remark` 字段值
- **THEN** RabbitMQ 实例的备注信息应该更新为新值
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回更新后的备注信息

#### Scenario: 删除备注
- **WHEN** 用户将 `remark` 字段从配置中移除
- **THEN** RabbitMQ 实例的备注信息应该保持不变（Terraform 不管理该字段）
- **THEN** `terraform apply` 成功完成，无错误

### Requirement: RabbitMQ 实例支持删除保护字段
The RabbitMQ instance resource SHALL support an `enable_deletion_protection` field to allow users to configure deletion protection via Terraform, and SHALL allow toggling this protection state when needed.

#### Scenario: 创建实例时启用删除保护
- **WHEN** 用户在创建 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源时设置 `enable_deletion_protection = true`
- **THEN** 创建的 RabbitMQ 实例应该启用删除保护
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回 `enable_deletion_protection = true`

#### Scenario: 创建实例时禁用删除保护
- **WHEN** 用户在创建 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源时设置 `enable_deletion_protection = false` 或不设置该字段
- **THEN** 创建的 RabbitMQ 实例应该禁用删除保护
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回 `enable_deletion_protection = false`

#### Scenario: 切换删除保护状态
- **WHEN** 用户修改现有 RabbitMQ 实例的 `enable_deletion_protection` 字段值（从 false 改为 true，或从 true 改为 false）
- **THEN** RabbitMQ 实例的删除保护状态应该更新为新值
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回更新后的删除保护状态

#### Scenario: 删除保护时尝试删除实例
- **WHEN** RabbitMQ 实例启用了删除保护（`enable_deletion_protection = true`）
- **WHEN** 用户执行 `terraform destroy` 尝试删除实例
- **THEN** 删除操作应该失败，并返回相应的错误信息
- **THEN** Terraform 应该明确告知用户实例受删除保护

### Requirement: RabbitMQ 实例支持集群风险提示字段
The RabbitMQ instance resource SHALL support an `enable_risk_warning` field to allow users to configure cluster risk warning via Terraform, and SHALL allow toggling this warning state when needed.

#### Scenario: 创建实例时启用风险提示
- **WHEN** 用户在创建 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源时设置 `enable_risk_warning = true`
- **THEN** 创建的 RabbitMQ 实例应该启用集群风险提示
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回 `enable_risk_warning = true`

#### Scenario: 创建实例时禁用风险提示
- **WHEN** 用户在创建 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源时设置 `enable_risk_warning = false` 或不设置该字段
- **THEN** 创建的 RabbitMQ 实例应该禁用集群风险提示
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回 `enable_risk_warning = false`

#### Scenario: 切换风险提示状态
- **WHEN** 用户修改现有 RabbitMQ 实例的 `enable_risk_warning` 字段值（从 false 改为 true，或从 true 改为 false）
- **THEN** RabbitMQ 实例的集群风险提示状态应该更新为新值
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回更新后的集群风险提示状态

## MODIFIED Requirements

### Requirement: 现有字段保持不变
The existing `cluster_name` and `resource_tags` fields SHALL maintain their current update logic and behavior.

#### Scenario: cluster_name 更新
- **WHEN** 用户修改现有 RabbitMQ 实例的 `cluster_name` 字段值
- **THEN** RabbitMQ 实例的集群名称应该更新为新值
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回更新后的集群名称

#### Scenario: resource_tags 更新
- **WHEN** 用户修改现有 RabbitMQ 实例的 `resource_tags` 字段
- **THEN** RabbitMQ 实例的标签应该更新为新值
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回更新后的标签

#### Scenario: resource_tags 清除
- **WHEN** 用户将 `resource_tags` 设置为空列表
- **THEN** RabbitMQ 实例的所有标签应该被删除
- **THEN** `terraform apply` 成功完成，无错误
- **THEN** 读取实例时应该返回空的标签列表
