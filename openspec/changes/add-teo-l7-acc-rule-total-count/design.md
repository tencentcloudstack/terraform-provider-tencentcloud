## Context

当前 tencentcloud_teo_l7_acc_rule 数据源通过 DescribeL7AccRules API 查询七层加速规则列表。用户在查询时可能需要知道总记录数，以便进行分页查询或数据统计。腾讯云 API 返回的响应中包含 TotalCount 字段，但目前数据源未将其导出。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 数据源的 schema 中新增 TotalCount 字段
- 从 DescribeL7AccRules API 响应中解析 TotalCount 并填充到 schema 中
- 更新相关文档以说明 TotalCount 字段的使用
- 添加或更新单元测试以验证 TotalCount 字段的正确性

**Non-Goals:**
- 不修改数据源的查询逻辑（除新增字段外）
- 不修改现有的其他字段
- 不改变数据源的行为和 API 调用方式

## Decisions

1. **字段类型选择**：使用 `schema.TypeInt` 作为 TotalCount 字段的类型，因为腾讯云 API 返回的 TotalCount 是整数类型
2. **字段属性**：设置为 `Computed`（只读），因为 TotalCount 是由 API 返回的计算值，不由用户输入
3. **字段映射**：在 schema 的 `Read` 函数中，将 API 响应中的 TotalCount 字段映射到 Terraform 状态
4. **文档更新**：在现有的 data_source_tc_teo_l7_acc_rule.md 文档中添加 TotalCount 字段的说明

## Risks / Trade-offs

**风险：**
- TotalCount 字段可能在某些情况下返回 null 或 0
  - **缓解措施**：使用 `d.Set("total_count", *response.TotalCount)` 时，先检查指针是否为 nil

**权衡：**
- 新增字段会增加数据源的状态大小
  - **影响**：影响较小，因为 TotalCount 只是单个整数值
