## Why

当前 RabbitMQ VIP Instance 资源的 update 函数仅支持更新 `cluster_name` 和 `resource_tags` 两个参数，大量本应支持更新的参数（如 `node_spec`、`node_num`、`storage_size`、`band_width` 等）被错误地标记为不可变，导致用户无法通过 Terraform 灵活调整实例配置，降低了用户体验和资源管理的便利性。

## What Changes

- 移除 RabbitMQ VIP Instance update 函数中不必要的不可变参数限制
- 新增对以下参数的更新支持：
  - `node_spec`: 节点规格升级
  - `node_num`: 节点数量扩容/缩容
  - `storage_size`: 存储大小调整
  - `band_width`: 公网带宽调整
  - `auto_renew_flag`: 自动续费开关
  - `enable_public_access`: 公网访问开关
- 保持向后兼容，不破坏现有 Terraform 配置和 state

## Capabilities

### New Capabilities

（无新能力）

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 修改现有资源的行为，扩展 update 函数支持的参数范围

## Impact

- **修改文件**:
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`: 主要变更，修改 update 函数逻辑
- **可能影响**:
  - `tencentcloud/services/trabbit/service_tencentcloud_tdmq.go`: 可能需要新增或修改服务层函数以支持参数更新
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`: 可能需要更新测试用例
- **依赖影响**: 需要确认腾讯云 TDMQ API 的 `ModifyRabbitMQVipInstance` 接口是否支持上述参数的修改
