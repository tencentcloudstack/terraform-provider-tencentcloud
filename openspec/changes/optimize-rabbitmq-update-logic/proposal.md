## Why

当前 RabbitMQ VIP 实例资源的 update 逻辑不完整。腾讯云 TDMQ API 的 `ModifyRabbitMQVipInstance` 接口支持多个字段的修改（包括 `remark`、`enable_deletion_protection`、`enable_risk_warning`），但当前的 Terraform Provider 实现仅支持修改 `cluster_name` 和 `resource_tags` 字段。

这导致用户无法通过 Terraform 声明式地管理 RabbitMQ 实例的以下重要配置：
- **remark**: 实例备注信息，用于资源标识和管理
- **enable_deletion_protection**: 删除保护开关，防止误删除关键实例
- **enable_risk_warning**: 集群风险提示开关，增强运维安全性

同时，当前 update 逻辑在 `immutableArgs` 列表中包含了 `auto_renew_flag` 等参数，但这些参数实际上不应该在 update 时被严格禁止（因为它们可以通过其他 API 修改）。

通过完善 RabbitMQ 实例的 update 逻辑，可以：
- 提供更完整的资源管理能力
- 支持更多实例配置的声明式管理
- 提高运维安全性和管理效率
- 与腾讯云 API 的能力保持一致

## What Changes

优化 RabbitMQ VIP 实例资源的 update 逻辑，完善字段支持和错误处理：

### 新增 Schema 字段

在 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` 中添加以下字段：

1. **`remark`** (Computed, Optional)
   - 实例备注信息
   - 说明：用于资源标识和管理
   - 类型：String

2. **`enable_deletion_protection`** (Computed, Optional)
   - 删除保护开关
   - 说明：防止误删除关键实例，默认为 false
   - 类型：Bool

3. **`enable_risk_warning`** (Computed, Optional)
   - 集群风险提示开关
   - 说明：开启集群风险提示，增强运维安全性
   - 类型：Bool

### Update 逻辑优化

1. 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中添加对新字段的支持：
   - 检测 `remark` 字段变化并调用 API 修改
   - 检测 `enable_deletion_protection` 字段变化并调用 API 修改
   - 检测 `enable_risk_warning` 字段变化并调用 API 修改

2. 优化 `immutableArgs` 列表，移除不应在 update 阶段禁止的字段（如 `auto_renew_flag`），因为这些字段可以通过其他 API 单独修改

### Read 逻辑优化

在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中添加对新字段的读取：
- 从 API 响应中读取 `remark` 字段
- 从 API 响应中读取 `enable_deletion_protection` 字段
- 从 API 响应中读取 `enable_risk_warning` 字段

### 修改文件清单

- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - 添加 3 个新 Schema 字段
  - 更新 Update 函数，支持新字段的修改
  - 更新 Read 函数，读取新字段值
  - 优化 immutableArgs 列表

### Schema 更新示例

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                         = [1]
  vpc_id                           = "vpc-xxxxxx"
  subnet_id                        = "subnet-xxxxxx"
  cluster_name                     = "rabbitmq-cluster"
  node_spec                        = "rabbit-vip-basic-1"
  node_num                         = 1
  storage_size                     = 200

  # 新增字段
  remark                           = "Production RabbitMQ instance"
  enable_deletion_protection       = true
  enable_risk_warning              = true
}
```

## Capabilities

### New Capabilities

无新增功能模块，这是对现有 RabbitMQ 实例资源的功能增强。

### Modified Capabilities

- **tdmq-rabbitmq-vip-instance**: 增强 RabbitMQ VIP 实例资源的 update 能力，支持备注、删除保护和风险提示的声明式管理

## Impact

### 受影响的规范

- 修改现有规范：`tdmq-rabbitmq-vip-instance` - 在现有能力基础上增强字段支持

### 受影响的代码

- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 唯一需要修改的文件

### 向后兼容性

- ✅ 完全向后兼容
- ✅ 新增字段均为 Optional 和 Computed，不影响现有配置
- ✅ 不修改现有资源的 Schema 结构
- ✅ 不影响现有 Terraform 状态文件

### API 依赖

- 依赖腾讯云 TDMQ API `ModifyRabbitMQVipInstance` 的完整字段支持
- 依赖腾讯云 TDMQ API `DescribeRabbitMQVipInstance` 返回新增字段的值

### 测试影响

- 需要更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中的测试用例
- 需要添加对新字段的 update 测试
- 需要添加对 read 新字段的测试
