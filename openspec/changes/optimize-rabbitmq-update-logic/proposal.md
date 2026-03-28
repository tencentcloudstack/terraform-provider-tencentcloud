## Why

RabbitMQ VIP Instance 资源的 update 函数当前将大量字段标记为不可变，包括 `node_spec`、`node_num`、`storage_size`、`band_width` 等关键规格参数。这些字段实际上可以通过腾讯云 API 的 `ModifyRabbitMQVipInstance` 接口进行修改，但当前实现错误地阻止了这些合法的更新操作，限制了用户的灵活性和实例管理能力。

## What Changes

- 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数，将以下字段从不可变列表中移除，支持动态更新：
  - `node_spec`: 节点规格（2C4G、4C8G 等）
  - `node_num`: 节点数量
  - `storage_size`: 存储容量
  - `band_width`: 公网带宽
  - `enable_public_access`: 公网访问开关
  - `cluster_version`: 集群版本（如果 API 支持）
  - `enable_create_default_ha_mirror_queue`: 镜像队列开关（如果 API 支持）
- 保持以下字段为不可变（需要重建实例）：
  - `zone_ids`: 可用区（基础架构属性）
  - `vpc_id` 和 `subnet_id`: 网络配置
  - `pay_mode`: 付费模式
  - `time_span`: 购买时长
  - `auto_renew_flag`: 自动续费标志
- 更新字段描述，明确哪些字段可更新、哪些需要重建
- 更新相关文档，说明更新行为和限制

## Capabilities

### New Capabilities
- `rabbitmq-instance-dynamic-update`: 支持 RabbitMQ VIP 实例的规格动态更新能力，包括节点规格、节点数量、存储容量、公网带宽和公网访问开关的在线修改

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 修改规格变更行为，从完全不可变（需重建）改为支持部分字段的动态更新，明确区分可更新字段和不可更新字段

## Impact

- 受影响的代码：
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`: 更新函数逻辑和不可变字段列表
  - 可能需要更新 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` 文档
- API 影响：使用 `ModifyRabbitMQVipInstance` API 进行更新调用
- 用户影响：现有资源在 update 时不再报错，可以进行规格调整；升级 provider 后，terraform plan 可能会显示某些字段的 update 操作，实际行为更符合用户预期
- 依赖影响：需要验证 Tencent Cloud API 的实际支持能力，确保字段确实可以通过 Modify API 更新
