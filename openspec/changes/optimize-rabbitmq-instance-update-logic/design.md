## Context

### Background

Tencent Cloud Terraform Provider 中的 RabbitMQ VIP 实例资源 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 当前支持有限的更新能力。通过分析腾讯云 SDK (tencentcloud-sdk-go) 的 `ModifyRabbitMQVipInstance` API，发现该 API 实际上支持更多参数的更新，但 Terraform Provider 的实现未完全利用这些能力。

### Current State

当前 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数的实现：

1. **硬编码的不可变参数列表**：包含大量参数（zone_ids, vpc_id, subnet_id, node_spec, node_num, storage_size 等），这些参数在更新时会直接返回错误。

2. **仅支持有限的更新参数**：只支持 `cluster_name` 和 `resource_tags` 的更新。

3. **缺失腾讯云 API 支持的参数**：
   - `Remark`：实例备注信息
   - `EnableDeletionProtection`：是否开启删除保护
   - `EnableRiskWarning`：是否开启集群风险提示

4. **代码结构问题**：
   - 标签更新逻辑较为复杂，需要优化
   - 参数处理逻辑分散，缺乏统一的变更检测机制
   - 错误处理和日志记录可以改进

### Stakeholders

- Terraform 用户管理 RabbitMQ 实例
- DevOps 团队需要灵活配置实例属性
- 安全团队需要控制删除保护设置

### Constraints

- **必须保持向后兼容**：不能破坏现有 Terraform 配置和状态
- **必须遵循现有模式**：与 Provider 中其他资源的实现保持一致
- **API 限制**：只能使用 `ModifyRabbitMQVipInstance` API 支持的参数
- **Schema 限制**：只能新增 Optional 字段，不能修改已有字段的属性

## Goals / Non-Goals

### Goals

1. **增强更新能力**：支持通过 Terraform 管理实例的备注、删除保护和风险提示开关
2. **优化代码结构**：重构 update 函数，提高代码的可维护性和可扩展性
3. **改进错误处理**：增强错误处理和日志记录，便于调试和故障排查
4. **完善测试覆盖**：为新增功能添加单元测试

### Non-Goals

- **不支持规格变更**：node_spec、node_num、storage_size 等规格参数的更新（腾讯云 API 不支持）
- **不支持基础设施变更**：zone_ids、vpc_id、subnet_id 等基础设施参数的更新（腾讯云 API 不支持）
- **不支持付费模式变更**：pay_mode、time_span、auto_renew_flag 等付费相关参数的更新（腾讯云 API 不支持）
- **不支持版本升级**：cluster_version 的更新（腾讯云 API 不支持）
- **不支持公网配置变更**：band_width、enable_public_access 的更新（腾讯云 API 不支持）

## Decisions

### Decision 1: 新增参数的 Schema 定义

**Choice**：为三个新参数添加 Optional 类型的 schema 字段

**Rationale**：
- 这三个参数都是 Optional 的，不破坏向后兼容性
- 用户可以选择性地配置这些属性
- 符合 Terraform Provider 的常见模式

**Schema 设计**：
```go
"remark": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    Description: "Instance remark information.",
},
"enable_deletion_protection": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Whether to enable deletion protection. Default is false.",
},
"enable_risk_warning": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Whether to enable cluster risk warning.",
},
```

**Alternatives Considered**：
1. **不添加 schema 字段**：用户无法通过 Terraform 管理这些属性
   - **Rejected**：违背了本次变更的目标
2. **添加为 Required 字段**：强制用户必须配置
   - **Rejected**：破坏向后兼容性，且这些属性应该有默认值

### Decision 2: 不可变参数列表的调整

**Choice**：保持当前的不可变参数列表不变

**Rationale**：
- 通过验证 SDK 确认，`ModifyRabbitMQVipInstance` API 确实不支持这些参数的更新
- 保持当前列表可以避免用户尝试更新不支持的参数而得到 API 错误
- 在 Terraform 层面提前返回错误，提供更好的用户体验

**不可变参数**：
```go
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num",
    "storage_size", "enable_create_default_ha_mirror_queue",
    "auto_renew_flag", "time_span", "pay_mode", "cluster_version",
    "band_width", "enable_public_access",
}
```

**Alternatives Considered**：
1. **移除不可变参数检查**：允许用户尝试更新这些参数，让 API 返回错误
   - **Rejected**：用户体验差，API 错误不够友好
2. **动态检测 API 支持的参数**：根据 API 返回的错误来判断哪些参数不可变
   - **Rejected**：增加了复杂度，且不可靠（API 错误可能因其他原因）

### Decision 3: Update 函数的结构重构

**Choice**：采用参数处理函数提取的方式重构

**Rationale**：
- 将每个参数的处理逻辑封装成独立函数，提高代码可读性
- 便于未来添加新的可更新参数
- 遵循单一职责原则

**重构后的结构**：
```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
    // 设置日志和不一致性检查
    defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
    defer tccommon.InconsistentCheck(d, meta)()

    // 检查不可变参数
    if err := checkImmutableArgs(d); err != nil {
        return err
    }

    // 构建更新请求
    request, needUpdate := buildUpdateRequest(d)

    // 执行更新
    if needUpdate {
        if err := executeUpdate(request, d, meta); err != nil {
            return err
        }
    }

    // 读取最新状态
    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}
```

**Alternatives Considered**：
1. **保持当前的单体函数**：所有逻辑都在一个函数中
   - **Rejected**：可读性差，难以维护
2. **使用策略模式**：为每个参数类型创建策略类
   - **Rejected**：过度设计，增加了不必要的复杂度

### Decision 4: 标签更新逻辑的优化

**Choice**：保持当前的标签更新逻辑，但提取为独立函数

