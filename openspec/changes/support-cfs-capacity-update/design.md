## Context

当前 `tencentcloud_cfs_file_system` 资源将 `capacity` 参数列为不可修改（immutableArgs），导致用户无法通过 Terraform 对 Turbo 系列文件系统进行扩容。腾讯云 CFS 服务提供了 `ScaleUpFileSystem` API 用于在线扩容 Turbo 文件系统，但 Provider 未暴露此能力。

**约束：**
- ScaleUpFileSystem 仅支持 Turbo 系列（storage_type 为 "TB" 或 "TP"）
- 仅支持扩容，不支持缩容
- 扩容步长固定：TB = 20,480 GiB, TP = 10,240 GiB
- 扩容操作是异步的，需要通过 DescribeCfsFileSystems 接口轮询 LifeCycleState 字段
  - `available`: 文件系统可用（扩容完成）
  - `expanding`: 扩容中

**关键文件：**
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go:312` - immutableArgs 列表
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go:303-367` - Update 函数
- `tencentcloud/services/cfs/service_tencentcloud_cfs.go` - 服务层 API 封装

## Goals / Non-Goals

**Goals:**
- 支持通过修改 `capacity` 参数对 Turbo 文件系统进行扩容
- 提供完整的扩容验证（容量增加、文件系统类型、扩容步长）
- 实现异步状态轮询，通过 DescribeCfsFileSystems 检查 LifeCycleState
- 支持 timeouts 块的 update 配置，控制扩容超时时间（默认 30 分钟）
- 保持向后兼容性，不影响现有资源行为

**Non-Goals:**
- 不支持缩容操作（云平台 API 不支持）
- 不支持标准型/高性能型文件系统的容量修改（仅 Turbo 支持）
- 不自动修正扩容步长（由用户确保容量符合步长要求）

## Decisions

### Decision 1: 验证逻辑位置
**决策：** 在 Update 函数中进行验证，而非 Schema 的 ValidateFunc。

**理由：**
- 验证需要访问当前状态值（旧容量、storage_type）
- ValidateFunc 仅能验证单个字段的值，无法进行跨字段验证
- 在 Update 函数中可以提供更详细的错误信息

**替代方案：**
- 使用 CustomizeDiff 进行验证 → 过于复杂，且无法访问远程状态
- 在 Schema 中使用 ValidateFunc → 无法获取当前容量和文件系统类型

### Decision 2: 扩容步长验证策略
**决策：** 验证目标容量是否为最小容量+扩容步长倍数的和。

**实现：**
```go
// Turbo 标准型 (TB)
minCapacity := 40960  // 40 TiB
increment := 20480    // 20 TiB
if (newCapacity - minCapacity) % increment != 0 {
    return error
}

// Turbo 性能型 (TP)
minCapacity := 20480  // 20 TiB
increment := 10240    // 10 TiB
if (newCapacity - minCapacity) % increment != 0 {
    return error
}
```

**理由：**
- 符合云平台 API 的容量要求
- 提供明确的错误提示，帮助用户纠正配置

### Decision 3: 异步状态轮询机制
**决策：** 调用 ScaleUpFileSystem 后，通过 DescribeCfsFileSystems 轮询 LifeCycleState 字段。

**实现：**
```go
// 调用扩容 API
err := cfsService.ScaleUpFileSystem(ctx, fsId, newCapacity)
if err != nil {
    return err
}

// 异步等待扩容完成
updateTimeout := d.Timeout(schema.TimeoutUpdate)
err = resource.Retry(updateTimeout, func() *resource.RetryError {
    instance, err := cfsService.DescribeFileSystemById(ctx, fsId)
    if err != nil {
        return resource.NonRetryableError(err)
    }
    
    state := *instance.LifeCycleState
    if state == "available" {
        return nil  // 扩容完成
    }
    if state == "expanding" {
        return resource.RetryableError(fmt.Errorf("waiting for expansion, current state: %s", state))
    }
    return resource.NonRetryableError(fmt.Errorf("unexpected state: %s", state))
})
```

**理由：**
- LifeCycleState 字段是腾讯云 CFS 的官方状态指示器
- `expanding` 状态明确表示扩容操作正在进行
- `available` 状态确认扩容已完成且文件系统可用
- 使用轮询机制符合 Terraform Provider 的异步操作最佳实践

### Decision 4: Update Timeout 配置
**决策：** Update timeout 默认值设为 30 分钟，支持用户通过 timeouts 块自定义。

**理由：**
- 扩容操作通常在 10-20 分钟内完成
- 30 分钟提供足够的缓冲，避免偶发的慢速操作超时
- 与 Create timeout 的 20 分钟保持一致的数量级
- 用户可根据文件系统规模自定义超时时间

**替代方案：**
- 60 分钟 → 过长，会延迟错误发现
- 20 分钟 → 可能不够，扩容大容量时可能超时

### Decision 5: API 错误处理
**决策：** 对 ScaleUpFileSystem API 调用使用标准重试逻辑（`resource.Retry` + `tccommon.WriteRetryTimeout`）。

**理由：**
- 与项目中其他资源的 Update 操作保持一致
- 处理瞬态网络错误和 API 限流
- 非重试错误（如参数错误）会立即返回

## Risks / Trade-offs

### Risk 1: 扩容失败导致状态不一致
**风险：** ScaleUpFileSystem API 调用成功，但等待 LifeCycleState 变为 available 超时。

**缓解措施：**
- 使用足够长的 timeout（30 分钟默认，用户可通过 timeouts.update 自定义）
- 在错误信息中提示用户手动检查云端资源状态
- Read 函数会从云端同步实际容量和状态，用户可以 `terraform refresh` 更新状态
- 下次 apply 会基于云端实际状态继续操作

### Risk 2: 用户误操作缩容
**风险：** 用户意外将容量改小，期望 Terraform 执行缩容。

**缓解措施：**
- 在 Update 函数中明确拒绝缩容操作
- 错误信息清楚说明"仅支持扩容"
- 文档中突出说明此限制

### Risk 3: 非 Turbo 系列修改 capacity
**风险：** 用户对标准型/高性能型文件系统修改 capacity，期望生效。

**缓解措施：**
- 在 Update 函数中检查 storage_type
- 返回明确错误："仅 Turbo 系列（TB/TP）支持容量修改"
- 文档中明确说明此限制

## Migration Plan

**对现有资源的影响：**
1. **已创建的资源**：
   - 不修改 capacity → 无影响
   - 修改 capacity → 从"触发重建"变为"原地更新"

2. **用户迁移步骤**：
   - 无需特殊迁移步骤
   - 用户首次修改 capacity 时会触发扩容而非重建
   - 建议在变更日志和文档中突出说明此行为变更

3. **回滚计划**：
   - 如发现严重问题，可将 `capacity` 重新加入 immutableArgs
   - 已扩容的资源不受影响（只是后续修改会触发重建）

## Open Questions

1. **Q**: 是否需要支持自动修正扩容步长（如用户输入 50000，自动向上取整到 61440）？
   **A**: 不支持。要求用户明确指定正确的容量值，避免意外的容量分配。

2. **Q**: 是否需要在 Plan 阶段显示扩容警告？
   **A**: 不需要。Terraform 的 diff 输出已经清楚显示容量变更，足以让用户理解操作。

3. **Q**: 扩容失败后是否需要回滚？
   **A**: 不需要。云平台不支持回滚，且文件系统已存储数据，强制回滚会导致数据丢失风险。
