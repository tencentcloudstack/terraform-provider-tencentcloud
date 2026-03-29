## Context

tencentcloud_teo_l7_acc_rule 是一个用于查询边缘访问规则的数据源。当前数据源从 DescribeL7AccRules API 获取规则列表，但未利用 API 响应中的 TotalCount 字段。该字段提供了规则总数，对于需要了解规则数量的场景非常有用（例如：分页处理、资源配额检查等）。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 数据源的 schema 中添加 TotalCount 字段
- 从 DescribeL7AccRules API 响应中读取 TotalCount 值并正确映射到 schema
- 确保向后兼容，不影响现有用户的使用

**Non-Goals:**
- 不修改 TEOS 服务的基础架构
- 不改变现有 API 调用逻辑
- 不涉及资源的创建、更新、删除操作

## Decisions

1. **字段类型**: 使用 TypeInt 类型，因为 API 返回的 TotalCount 是整数类型

2. **字段位置**: 将 TotalCount 添加到 schema 的顶层，作为数据源的输出字段（Computed 属性）

3. **映射逻辑**: 在 Read 函数中，从 API 响应的 TotalCount 字段直接读取并设置到 state 中

4. **依赖管理**: 不需要新增外部依赖，利用现有的 tencentcloud-sdk-go 中的 teo 包

## Risks / Trade-offs

**Risks:**
- TotalCount 字段可能在某些 API 响应中不存在或为 null → 在代码中添加 nil 检查，设置默认值 0

**Trade-offs:**
- 使用 Computed 而非 Optional，因为该字段由 API 返回，用户无法输入 → 这是正确的选择，符合数据源模式

## Migration Plan

无迁移需求，这是一个新增字段的变更，向后兼容。

## Open Questions

无
