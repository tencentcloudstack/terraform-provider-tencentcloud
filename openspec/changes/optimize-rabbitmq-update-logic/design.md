# Design: Optimize RabbitMQ Instance Update Logic

## Architecture Overview

本文档详细说明优化 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源 Update 逻辑的技术设计和实现方案。

## Current Implementation Analysis

### 现有 Update 函数结构

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    logId := tccommon.GetLogId(ctx)
    request := tdmq.NewModifyRabbitMQVipInstanceRequest()

    // 填充请求参数
    request.ClusterId = helper.String(d.Id())
    request.ClusterName = helper.String(d.Get("cluster_name").(string))
    request.Remark = helper.String(d.Get("remark").(string))
    // ... 其他字段

    // 调用 API
    _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstance(request)
    if err != nil {
        return diag.FromErr(err)
    }

    // 调用 Read 刷新状态
    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(ctx, d, meta)
}
```

### 存在的问题

1. **无条件调用 API**：无论配置是否变化都会调用 `ModifyRabbitMQVipInstance`
2. **参数冗余**：所有参数都包含在请求中，即使没有变化
3. **无超时控制**：长时间运行的 Update 操作可能卡住
4. **状态刷新低效**：每次 Update 后都调用 Read，即使只是更新了状态

## Optimization Design

### 1. 分层架构设计

```
Update 函数 (入口层)
    ├── 变更检测层 (Change Detection)
    │   ├── 检测 cluster_name 变更
    │   ├── 检测 remark 变更
    │   ├── 检测 spec 变更
    │   └── 检测其他可更新属性
    ├── 更新执行层 (Update Execution)
    │   ├── updateClusterName()
    │   ├── updateRemark()
    │   ├── updateSpec()
    │   └── updateOtherProps()
    └── 状态管理层 (State Management)
        ├── 本地状态更新
        └── 最终 Read 调用
```

### 2. 优化后的 Update 函数结构

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    logId := tccommon.GetLogId(ctx)

    // 创建带有超时的上下文
    updateCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
    defer cancel()

    // 变更检测和执行
    hasChanges := false

    // 1. 检测并更新集群名称
    if d.HasChange("cluster_name") {
        if err := updateClusterName(updateCtx, d, meta); err != nil {
            return diag.FromErr(fmt.Errorf("update cluster name failed: %v", err))
        }
        hasChanges = true
    }

    // 2. 检测并更新备注
    if d.HasChange("remark") {
        if err := updateRemark(updateCtx, d, meta); err != nil {
            return diag.FromErr(fmt.Errorf("update remark failed: %v", err))
        }
        hasChanges = true
    }

    // 3. 检测并更新规格
    if d.HasChange("node_count") || d.HasChange("spec_name") {
        if err := updateSpec(updateCtx, d, meta); err != nil {
            return diag.FromErr(fmt.Errorf("update spec failed: %v", err))
        }
        hasChanges = true
    }

    // 4. 检测并更新自动续费
    if d.HasChange("auto_renew_flag") {
        if err := updateAutoRenewFlag(updateCtx, d, meta); err != nil {
            return diag.FromErr(fmt.Errorf("update auto renew flag failed: %v", err))
        }
        hasChanges = true
    }

    // 如果没有实际变更，直接返回
    if !hasChanges {
        log.Printf("[WARN] no actual changes detected for rabbitmq instance %s", d.Id())
        return nil
    }

    // 最终刷新状态
    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(updateCtx, d, meta)
}
```

### 3. 独立更新函数实现

#### 3.1 更新集群名称

```go
func updateClusterName(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
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

        // 等待任务完成
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

    // 更新本地状态
    return d.Set("cluster_name", d.Get("cluster_name"))
}
```

#### 3.2 更新备注

```go
func updateRemark(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    logId := tccommon.GetLogId(ctx)
    request := tdmq.NewModifyRabbitMQVipInstanceRequest()

    request.ClusterId = helper.String(d.Id())
    request.Remark = helper.String(d.Get("remark").(string))

    // 使用重试机制
    err := resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstance(request)
        if err != nil {
            return tccommon.RetryError(ctx, err)
        }

        // 等待任务完成
        if response.Response.TaskId != nil {
            taskId := *response.Response.TaskId
            if err := waitForTaskCompletion(ctx, meta, taskId); err != nil {
                return resource.NonRetryableError(fmt.Errorf("wait for task failed: %v", err))
            }
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("modify remark failed: %v", err)
    }

    // 更新本地状态
    return d.Set("remark", d.Get("remark"))
}
```

