## Context

### Background

Terraform Provider for TencentCloud 的 RabbitMQ VIP 实例资源 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 目前支持通过 `ModifyRabbitMQVipInstance` API 修改部分字段。但当前实现不完整，遗漏了 API 支持的几个重要字段：

- `remark`: 实例备注
- `enable_deletion_protection`: 删除保护
- `enable_risk_warning`: 风险提示

腾讯云 TDMQ API 的 `ModifyRabbitMQVipInstance` 接口（SDK v20200217）已支持这些字段的修改，但 Provider 层未实现。

### Current State

**当前实现位置**：`tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`

**Update 函数现状**：
```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
    // 仅支持修改 cluster_name 和 resource_tags
    // immutableArgs 列表包含多个字段
}
```

**问题分析**：
1. Schema 中未定义 `remark`、`enable_deletion_protection`、`enable_risk_warning` 字段
2. Update 函数未处理这些字段的变化
3. Read 函数未从 API 响应中读取这些字段
4. `immutableArgs` 列表包含了 `auto_renew_flag` 等不应在 update 阶段严格禁止的字段

### Constraints

1. **向后兼容性**：必须保持 100% 向后兼容，不能破坏现有 Terraform 配置和 state
2. **Schema 限制**：只能新增 Optional/Computed 字段，不能修改或删除现有字段
3. **API 依赖**：依赖腾讯云 TDMQ SDK (v20200217) 的完整 API 支持
4. **代码规范**：遵循项目现有的代码模式和错误处理机制

## Goals / Non-Goals

### Goals

1. **完善 Schema 定义**：在 RabbitMQ VIP 实例资源中添加 3 个新字段的 Schema 定义
2. **增强 Update 逻辑**：在 Update 函数中支持新字段的修改
3. **完善 Read 逻辑**：在 Read 函数中读取新字段的值
4. **优化错误处理**：改进 `immutableArgs` 列表，确保合理的字段禁止策略

### Non-Goals

1. 不修改 RabbitMQ 实例的 Create 和 Delete 逻辑
2. 不添加新的 Terraform Provider 资源或数据源
3. 不引入新的外部依赖或 API 调用
4. 不修改测试文件（测试更新将作为独立任务）

## Decisions

### 1. Schema 字段设计

**决策**：新字段定义为 Optional + Computed，以确保向后兼容性

**理由**：
- `Optional`：允许用户在配置中显式设置这些字段
- `Computed`：允许用户不设置这些字段时从 API 读取实际值
- 这种组合确保既支持用户配置，也支持 API 返回值的读取

**Schema 定义**：
```go
"remark": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    Description: "Remarks for the RabbitMQ instance.",
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

### 2. Update 逻辑实现

**决策**：使用 `d.HasChange()` 检测字段变化，仅在有变化时调用 API

**理由**：
- 遵循项目的标准模式（参考 `cluster_name` 的处理）
- 避免不必要的 API 调用，提高性能
- 使用 `needUpdate` 标志控制 API 调用

**实现模式**：
```go
if d.HasChange("remark") {
    if v, ok := d.GetOk("remark"); ok {
        request.Remark = helper.String(v.(string))
        needUpdate = true
    }
}
// 类似处理 enable_deletion_protection 和 enable_risk_warning
```

### 3. Read 逻辑实现

**决策**：从 API 响应的两个来源读取新字段

**理由**：
- `DescribeRabbitMQVipInstance` API 返回详细信息（包含 remark）
- `DescribeRabbitMQVipInstances` API 返回列表信息（包含 enable_deletion_protection 和 enable_risk_warning）
- 需要从两个 API 响应中分别读取不同字段

**读取映射**：
```go
// 从 rabbitmqVipInstance (DescribeRabbitMQVipInstance) 读取
if rabbitmqVipInstance.ClusterInfo.Remark != nil {
    _ = d.Set("remark", rabbitmqVipInstance.ClusterInfo.Remark)
}

// 从 result[0] (DescribeRabbitMQVipInstances) 读取
if result[0].EnableDeletionProtection != nil {
    _ = d.Set("enable_deletion_protection", result[0].EnableDeletionProtection)
}
if result[0].EnableRiskWarning != nil {
    _ = d.Set("enable_risk_warning", result[0].EnableRiskWarning)
}
```

### 4. immutableArgs 列表优化

**决策**：保留当前列表不变，仅在必要时调整

**理由**：
- `auto_renew_flag` 虽然可以通过其他 API 修改，但 Terraform Provider 中通常在 update 时禁止修改
- 腾讯云的 `ModifyRabbitMQVipInstance` API 确实不支持修改大部分基础设施参数
- 过度放宽限制可能导致用户误操作

**当前列表**（保持不变）：
```go
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num",
    "storage_size", "enable_create_default_ha_mirror_queue",
    "auto_renew_flag", "time_span", "pay_mode", "cluster_version",
    "band_width", "enable_public_access",
}
```

### 5. 错误处理策略

**决策**：保持现有的错误处理模式，不引入新的错误处理逻辑

**理由**：
- 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行重试
- 使用 `tccommon.RetryError(e)` 处理可重试错误
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行日志和一致性检查
- 这些模式在项目中广泛使用，保持一致性

## Risks / Trade-offs

### Risk 1: API 返回字段可能为 nil

**描述**：API 可能不返回某些新字段（如 remark），导致读取时出现空指针异常

**缓解措施**：
- 使用 nil 检查：`if rabbitmqVipInstance.ClusterInfo.Remark != nil`
- Schema 中字段定义为 Optional + Computed，允许 nil 值
- 遵循项目中现有代码的 nil 处理模式

### Risk 2: 新字段名称与未来 API 变更冲突

**描述**：腾讯云 API 可能在未来更改字段名称或结构

**缓解措施**：
- 使用 SDK 类型（如 `*string`, `*bool`）而非原始类型
- 依赖 SDK 的版本化，SDK 会处理 API 变更
- 字段名称与腾讯云 API 文档保持一致

### Trade-off: 字段是否强制用户配置

**选择**：Optional + Computed（不强制配置，支持读取）

**替代方案**：
- **Optional only**：用户配置才生效，不读取 API 值 → 用户体验差
- **Required**：强制用户配置 → 破坏向后兼容性

**理由**：Optional + Computed 既提供灵活性，又保证完整的功能支持

## Migration Plan

### 部署步骤

1. **代码变更**：修改 `resource_tc_tdmq_rabbitmq_vip_instance.go`
   - 添加 3 个 Schema 字段定义
   - 更新 Update 函数
   - 更新 Read 函数

2. **测试验证**：
   - 运行单元测试：`go test -v ./tencentcloud/services/trabbit/`
   - 运行验收测试：`TF_ACC=1 go test -v ./tencentcloud/services/trabbit/`
   - 手动测试：创建实例并验证新字段的 update 和 read 功能

3. **文档更新**（可选）：
   - 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 文档
   - 更新 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` 文档

### Rollback 策略

- 代码变更是向后兼容的，可以直接回滚
- 如果 API 行为异常，可以暂时在 Update 中跳过这些字段的处理
- 不涉及数据库迁移或状态迁移，回滚风险低

## Open Questions

无。所有技术决策已在上述章节中明确。
