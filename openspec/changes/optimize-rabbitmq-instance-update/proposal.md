## Why

RabbitMQ VIP 实例资源的 update 逻辑不完整，云 API 支持的多个可更新参数未被利用。当前 update 函数仅支持 cluster_name 和 resource_tags 的更新，缺少对 remark、enable_deletion_protection 和 enable_risk_warning 等字段的 CRUD 支持，限制了用户通过 Terraform 管理这些属性的能力。

## What Changes

- 在 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 Schema 中新增以下字段：
  - `remark`: 实例备注（可选）
  - `enable_deletion_protection`: 是否开启删除保护（可选，Computed）
  - `enable_risk_warning`: 是否开启集群风险提示（可选，Computed）
- 在 Create 函数中支持这些字段的初始化
- 在 Read 函数中从云 API 响应中读取这些字段的值
- 在 Update 函数中支持对这些字段的更新操作，调用 ModifyRabbitMQVipInstance API

## Capabilities

### New Capabilities

- `rabbitmq-instance-update`: 扩展 RabbitMQ VIP 实例资源的 CRUD 能力，支持 remark、enable_deletion_protection 和 enable_risk_warning 字段的完整生命周期管理。

### Modified Capabilities

None - This is an enhancement to existing capabilities without changing their core requirements.

## Impact

- Affected code: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- Affected tests: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- Used APIs:
  - Create: CreateRabbitMQVipInstance
  - Read: DescribeRabbitMQVipInstance, DescribeRabbitMQVipInstances
  - Update: ModifyRabbitMQVipInstance
  - Delete: DeleteRabbitMQVipInstance
- Documentation: Needs update to `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