#### 3.3 更新规格

```go
func updateSpec(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    logId := tccommon.GetLogId(ctx)
    request := tdmq.NewModifyRabbitMQVipInstanceSpecRequest()

    request.ClusterId = helper.String(d.Id())

    // 设置新的规格参数
    if d.HasChange("node_count") {
        request.NodeCount = helper.IntUint64(uint64(d.Get("node_count").(int)))
    }

    if d.HasChange("spec_name") {
        request.SpecName = helper.String(d.Get("spec_name").(string))
    }

    // 使用重试机制
    err := resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstanceSpec(request)
        if err != nil {
            return tccommon.RetryError(ctx, err)
        }

        // 等待任务完成
        if response.Response.TaskId != nil {
            taskId := *response.Response.TaskId
            if err := waitForTaskCompletion(ctx, meta, taskId); err != nil {
                return resource.NonRetryableError(fmt.Errorf("wait for task failed: %v", err))
            }
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("modify spec failed: %v", err)
    }

    return nil
}
```

#### 3.4 更新自动续费

```go
func updateAutoRenewFlag(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
    logId := tccommon.GetLogId(ctx)
    request := billing.NewModifyAutoRenewFlagRequest()

    request.ResourceId = helper.String(d.Id())
    request.AutoRenewFlag = helper.IntUint64(uint64(d.Get("auto_renew_flag").(int)))

    // 使用重试机制
    err := resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *resource.RetryError {
        _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingClient().ModifyAutoRenewFlag(request)
        if err != nil {
            return tccommon.RetryError(ctx, err)
        }
        return nil
    })

    if err != nil {
        return fmt.Errorf("modify auto renew flag failed: %v", err)
    }

    return nil
}
```

### 4. 任务等待函数

```go
func waitForTaskCompletion(ctx context.Context, meta interface{}, taskId string) error {
    request := tdmq.NewDescribeTaskDetailRequest()
    request.TaskId = helper.String(taskId)

    // 任务状态常量
    const (
        taskStatusSuccess = "success"
        taskStatusRunning = "running"
        taskStatusFailed  = "failed"
    )

    err := resource.RetryContext(ctx, tccommon.ReadRetryTimeout, func() *resource.RetryError {
        response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().DescribeTaskDetail(request)
        if err != nil {
            return tccommon.RetryError(ctx, err)
        }

        if response.Response.TaskStatus == nil {
            return resource.NonRetryableError(fmt.Errorf("task status is nil"))
        }

        status := *response.Response.TaskStatus

        switch status {
        case taskStatusSuccess:
            return nil
        case taskStatusRunning:
            return resource.RetryableError(fmt.Errorf("task is still running"))
        case taskStatusFailed:
            if response.Response.ErrorMessage != nil {
                return resource.NonRetryableError(fmt.Errorf("task failed: %s", *response.Response.ErrorMessage))
            }
            return resource.NonRetryableError(fmt.Errorf("task failed"))
        default:
            return resource.NonRetryableError(fmt.Errorf("unknown task status: %s", status))
        }
    })

    return err
}
```

