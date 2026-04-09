# Design: 优化 RabbitMQ 实例的 update 逻辑

## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑存在限制，只支持更新 `cluster_name` 和 `resource_tags` 两个字段。然而，腾讯云的 `ModifyRabbitMQVipInstance` API 实际上支持更多字段的更新，包括：
- `Remark`（备注信息）
- `EnableDeletionProtection`（删除保护）
- `EnableRiskWarning`（集群风险提示）

这种限制导致用户无法通过 Terraform 完整管理 RabbitMQ 实例的所有可更新属性，降低了基础设施即代码的完整性。

## Goals / Non-Goals

**Goals:**
1. 为 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源新增 `remark` 字段支持
2. 为 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源新增 `enable_deletion_protection` 字段支持
3. 为 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源新增 `enable_risk_warning` 字段支持
4. 实现这三个字段在创建、读取和更新操作中的完整支持
5. 确保向后兼容性，不影响现有配置
6. 更新测试用例，覆盖新增字段的所有操作场景

**Non-Goals:**
1. 不修改现有字段的行为和限制（如 zone_ids、vpc_id 等不可变字段）
2. 不修改资源的删除逻辑
3. 不添加任何非腾讯云 API 支持的字段

## Decisions

### 1. 字段类型和默认值

**决策**：根据腾讯云 API 的定义和 Terraform 最佳实践，确定字段类型：
- `remark`：`schema.TypeString`，无默认值
- `enable_deletion_protection`：`schema.TypeBool`，默认值为 `false`
- `enable_risk_warning`：`schema.TypeBool`，默认值为 `false`

**理由**：
- 备注信息是文本类型，布尔类型用于开关字段符合 Terraform 的常见实践
- 根据 API 文档，删除保护和风险提示的默认值均为 false

### 2. 字段可选性

**决策**：所有新增字段均设置为 Optional，不设为 Required

**理由**：
- 保持向后兼容性，现有配置无需修改即可继续使用
- 这些字段对于实例的基本功能运行不是必需的
- 用户可以根据需要选择性配置这些字段

### 3. Update 逻辑实现

**决策**：在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，检查这三个字段的变更，并在变更时调用相应的 API

**理由**：
- 遵循 Terraform 的 update 模式：只有在字段值发生变化时才调用 API
- 减少不必要的 API 调用，提高性能
- 与现有 `cluster_name` 和 `resource_tags` 的实现方式保持一致

### 4. Read 逻辑实现

**决策**：在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，从 API 响应中读取这三个字段的值，并设置到 resource data 中

**理由**：
- 确保 Terraform state 与实际资源状态同步
- 支持资源的导入功能（Import）
- 遵循 Terraform 的 read 模式

### 5. Create 逻辑实现

**决策**：在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中，支持在创建实例时设置这三个字段

**理由**：
- 用户可以在创建实例时就配置这些字段，而不需要额外执行 update 操作
- 提高用户体验，减少操作步骤

## Risks / Trade-offs

### Risk 1: API 行为变更

**风险**：腾讯云 API 的行为可能在将来发生变化，导致实现与 API 不兼容

**缓解措施**：
- 使用腾讯云官方 SDK，确保与 API 的兼容性
- 在测试用例中验证 API 的行为
- 关注腾讯云的 API 更新公告，及时适配变更

### Risk 2: 删除保护导致的删除失败

**风险**：如果启用了删除保护，用户在执行 `terraform destroy` 时可能遇到失败，导致困惑

**缓解措施**：
- 在字段的 Description 中明确说明删除保护的影响
- 在删除失败的错误信息中明确提示用户删除保护已启用
- 在文档中提供删除保护的说明和操作指南

### Risk 3: 默认值处理

**风险**：如果用户不设置 `enable_deletion_protection` 或 `enable_risk_warning`，Terraform 的 Computed 字段可能返回 nil，导致状态不一致

**缓解措施**：
- 在 Read 函数中处理 nil 值，确保返回合理的默认值
- 在 Update 函数中检查字段是否真正需要更新，避免不必要的 API 调用

### Trade-off 1: 性能 vs. 完整性

**权衡**：在 Read 函数中，是否需要额外调用 API 获取这些字段

**决策**：不额外调用 API，直接从现有 API 响应中读取

**理由**：
- 腾讯云的 DescribeRabbitMQVipInstance API 返回的响应已经包含了这些字段
- 额外调用 API 会增加延迟和成本
- 从现有响应中读取即可满足需求

### Trade-off 2: 向后兼容性 vs. 新功能

**权衡**：是否应该将这些字段设为 Computed 而不是 Optional

