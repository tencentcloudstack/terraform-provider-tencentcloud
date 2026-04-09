# Change: 优化 RabbitMQ 实例的 update 逻辑

## Why

当前 RabbitMQ 实例资源的 update 逻辑存在以下问题：

1. **缺少部分可更新字段的支持**：腾讯云 API 的 `ModifyRabbitMQVipInstanceRequest` 接口支持修改多个字段，但当前 Terraform Provider 实现只支持 `cluster_name` 和 `resource_tags` 的更新。

2. **缺少关键字段支持**：以下 API 支持的字段在当前实现中缺失：
   - `remark`：备注信息
   - `enable_deletion_protection`：删除保护开关
   - `enable_risk_warning`：集群风险提示开关

3. **限制用户灵活性**：由于缺少这些字段的支持，用户无法通过 Terraform 完整管理 RabbitMQ 实例的所有可更新属性，需要手动通过腾讯云控制台或 API 进行操作，破坏了基础设施即代码的完整性。

## What Changes

### 新增字段支持

1. **remark**（可选，字符串类型）
   - 描述：实例备注信息
   - API 字段：`Remark`
   - 支持更新：是

2. **enable_deletion_protection**（可选，布尔类型）
   - 描述：是否开启删除保护
   - API 字段：`EnableDeletionProtection`
   - 支持更新：是
   - 默认值：false

3. **enable_risk_warning**（可选，布尔类型）
   - 描述：是否开启集群风险提示
   - API 字段：`EnableRiskWarning`
   - 支持更新：是
   - 默认值：false

### Schema 变更

在 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` 的 Schema 中新增以下字段：

```go
"remark": {
    Optional:    true,
    Type:        schema.TypeString,
    Description: "Instance remark information.",
},

"enable_deletion_protection": {
    Optional:    true,
    Type:        schema.TypeBool,
    Description: "Whether to enable deletion protection. Default is false.",
},

"enable_risk_warning": {
    Optional:    true,
    Type:        schema.TypeBool,
    Description: "Whether to enable cluster risk warning. Default is false.",
},
```

### Update 逻辑变更

在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中：
1. 从 `immutableArgs` 列表中移除 `remark`、`enable_deletion_protection`、`enable_risk_warning`（如果存在）
2. 添加这三个字段的更新逻辑，检查 `d.HasChange()` 并在变更时调用相应的 API

### Read 逻辑变更

在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中：
1. 从 API 响应中读取 `remark`、`enable_deletion_protection`、`enable_risk_warning` 字段
2. 将这些字段的值设置到 resource data 中

## Capabilities

### Modified Capabilities

- `tdmq-rabbitmq-vip-instance`: 扩展了 RabbitMQ VIP 实例的可更新属性范围，新增备注、删除保护和风险提示三个字段的更新支持。

## Impact

### 受影响的规范

- 修改规范：`tdmq-rabbitmq-vip-instance` - RabbitMQ VIP 实例管理

### 受影响的代码

- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 资源实现文件
  - 修改 `ResourceTencentCloudTdmqRabbitmqVipInstance` 函数，新增三个字段的 Schema 定义
  - 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数，支持创建时设置这三个字段
  - 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数，支持读取这三个字段
  - 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数，支持更新这三个字段

### 向后兼容性

- ✅ 完全向后兼容，所有变更都是新增字段，不影响现有配置
- ✅ 新增字段均为 Optional，不会破坏现有的 Terraform 配置
- ✅ 不修改现有字段的定义和行为

### 依赖关系

- 无新增外部依赖
- 依赖现有的腾讯云 TDMQ SDK：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`

### 测试影响

- 需要更新验收测试用例 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`
- 新增测试场景：验证新增字段的创建、读取和更新功能
- 测试需要确保新字段能够正确地在创建时设置、更新时修改、读取时获取
