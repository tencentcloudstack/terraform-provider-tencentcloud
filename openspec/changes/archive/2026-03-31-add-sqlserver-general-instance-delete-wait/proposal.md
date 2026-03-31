# Proposal: Add Isolation Wait Mechanism for SQL Server General Instance Deletion

## Overview

为 `tencentcloud_sqlserver_general_cloud_instance` 资源的删除流程添加等待隔离成功的机制。在执行 `TerminateSqlserverInstanceById` 后，需要等待实例状态变为隔离状态（Status = 4）后，再执行 `DeleteSqlserverInstanceById`。

## Problem Statement

当前 `tencentcloud_sqlserver_general_cloud_instance` 资源的删除模块存在以下问题：

```go
func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
    // ...
    
    // 调用 TerminateSqlserverInstanceById 隔离实例
    if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
        return err
    }
    
    // 立即调用 DeleteSqlserverInstanceById 删除实例
    if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
        return err
    }
    
    return nil
}
```

**当前问题**：
1. `TerminateSqlserverInstanceById` 调用后立即执行 `DeleteSqlserverInstanceById`
2. 没有等待实例隔离成功（Status = 4）
3. 可能导致删除操作失败或不稳定

## Proposed Solution

在 `TerminateSqlserverInstanceById` 和 `DeleteSqlserverInstanceById` 之间添加等待机制：

```go
func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
    // ...
    
    // Step 1: 隔离实例
    if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
        return err
    }
    
    // Step 2: 等待隔离成功（新增）
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        instance, err := service.DescribeDBInstanceById(ctx, instanceId)
        if err != nil {
            return tccommon.RetryError(err)
        }
        
        if instance.Status != nil && *instance.Status == 4 {
            // 隔离成功
            return nil
        }
        
        // 继续等待
        return resource.RetryableError(fmt.Errorf("waiting for instance to be isolated"))
    })
    
    if err != nil {
        return err
    }
    
    // Step 3: 删除实例
    if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
        return err
    }
    
    return nil
}
```

## Goals

1. ✅ 在删除流程中添加等待隔离成功的机制
2. ✅ 调用 `DescribeDBInstances` 接口查询实例状态
3. ✅ 等待 Status 值等于 4（已隔离）后再继续
4. ✅ 使用 Terraform 的 retry 机制处理等待逻辑
5. ✅ 确保代码格式化（执行 `go fmt`）

## Non-Goals

- 不修改 `TerminateSqlserverInstanceById` 和 `DeleteSqlserverInstanceById` 的实现
- 不修改 service 层的其他方法
- 不添加新的配置参数

## Success Criteria

1. 删除流程中正确等待实例隔离成功
2. 使用 `DescribeDBInstances` 接口正确查询实例状态
3. Status = 4 时继续执行删除操作
4. 代码通过 `go fmt` 格式化
5. 代码编译成功，无语法错误

## Dependencies

- 需要使用现有的 `SqlserverService.DescribeSqlserverInstances` 方法或创建新的查询方法
- 依赖 Terraform SDK 的 `resource.Retry` 机制
- 依赖腾讯云 SDK 的 `DescribeDBInstances` 接口

## Timeline

- 预计实施时间：1-2 小时
- 主要工作：修改 delete 函数，添加等待逻辑

## Risks and Mitigation

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 等待超时 | 删除失败 | 使用合理的超时时间（ReadRetryTimeout） |
| Status 值不正确 | 逻辑错误 | 参考腾讯云 API 文档确认 Status = 4 表示已隔离 |
| 兼容性问题 | 影响现有用户 | 仅添加等待逻辑，不改变删除行为 |

## References

- 腾讯云 SQL Server API 文档：DescribeDBInstances 接口
- 实例状态枚举：
  - 1：申请中
  - 2：运行中
  - 3：受限运行中
  - 4：已隔离 ✅ (目标状态)
  - 5：回收中
  - 6：已回收
  - 其他状态...