**决策**：设为 Optional，不是 Computed

**理由**：
- 保持向后兼容性是最重要的
- Computed 字段在 Terraform 中通常用于 API 返回但用户不能设置的值
- 这三个字段用户应该可以设置和修改

## Implementation Details

### Schema 变更

在 `ResourceTencentCloudTdmqRabbitmqVipInstance` 函数的 Schema map 中新增：

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

### Create 函数变更

在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中，在创建 API 调用前设置这些字段：

```go
if v, ok := d.GetOkExists("remark"); ok {
    request.Remark = helper.String(v.(string))
}

if v, ok := d.GetOkExists("enable_deletion_protection"); ok {
    request.EnableDeletionProtection = helper.Bool(v.(bool))
}

if v, ok := d.GetOkExists("enable_risk_warning"); ok {
    request.EnableRiskWarning = helper.Bool(v.(bool))
}
```

**注意**：需要确认 `CreateRabbitMQVipInstanceRequest` 是否支持这些字段。如果不支持，则不需要在 Create 中设置。

### Read 函数变更

在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，从 API 响应中读取这些字段：

```go
if rabbitmqVipInstance.Remark != nil {
    _ = d.Set("remark", rabbitmqVipInstance.Remark)
}

if rabbitmqVipInstance.EnableDeletionProtection != nil {
    _ = d.Set("enable_deletion_protection", rabbitmqVipInstance.EnableDeletionProtection)
}

if rabbitmqVipInstance.EnableRiskWarning != nil {
    _ = d.Set("enable_risk_warning", rabbitmqVipInstance.EnableRiskWarning)
}
```

### Update 函数变更

在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中：

1. 从 `immutableArgs` 中移除这三个字段（如果它们被错误地添加到该列表）

2. 添加更新逻辑：

```go
if d.HasChange("remark") {
    if v, ok := d.GetOk("remark"); ok {
        request.Remark = helper.String(v.(string))
        needUpdate = true
    } else {
        // 如果 remark 被设置为空字符串，不更新该字段
        // 或者根据 API 行为决定是否设置为 nil
    }
}

if d.HasChange("enable_deletion_protection") {
    if v, ok := d.GetOkExists("enable_deletion_protection"); ok {
        request.EnableDeletionProtection = helper.Bool(v.(bool))
        needUpdate = true
    }
}

if d.HasChange("enable_risk_warning") {
    if v, ok := d.GetOkExists("enable_risk_warning"); ok {
        request.EnableRiskWarning = helper.Bool(v.(bool))
        needUpdate = true
    }
}
```

### 测试用例更新

更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`：

1. 新增测试场景：测试创建实例时设置 `remark`
2. 新增测试场景：测试更新 `remark`
3. 新增测试场景：测试启用和禁用 `enable_deletion_protection`
4. 新增测试场景：测试启用和禁用 `enable_risk_warning`
5. 新增测试场景：测试删除保护启用时的删除失败
6. 确保现有测试用例仍然通过

## API Compatibility

需要验证以下 API 字段的可用性：

### CreateRabbitMQVipInstanceRequest
- `Remark` - 需要确认是否支持
- `EnableDeletionProtection` - 需要确认是否支持
- `EnableRiskWarning` - 需要确认是否支持

### ModifyRabbitMQVipInstanceRequest
- `Remark` - ✅ 已确认支持
- `EnableDeletionProtection` - ✅ 已确认支持
- `EnableRiskWarning` - ✅ 已确认支持

### DescribeRabbitMQVipInstanceResponse
- `Remark` - 需要确认是否返回
- `EnableDeletionProtection` - 需要确认是否返回
- `EnableRiskWarning` - 需要确认是否返回

**行动项**：
1. 查看 `CreateRabbitMQVipInstanceRequest` 的定义，确认是否支持这些字段
2. 查看 `DescribeRabbitMQVipInstanceResponse` 的定义，确认是否返回这些字段
3. 如果不支持，需要调整实现策略

## Testing Strategy

### 单元测试
- 测试 Create 函数中这些字段的设置逻辑
- 测试 Read 函数中这些字段的读取逻辑
- 测试 Update 函数中这些字段的更新逻辑
- 测试边界情况：nil 值、空字符串、默认值

### 集成测试
- 测试完整的 CRUD 流程
- 测试字段的独立更新（只更新其中一个字段）
- 测试多个字段的联合更新（同时更新多个字段）
- 测试删除保护启用时的删除失败场景

### 验收测试
- 使用真实的腾讯云资源进行测试
- 确保所有测试用例通过
- 验证与 Terraform CLI 的兼容性
