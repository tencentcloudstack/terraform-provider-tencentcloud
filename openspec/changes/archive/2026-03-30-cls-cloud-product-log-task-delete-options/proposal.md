# Proposal: Add Deletion Options for CLS Cloud Product Log Task

## What

为 `tencentcloud_cls_cloud_product_log_task_v2` 资源添加更精细的删除控制选项，支持在删除日志投递任务时选择性地删除关联的 Topic（日志主题）和 Logset（日志集）。

## Why

### 背景

目前 `tencentcloud_cls_cloud_product_log_task_v2` 资源只有一个 `force_delete` 字段，当设置为 `true` 时会强制删除关联的 Topic 和 Logset。这种"全部删除"的方式缺乏灵活性。

### 问题

1. **缺乏精细控制**：用户可能只想删除 Topic 而保留 Logset，或者只删除 Logset 而保留 Topic，但现有实现是"要么全删，要么都不删"。
2. **API 能力未充分利用**：腾讯云 DeleteCloudProductLogCollection API 已经提供了 `IsDeleteTopic` 和 `IsDeleteLogset` 两个独立的布尔参数，但 Terraform Provider 没有暴露这些能力。

### API 更新

DeleteCloudProductLogCollection 接口（文档：https://cloud.tencent.com/document/api/614/117420）新增了两个可选字段：

- **IsDeleteTopic** (Boolean): 是否删除关联的 Topic（日志主题）
- **IsDeleteLogset** (Boolean): 是否删除关联的 Logset（日志集）。如果 Logset 下还有其他 Topic，则不会删除

### 目标

1. 在 schema 中新增 `is_delete_topic` 和 `is_delete_logset` 字段，默认值为 `false`
2. 保持向后兼容：`force_delete` 为 `true` 时，自动设置 `is_delete_topic` 和 `is_delete_logset` 为 `true`
3. 添加逻辑约束：只有当 `force_delete` 为 `false` 时，用户才能自定义设置 `is_delete_topic` 或 `is_delete_logset`
4. 更新 Delete 模块，使用 API 的原生字段而非手动调用 DeleteTopic/DeleteLogset

## How

### Schema 变更

添加两个新的可选布尔字段：

```hcl
is_delete_topic = false   # 默认不删除 Topic
is_delete_logset = false  # 默认不删除 Logset
```

### 逻辑关系

- **force_delete = true**：忽略 `is_delete_topic` 和 `is_delete_logset` 的值，强制删除 Topic 和 Logset（向后兼容旧行为）
- **force_delete = false**：
  - 可以单独设置 `is_delete_topic = true` 来只删除 Topic
  - 可以单独设置 `is_delete_logset = true` 来只删除 Logset
  - 可以两者都设置为 `true` 来删除两者
  - 默认两者都是 `false`，不删除

### Delete 模块重构

**旧实现**（手动调用）：
```go
if force_delete {
    // 手动调用 DeleteTopic
    // 手动调用 DeleteLogset
}
```

**新实现**（使用 API 原生字段）：
```go
// 在 DeleteCloudProductLogCollection 请求中设置
if force_delete {
    request.IsDeleteTopic = helper.Bool(true)
    request.IsDeleteLogset = helper.Bool(true)
} else {
    if is_delete_topic {
        request.IsDeleteTopic = helper.Bool(true)
    }
    if is_delete_logset {
        request.IsDeleteLogset = helper.Bool(true)
    }
}
```

## Benefits

1. **更灵活的资源管理**：用户可以根据实际需求选择性删除资源
2. **更简洁的代码**：不再需要手动调用 DeleteTopic 和 DeleteLogset，API 内部处理
3. **更好的性能**：单次 API 调用完成所有删除操作，减少网络往返
4. **向后兼容**：现有使用 `force_delete` 的代码行为不变

## Impact

- **Breaking Changes**: 无
- **Deprecations**: 无（保留 `force_delete` 以保持向后兼容）
- **Migration**: 不需要，现有代码可以继续工作

## Success Criteria

1. 新增的 `is_delete_topic` 和 `is_delete_logset` 字段在 schema 中正确定义
2. Delete 模块使用 API 的原生字段而非手动删除
3. 逻辑约束正确实现：`force_delete` 优先级最高
4. 代码通过 `go fmt` 格式化
5. 编译通过，无语法错误
