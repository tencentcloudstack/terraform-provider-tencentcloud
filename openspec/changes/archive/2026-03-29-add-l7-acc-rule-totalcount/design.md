## Context

当前 DescribeL7AccRules API 的响应中包含 Rules 列表，但不包含 TotalCount 字段。SDK 中的 `DescribeL7AccRulesResponseParams` 结构体已经定义了 TotalCount 字段，但当前的实现中没有使用它。这导致用户在查询规则时无法获取总条数信息。

相关文件位于 `tencentcloud/services/teo/` 目录下：
- `service_tencentcloud_teo.go` - 包含 DescribeTeoL7AccRuleById 函数
- `resource_tc_teo_l7_acc_rule_v2.go` - 资源实现

## Goals / Non-Goals

**Goals:**
- 确保从 DescribeL7AccRules API 调用中返回 TotalCount 字段
- 如果存在数据源，更新其 schema 以包含 TotalCount 作为输出字段
- 确保相关测试覆盖 TotalCount 字段

**Non-Goals:**
- 不修改 API 的调用方式或请求参数
- 不改变现有的规则列表处理逻辑
- 不影响现有的资源和数据源功能

## Decisions

### 1. TotalCount 字段处理
- **决策**: 在 API 响应处理中保留 TotalCount 字段，但不将其暴露到资源的 schema 中
- **理由**: TotalCount 是一个用于统计和分页的字段，对单个资源实例没有意义。它主要用于数据源查询场景。
- **备选方案**: 也可以在资源中暴露该字段，但这会增加不必要的复杂性

### 2. 数据源支持（如果存在）
- **决策**: 如果存在 teo_l7_acc_rule 数据源，在其 schema 中添加 TotalCount 作为 Computed 字段
- **理由**: 数据源通常用于查询多个实例，TotalCount 字段对了解数据总量很有价值
- **实现**: 在 schema 定义中添加 `"total_count"` 字段，类型为 TypeInt，属性为 Computed

### 3. 代码修改范围
- **决策**: 修改 `service_tencentcloud_teo.go` 中的相关函数以正确处理 TotalCount
- **理由**: 这是处理 API 响应的中心位置，确保所有调用都能受益
- **备选方案**: 在调用处单独处理，但这会导致代码重复

## Risks / Trade-offs

### Risk 1: TotalCount 字段值可能不准确
- **风险**: 在某些边缘情况下，API 返回的 TotalCount 可能与实际 Rules 列表长度不一致
- **缓解**: 添加日志记录以便调试，同时确保错误处理逻辑不会因 TotalCount 问题而失败

### Risk 2: 向后兼容性
- **风险**: 修改现有结构体的使用方式可能引入意外行为
- **缓解**: TotalCount 是只读字段，不会影响现有的写入操作，因此风险很低

### Trade-off: 字段命名
- **权衡**: SDK 中的字段名为 `TotalCount` (驼峰命名)，但在 Terraform provider 中我们通常使用蛇形命名 `total_count`
- **决定**: 在 Terraform schema 中使用蛇形命名以保持一致性，在内部处理时进行转换

## Migration Plan

1. **开发阶段**:
   - 更新 `service_tencentcloud_teo.go` 中的相关函数
   - 如果存在数据源，更新其 schema
   - 编写或更新单元测试

2. **测试阶段**:
   - 运行验收测试（`TF_ACC=1`）
   - 验证 TotalCount 字段正确返回
   - 确保现有测试仍然通过

3. **部署阶段**:
   - 代码审查
   - 合并到主分支
   - 发布新版本

## Open Questions

无。这是一个相对简单的字段添加变更，技术实现路径清晰。
