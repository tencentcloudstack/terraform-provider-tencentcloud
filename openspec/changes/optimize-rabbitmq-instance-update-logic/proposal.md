## Why

当前 RabbitMQ 实例的 update 逻辑存在局限性，将很多本可通过 API 更新的字段（如 node_spec、node_num、storage_size、band_width、enable_public_access 等）标记为不可变，强制用户删除重建资源。这不仅增加了运维成本，还导致数据迁移风险和业务中断。腾讯云 API 实际上支持部分字段的热更新，因此有必要优化 update 逻辑，充分利用 API 能力。

## What Changes

- 修改 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 Update 函数，移除不必要的不可变字段限制
- 为可更新的字段（如 node_spec、node_num、storage_size、band_width、enable_public_access 等）添加适当的 API 调用支持
- 优化错误处理和状态检查逻辑，确保更新操作的安全性和可靠性
- 更新资源文档，明确哪些字段可更新、哪些不可更新

## Capabilities

### New Capabilities
- `tdmq-rabbitmq-vip-instance-update-optimization`: 优化 RabbitMQ VIP 实例的更新能力，支持更多字段的热更新

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 扩展现有实例资源的更新能力，修改部分字段的可更新性约束

## Impact

- 受影响的代码：`/repo/tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- 受影响的 API：`ModifyRabbitMQVipInstance` API 的使用方式
- 受影响的文档：`/repo/website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
- 依赖系统：Tencent Cloud TDMQ API
- 向后兼容性：保持向后兼容，已存在的资源不会受到负面影响
