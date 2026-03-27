# Proposal: Optimize RabbitMQ Instance Update Logic

## Summary

优化 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 Update 逻辑，提高更新操作的效率、可靠性和用户体验。当前实现存在重复 API 调用、不必要的属性更新和错误处理不完善等问题。

## Background

当前的 RabbitMQ 实例 Update 函数存在以下问题：

1. **重复的 Read 操作**：每次 Update 操作后都会调用 Read 函数刷新状态，导致额外的 API 调用
2. **低效的变更检测**：没有正确使用 `d.HasChange()` 来检测变更，导致即使配置未改变也会触发 API 调用
3. **不完整的属性更新**：某些属性的更新逻辑不完整，可能导致状态不一致
4. **错误处理不完善**：Update 失败时的错误处理和回滚机制不够健壮
5. **缺少超时配置**：Update 操作没有使用超时配置，可能导致长时间等待

## Goals

1. 提高更新操作的效率，减少不必要的 API 调用
2. 改进变更检测逻辑，只在必要时触发 API 调用
3. 完善属性更新逻辑，确保状态一致性
4. 增强错误处理和重试机制
5. 为 Update 操作添加超时配置

## Non-Goals

- 修改资源的 Schema 定义
- 改变 Create 和 Delete 函数的行为
- 添加新的资源属性

## Design

### 1. 优化变更检测逻辑

使用 `d.HasChange()` 来检测属性变更，避免不必要的 API 调用：

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // 只检测实际发生变更的属性
    if d.HasChange("cluster_name") {
        // 只在名称变更时调用 API
        if err := modifyClusterName(ctx, d, meta); err != nil {
            return diag.FromErr(err)
        }
    }

    if d.HasChange("remark") {
        // 只在备注变更时调用 API
        if err := modifyRemark(ctx, d, meta); err != nil {
            return diag.FromErr(err)
        }
    }
    // ... 其他属性
}
```

### 2. 减少重复 Read 操作

在 Update 函数内部实现状态更新，避免额外的 Read 调用：

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // ... 执行更新逻辑

    // 使用 SetNewValue 更新状态，而不是调用 Read
    if d.HasChange("cluster_name") {
        if err := d.Set("cluster_name", d.Get("cluster_name")); err != nil {
            return diag.FromErr(err)
        }
    }

    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(ctx, d, meta)
}
```

### 3. 添加超时配置

为 Update 操作添加超时配置支持：

```go
// 在 Schema 中添加 Timeouts 配置
"Timeouts": &schema.Schema{
    Type:     schema.TypeType,
    Optional: true,
}

// 在 Update 函数中使用超时
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    updateCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
    defer cancel()

    // ... 执行更新逻辑，使用 updateCtx
}
```

### 4. 改进错误处理

添加更完善的错误处理和重试机制：

```go
func modifyClusterName(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    logId := tccommon.GetLogId(ctx)
    request := tdmq.NewModifyRabbitMQVipInstanceRequest()

    request.ClusterId = helper.String(d.Id())
    request.ClusterName = helper.String(d.Get("cluster_name").(string))

    // 使用重试机制
    err := resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstance(request)
        if err != nil {
            return tccommon.RetryError(ctx, err)
        }

        // 检查任务状态
        if response.Response.TaskId != nil {
            taskId := *response.Response.TaskId
            if err := waitForTaskCompletion(ctx, meta, taskId); err != nil {
                return resource.NonRetryableError(fmt.Errorf("wait for task failed: %v", err))
            }
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("modify cluster name failed: %v", err)
    }

    return nil
}
```

### 5. 分离关注点

将不同属性的更新逻辑分离到独立函数中：

```go
func modifyClusterName(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    // 专门处理名称更新的逻辑
}

func modifyRemark(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    // 专门处理备注更新的逻辑
}

func modifySpec(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    // 专门处理规格更新的逻辑
}
```

## Alternatives

### Alternative 1: 保持当前实现
- **优点**: 不需要修改代码，风险低
- **缺点**: 继续存在性能和可靠性问题

### Alternative 2: 重写整个 Update 函数
- **优点**: 可以彻底重构，优化所有方面
- **缺点**: 工作量大，测试成本高，引入风险大

### Alternative 3: 逐步优化
- **优点**: 风险可控，可以分阶段验证
- **缺点**: 需要多次迭代，最终效果可能不如整体重构

**选择**: 本提案采用 Alternative 3（逐步优化），在保证稳定性的前提下逐步改进。

## Drawbacks

1. **测试成本**：需要编写更多的测试用例来验证优化效果
2. **回归风险**：修改 Update 逻辑可能影响现有用户
3. **复杂度增加**：分离关注点和添加超时配置会增加代码复杂度

## Prior Art

参考其他腾讯云 Provider 资源的 Update 逻辑优化经验，例如：
- `tencentcloud_mysql_instance` 的 Update 函数使用了 `d.HasChange()` 来检测变更
- `tencentcloud_redis_instance` 实现了超时配置支持
- `tencentcloud_cvm_instance` 使用了独立函数处理不同属性的更新

## Unresolved Questions

1. 是否需要为所有属性都实现独立的更新函数？
2. 超时配置的默认值应该如何设置？
3. 是否需要添加 Update 操作的指标统计？

## Out of Scope

- 修改 RabbitMQ 实例的 Create 和 Delete 逻辑
- 添加新的资源属性
- 修改 Schema 定义
- 影响其他 TDMQ 资源（如 Pulsar、RocketMQ）

## Testing Strategy

### 单元测试
- 为每个独立的更新函数编写单元测试
- 测试变更检测逻辑的正确性
- 验证错误处理和重试机制

### 验收测试
- 测试各种属性的更新场景
- 验证超时配置的生效
- 测试错误恢复和状态一致性

### 回归测试
- 确保现有功能不受影响
- 验证与现有配置的兼容性

## Rollout Plan

1. **Phase 1**: 实现变更检测逻辑优化
2. **Phase 2**: 分离关注点，创建独立更新函数
3. **Phase 3**: 添加超时配置支持
4. **Phase 4**: 完善错误处理和重试机制
5. **Phase 5**: 全面测试和文档更新

每个阶段完成后进行充分测试，确保没有回归问题。