**Rationale**：
- 当前的标签更新逻辑已经考虑了标签添加和删除的场景
- 将其提取为独立函数可以改善代码结构
- 与其他参数的处理逻辑保持一致

**标签更新逻辑**：
```go
func handleTagsUpdate(d *schema.ResourceData, request *tdmq.ModifyRabbitMQVipInstanceRequest) bool {
    if d.HasChange("resource_tags") {
        request.Tags = []*tdmq.Tag{}
        if v, ok := d.GetOk("resource_tags"); ok {
            for _, item := range v.([]interface{}) {
                dMap := item.(map[string]interface{})
                tag := tdmq.Tag{}
                if v, ok := dMap["tag_key"]; ok {
                    tag.TagKey = helper.String(v.(string))
                }
                if v, ok := dMap["tag_value"]; ok {
                    tag.TagValue = helper.String(v.(string))
                }
                request.Tags = append(request.Tags, &tag)
            }
        } else {
            request.RemoveAllTags = helper.Bool(true)
        }
        return true
    }
    return false
}
```

**Alternatives Considered**：
1. **重构为增量更新**：只发送变更的标签
   - **Rejected**：API 期望接收完整的标签列表，不支持增量更新
2. **使用 map 类型**：将标签从 TypeList 改为 TypeMap
   - **Rejected**：破坏向后兼容性，且现有用户已使用 TypeList

### Decision 5: Read 函数的更新

**Choice**：在 Read 函数中添加对新参数的读取逻辑

**Rationale**：
- 确保 Terraform 状态与实际云资源保持同步
- 在 schema 中标记为 Computed，需要从 API 读取后设置到状态中

**新增读取逻辑**：
```go
// 在 Read 函数中添加
if rabbitmqVipInstance.ClusterInfo.Remark != nil {
    _ = d.Set("remark", rabbitmqVipInstance.ClusterInfo.Remark)
}
if rabbitmqVipInstance.ClusterInfo.EnableDeletionProtection != nil {
    _ = d.Set("enable_deletion_protection", rabbitmqVipInstance.ClusterInfo.EnableDeletionProtection)
}
if rabbitmqVipInstance.ClusterInfo.EnableRiskWarning != nil {
    _ = d.Set("enable_risk_warning", rabbitmqVipInstance.ClusterInfo.EnableRiskWarning)
}
```

**Alternatives Considered**：
1. **不读取这些参数**：只写入不读取
   - **Rejected**：会导致状态不一致，Terraform 认为这些参数每次都需要更新
2. **使用 API 响应的不同字段**：从其他结构中读取
   - **Rejected**：需要验证 API 响应结构，当前方案最直接

## Risks / Trade-offs

### Risk 1: SDK 版本兼容性

**Issue**：如果使用的 SDK 版本较旧，可能不支持新的参数

**Likelihood**：低

**Mitigation**：
- 在实施前验证 SDK 版本，确认包含新参数
- 如果 SDK 版本不匹配，先升级 SDK 依赖

**Impact**：中等 - 可能需要升级 SDK，但不会破坏现有功能

### Risk 2: API 行为变化

**Issue**：腾讯云可能在未来修改 API 的行为，导致某些参数不再支持更新

**Likelihood**：低

**Mitigation**：
- 在文档中明确标注支持的更新参数
- 监控 API 变更通知
- 如果 API 发生变化，及时更新 Provider

**Impact**：中等 - 可能需要调整实现，但不会破坏现有功能

### Risk 3: 用户配置迁移问题

**Issue**：现有用户可能已经在 Terraform 外部通过 API 修改了这些参数

**Likelihood**：中

**Mitigation**：
- 新参数设置为 Computed，首次读取时会同步实际值
- 在文档中说明如何处理状态不一致的情况

**Impact**：低 - 不影响现有功能，只是状态同步

### Trade-off 1: 代码重构的范围

**Chosen**：仅重构 Update 函数，保持其他函数不变

**Trade-off**：减少重构范围 vs 可能遗漏其他可优化的地方

**Justification**：聚焦于本次变更的目标，降低风险，避免过度重构

### Trade-off 2: 新参数的默认值

**Chosen**：不设置默认值，使用 API 返回的实际值

**Trade-off**：明确性 vs 灵活性

**Justification**：让用户明确感知这些参数的存在，避免隐式的默认值行为

## Migration Plan

### For New Users

无需迁移 - 新功能为 opt-in（Optional 字段），新用户可以直接使用。

### For Existing Instances

**场景 1**：用户已有 RabbitMQ 实例，且未配置新参数

**行为**：
1. 第一次 `terraform plan` 不会显示任何变化（新参数未配置）
2. 用户可以选择性地添加 `remark`、`enable_deletion_protection`、`enable_risk_warning` 到配置中
3. `terraform plan` 会显示参数添加
4. `terraform apply` 会调用 `ModifyRabbitMQVipInstance` API 来更新这些参数

**状态迁移**：不需要 - 新字段为 Optional，默认行为安全

**场景 2**：用户已有 RabbitMQ 实例，且已通过 API 修改了这些参数

**行为**：
1. 第一次 `terraform plan` 会显示参数添加（因为在 state 中不存在）
2. 用户可以选择接受这些更改（同步 API 值到 state）
3. `terraform apply` 会将 API 返回的实际值设置到 state 中

**状态迁移**：不需要 - Computed 字段会自动同步 API 值

### Rollback

如果功能导致问题：

1. **代码回滚**：恢复到变更前的代码版本
2. **配置调整**：用户可以从配置中移除新参数
3. **状态清理**：使用 `terraform state rm` 命令移除 state 中的新参数字段

## Open Questions

无未决问题 - 所有必要的技术决策已在提案和设计中明确。
