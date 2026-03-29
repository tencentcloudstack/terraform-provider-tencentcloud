## Context

当前 `tencentcloud_teo_l7_acc_rule` 数据源的读取操作（`Read` 方法）通过调用 `DescribeL7AccRules` API 获取七层访问规则列表。该 API 返回的数据中包含 `TotalCount` 字段，表示符合查询条件的规则总数。但当前实现中未接入此参数，导致用户无法获取此重要信息。

本变更需要在现有数据源实现中添加对 `TotalCount` 字段的处理逻辑。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule` 数据源的 schema 中添加 `TotalCount` 字段（Computed 类型）
- 在 `Read` 方法中解析 API 响应并设置 `TotalCount` 值
- 更新相关测试用例以验证 `TotalCount` 字段的正确性

**Non-Goals:**
- 不修改 API 调用本身的逻辑
- 不修改其他字段的处理方式
- 不涉及状态存储的变化

## Decisions

### 1. Schema 设计
在数据源 schema 中添加 `TotalCount` 字段：
```go
"total_count": {
    Type:     schema.TypeInt,
    Computed: true,
    Description: "符合条件的规则总数",
}
```

**理由：**
- `TotalCount` 是 API 响应中的只读字段，设置为 `Computed` 类型符合语义
- 不需要在 schema 中标记为 `Required` 或 `Optional`，因为它完全由 API 响应填充

### 2. 数据处理位置
在 `Read` 方法的响应解析部分添加 `TotalCount` 字段的处理逻辑，位置应在解析规则列表之后。

**理由：**
- 与其他响应字段的处理保持一致的代码结构
- 便于维护和后续扩展

### 3. 测试策略
在单元测试中添加对 `TotalCount` 字段的验证，确保：
- API 响应中的 `TotalCount` 值能正确设置到 schema 中
- 值为 0 的情况也能正确处理

**理由：**
- 保证代码质量和正确性
- 防止回归问题

## Risks / Trade-offs

**Risk 1：** API 响应结构变化可能影响字段解析
**Mitigation：** 使用类型安全的方式解析 JSON 响应，避免直接字段访问。参考现有的解析逻辑。

**Risk 2：** 现有用户如果已在 schema 中自定义了 `total_count` 字段
**Mitigation：** 检查现有 schema 确认无冲突。由于 `Computed` 字段是只读的，不会影响用户写入的配置。

**Trade-off：** 添加新字段会增加 schema 的复杂度，但带来的价值（用户能够获取总数信息）远大于增加的复杂度。

## Migration Plan

1. 修改数据源 schema 添加 `TotalCount` 字段
2. 修改 `Read` 方法添加响应解析逻辑
3. 更新单元测试
4. 验证集成测试（如果存在）

**Rollback Strategy：**
- 如果出现问题，可以回滚 schema 修改和代码变更
- 由于只是添加了一个 `Computed` 字段，不会破坏现有的 Terraform 配置或 state

## Open Questions

无（本变更相对简单，无未决问题）
