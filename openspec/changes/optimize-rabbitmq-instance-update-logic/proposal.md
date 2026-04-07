## Why

RabbitMQ VIP 实例的 update 逻辑存在过度限制，将 node_spec、node_num、storage_size、band_width 等关键字段标记为不可修改。这导致用户无法在实例创建后进行规格升级、节点扩容、存储扩容等常见运维操作，降低了 Terraform 资源管理的灵活性。

## What Changes

- 移除或修改 immutableArgs 列表中可由腾讯云 API 支持修改的字段
- 支持通过 update 操作修改实例规格（node_spec）
- 支持通过 update 操作修改节点数量（node_num）
- 支持通过 update 操作修改存储大小（storage_size）
- 支持通过 update 操作修改公网带宽（band_width）
- 支持通过 update 操作开启/关闭公网访问（enable_public_access）
- 增强错误处理，提供更友好的错误提示

## Capabilities

### New Capabilities
- `tdmq-rabbitmq-instance-update`: 优化 RabbitMQ VIP 实例的 update 逻辑，支持更多字段的修改操作

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 扩展现有资源实例的 update 能力，允许修改更多配置参数

## Impact

- **Affected Code**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- **Affected Tests**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- **API Dependencies**: TencentCloud TDMQ ModifyRabbitMQVipInstance API 及相关规格变更 API
- **Backward Compatibility**: 需要确保现有用户配置不受影响，新增的 update 操作应该是可选的
