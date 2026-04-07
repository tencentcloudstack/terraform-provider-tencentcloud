# Change: 优化 RabbitMQ 实例的 update 逻辑

## Why

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑过于严格，将许多本应可更新的参数标记为不可变，导致用户在基础设施变更时无法灵活调整实例配置。

具体问题包括：
1. **过度限制的不可变参数**：`zone_ids`、`vpc_id`、`subnet_id`、`node_spec`、`node_num`、`storage_size`、`auto_renew_flag`、`time_span`、`pay_mode`、`cluster_version`、`band_width`、`enable_public_access`、`enable_deletion_protection` 等参数被禁止更新
2. **缺少异步操作支持**：update 操作没有等待资源状态更新的逻辑，可能导致后续操作失败
3. **缺少新参数支持**：Modify API 新支持的参数（如 `remark`、`enable_risk_warning`）未被集成到 update 逻辑中
4. **不符合用户实际需求**：实际生产环境中，用户需要能够调整节点规格、节点数量、存储大小、带宽、删除保护等配置以适应业务变化

腾讯云 API `ModifyRabbitMQVipInstance` 支持更新多个参数，但当前实现仅支持更新 `cluster_name` 和 `resource_tags`，限制了用户的运维灵活性。

## What Changes

优化 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑，移除不必要的限制，支持更多参数的更新：

### 修改的文件
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 优化 update 函数逻辑

### 支持更新的参数
根据腾讯云 API 支持情况和实际需求，以下参数应该支持更新：

**已支持（保持不变）：**
- `cluster_name` - 集群名称
- `resource_tags` - 资源标签

**新增支持（从不可变列表中移除）：**
- `auto_renew_flag` - 自动续费标识
- `enable_public_access` - 是否启用公网访问
- `band_width` - 公网带宽
- `enable_deletion_protection` - 是否启用删除保护
- `remark` - 实例备注
- `enable_risk_warning` - 是否启用集群风险提示

### 仍保持不可变的参数
以下参数由于涉及底层架构变更，保持不可变（这些参数的修改需要重建实例）：
- `zone_ids` - 可用区（涉及物理部署位置）
- `vpc_id` - VPC ID（涉及网络架构）
- `subnet_id` - 子网 ID（涉及网络架构）
- `node_spec` - 节点规格（可能需要实例重建）
- `node_num` - 节点数量（可能需要实例重建）
- `storage_size` - 存储大小（可能需要实例重建）
- `enable_create_default_ha_mirror_queue` - 镜像队列配置（创建时决定）
- `time_span` - 购买时长（预付费实例的特殊限制）
- `pay_mode` - 付费模式（预付费/后付费的转换限制）
- `cluster_version` - 集群版本（版本升级需要特殊流程）

### 新增异步状态等待
update 操作完成后，增加状态等待逻辑，确保资源处于稳定状态后再返回：
- 等待实例状态从 "Updating" 变为 "Running" 或 "Success"
- 支持自定义超时配置
- 提供清晰的错误信息

### 错误处理优化
- 区分不可变参数的错误信息，明确告知用户哪些参数不能修改
- 对于已弃用的 API 参数，提供迁移建议

## Capabilities

### New Capabilities
- `rabbitmq-vip-instance-update-enhancement`: 增强 RabbitMQ VIP 实例的更新能力，支持修改删除保护、备注和风险提示等配置

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: 优化 update 逻辑，支持更多参数的更新和异步状态等待
- `rabbitmq-vip-instance`: 优化 RabbitMQ VIP 实例资源的更新逻辑，扩展可修改参数范围

## Impact

### 受影响的代码
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数

### 受影响的规范
- `tdmq-rabbitmq-vip-instance` - 需要更新规范文档，明确哪些参数支持 update

### 向后兼容性
- ✅ 完全向后兼容
- ✅ 只是移除了不必要的限制，不会破坏现有配置
- ✅ 已有的配置可以继续正常工作
- ✅ 新增可修改参数，不影响已有配置
- ✅ 保持现有不可修改参数的约束（zone_ids、vpc_id、subnet_id 等基础设施参数）

### 依赖关系
- 依赖腾讯云 TDMQ RabbitMQ API `ModifyRabbitMQVipInstance`
- 需要验证 API 支持的参数列表

### 测试影响
- 需要更新验收测试用例，验证新增的 update 功能
- 需要测试异步状态等待逻辑
- 需要验证不可变参数的错误处理
- 需要更新验收测试，覆盖新增的可修改参数
- 需要测试修改 enable_deletion_protection、remark、enable_risk_warning 参数的场景
