## Why

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑过于严格，将 11 个参数标记为不可修改参数（immutableArgs），但实际上腾讯云 API 支持部分参数的修改。这导致用户无法通过 Terraform 灵活地修改实例配置，必须手动通过控制台或 API 进行修改，破坏了基础设施即代码的一致性。

## What Changes

- 将 `auto_renew_flag` 参数从 `immutableArgs` 列表中移除，支持通过 `ModifyRabbitMQVipInstance` API 修改自动续费标识
- 将 `band_width` 参数从 `immutableArgs` 列表中移除，支持公网带宽的动态调整
- 将 `enable_public_access` 参数从 `immutableArgs` 列表中移除，支持公网访问的开启/关闭操作
- 添加异步操作等待机制，对于需要较长时间的更新操作（如带宽修改、公网访问开关），支持等待操作完成
- 完善 update 函数的错误处理和日志记录

## Capabilities

### New Capabilities

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 支持更多参数的在线修改，包括自动续费标识、公网带宽、公网访问开关

## Impact

- **受影响的资源**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
- **受影响的测试**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- **受影响的文档**: `website/docs/r/tdmq_rabbitmq_vip_instance.md` (如果存在)
- **API 调用**: 扩展 `ModifyRabbitMQVipInstance` API 的参数支持范围
- **向后兼容**: 本次变更仅增加可修改参数的范围，不破坏现有配置，保持完全向后兼容