### 5. 错误处理增强

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    logId := tccommon.GetLogId(ctx)
    defer tccommon.LogElapsed(logId + "update")

    // 记录开始时间
    startTime := time.Now()

    // ... 更新逻辑 ...

    // 记录耗时
    duration := time.Since(startTime)
    log.Printf("[INFO] rabbitmq instance %s update completed in %v", d.Id(), duration)

    return nil
}
```

### 6. Schema 超时配置

```go
func resourceTencentCloudTdmqRabbitmqVipInstance() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudTdmqRabbitmqVipInstanceCreate,
        Read:   resourceTencentCloudTdmqRabbitmqVipInstanceRead,
        Update: resourceTencentCloudTdmqRabbitmqVipInstanceUpdate,
        Delete: resourceTencentCloudTdmqRabbitmqVipInstanceDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Timeout: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(30 * time.Minute),
            Update: schema.DefaultTimeout(20 * time.Minute),  // 为 Update 添加超时
            Delete: schema.DefaultTimeout(30 * time.Minute),
        },
        Schema: map[string]*schema.Schema{
            // ... 现有 Schema 字段 ...
        },
    }
}
```

## Performance Analysis

### 优化前的性能

- **API 调用次数**：每次 Update 至少 2 次（1 次 Modify + 1 次 Read）
- **网络延迟**：~500ms × 2 = 1000ms
- **不必要调用**：即使配置未变也会调用 API

### 优化后的性能

- **API 调用次数**：只在有实际变更时调用，最少 1 次（Read）
- **网络延迟**：~500ms × N（N 为实际变更的属性数）
- **无变更时**：仅 1 次 Read 调用，~500ms

**性能提升**：
- 无变更场景：50% 减少延迟（1000ms → 500ms）
- 单属性变更：相同延迟，但代码更清晰
- 多属性变更：取决于 API 设计，可能批量或分别调用

## Compatibility Considerations

### 向后兼容性

1. **配置兼容**：现有 Terraform 配置无需修改
2. **状态兼容**：现有状态可以正常迁移
3. **API 兼容**：仍然使用相同的 TencentCloud API

### 潜在影响

1. **错误消息变化**：更具体的错误消息（例如："update cluster name failed" vs "modify instance failed"）
2. **超时行为**：新增的超时配置可能导致某些长时间操作失败
3. **日志格式**：新增的性能日志

## Testing Strategy

### 单元测试

```go
func TestUpdateClusterName(t *testing.T) {
    // 测试集群名称更新逻辑
}

func TestUpdateRemark(t *testing.T) {
    // 测试备注更新逻辑
}

func TestUpdateSpec(t *testing.T) {
    // 测试规格更新逻辑
}

func TestWaitForTaskCompletion(t *testing.T) {
    // 测试任务等待逻辑
}
```

### 验收测试

```go
func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_clusterName(t *testing.T) {
    // 测试更新集群名称
}

func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_remark(t *testing.T) {
    // 测试更新备注
}

func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_spec(t *testing.T) {
    // 测试更新规格
}

func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_multipleProps(t *testing.T) {
    // 测试同时更新多个属性
}

func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_noChanges(t *testing.T) {
    // 测试无变更场景
}

func TestAccTencentCloudTdmqRabbitmqVipInstanceUpdate_timeout(t *testing.T) {
    // 测试超时配置
}
```

### 性能测试

```go
func BenchmarkUpdateLogic(b *testing.B) {
    // 基准测试，对比优化前后的性能
}
```

## Monitoring and Observability

### 日志增强

```go
log.Printf("[INFO] starting update for rabbitmq instance %s", instanceId)
log.Printf("[INFO] updating cluster name from %s to %s", oldName, newName)
log.Printf("[INFO] updating remark")
log.Printf("[INFO] update completed in %v", duration)
```

### 指标收集

建议收集以下指标：
- Update 操作的成功率
- Update 操作的平均耗时
- 各属性更新的频率
- 超时失败的数量

## Migration Guide

### 用户迁移步骤

1. **升级 Provider**：升级到包含此优化的版本
2. **验证配置**：运行 `terraform plan` 验证配置兼容性
3. **应用变更**：运行 `terraform apply` 应用任何必要的变更
4. **监控日志**：关注日志中的警告或错误消息

### 开发者迁移步骤

1. **Review 代码变更**：理解新的 Update 逻辑
2. **运行测试**：确保所有测试通过
3. **更新文档**：更新 API 文档和用户指南

## Future Enhancements

1. **批量更新优化**：如果 API 支持，实现批量更新多个属性
2. **条件更新**：实现条件性更新（例如：只在值真正改变时更新）
3. **并行更新**：对于独立的属性更新，实现并行执行
4. **智能回滚**：实现部分失败时的自动回滚机制
5. **缓存优化**：添加本地缓存减少不必要的 Read 调用
