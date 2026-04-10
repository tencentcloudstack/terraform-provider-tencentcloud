## Why

当前 RabbitMQ 实例的 update 逻辑存在两个主要问题：一是腾讯云 API 已支持部分新增参数（如 Remark、EnableDeletionProtection、EnableRiskWarning），但 Terraform Provider 未实现对这些参数的更新支持；二是现有 update 逻辑在代码结构上缺乏清晰的可扩展性，参数处理逻辑分散且难以维护。这导致用户无法通过 Terraform 管理这些新增的可更新属性，同时也增加了未来扩展更新功能的维护成本。

## What Changes

- **新增支持更新的参数**：
  - `remark`：实例备注信息
  - `enable_deletion_protection`：是否开启删除保护
  - `enable_risk_warning`：是否开启集群风险提示

- **代码结构优化**：
  - 重构 update 函数，将参数处理逻辑模块化
  - 添加统一的参数变更检测机制
  - 改进错误处理和日志记录
  - 优化标签更新逻辑的代码结构

- **Schema 更新**：
  - 为新增的可更新参数添加 schema 定义
  - 更新现有参数的描述文档

## Capabilities

### New Capabilities

- `tdmq-rabbitmq-instance-update-enhancement`: 增强 RabbitMQ 实例的更新能力，支持通过 Terraform 管理实例的备注信息、删除保护开关和集群风险提示开关。

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 更新现有 RabbitMQ 实例资源的可更新参数列表，移除部分不必要地标记为不可变的参数（基于实际 API 支持），并新增对 remark、enable_deletion_protection、enable_risk_warning 参数的更新支持。

## Impact

### 受影响的代码

- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`：主要修改文件，包括 schema 定义和 update 函数的重构
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`：需要添加对新参数的单元测试

### 受影响的 API

- `ModifyRabbitMQVipInstance`：利用此 API 的所有参数进行实例更新

### 向后兼容性

- ✅ 完全向后兼容，不破坏现有配置
- ✅ 新增参数均为 Optional 字段，不影响已有资源
- ✅ 不可变参数列表的调整不会影响现有用户（因为 API 本身就不支持这些参数的更新）

### 依赖关系

- 依赖 `tencentcloud-sdk-go` TDMQ 模块已支持新参数（已验证当前版本支持）
