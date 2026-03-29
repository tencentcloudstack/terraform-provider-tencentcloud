## Context

tencentcloud_teo_l7_acc_rule 是一个数据源，用于查询腾讯云边缘安全加速（TEO）服务的七层加速规则。目前该数据源从 DescribeL7AccRules API 读取规则列表，但未暴露 TotalCount 字段，导致用户无法获取规则总数信息。

该数据源位于 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go`。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 数据源中新增 TotalCount 输出字段
- 从 DescribeL7AccRules API 响应中正确读取并映射 TotalCount
- 保持向后兼容性，不影响现有功能
- 更新相应的文档和测试

**Non-Goals:**
- 不修改数据源的查询过滤逻辑
- 不改变现有字段的输出格式
- 不涉及 API 调用参数的变更

## Decisions

### Schema 扩展
在数据源 schema 中新增 TotalCount 字段，类型为 TypeInt，设置为 Computed，无需用户输入。

```go
"total_count": {
    Type:     schema.TypeInt,
    Computed: true,
},
```

### 数据映射
在 ReadResource 函数中，从 API 响应中读取 TotalCount 字段并映射到 schema：

```go
if response.Response.TotalCount != nil {
    d.Set("total_count", *response.Response.TotalCount)
}
```

### 文档更新
在 `website/docs/r/teo_l7_acc_rule.html.markdown` 中添加 TotalCount 字段说明：
- 在输出属性文档中添加 `(String, Computed)` TotalCount: 七层加速规则总数

### 测试更新
在 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go` 中添加验证：
- 验证 TotalCount 字段在响应中正确设置
- 确保现有测试用例仍然通过

## Risks / Trade-offs

**风险：API 响应中 TotalCount 可能为 nil**
- 缓解措施：使用 nil 检查确保安全，只有当字段非 nil 时才设置到 state 中

**风险：字段名不一致**
- 缓解措施：使用 SDK 中定义的字段名 TotalCount，确保与 API 响应一致

## Migration Plan

本变更为向后兼容的新增功能，无需迁移步骤：
- 现有配置无需修改即可继续工作
- 新字段自动可用，用户可以选择性地使用

## Open Questions

（无）
