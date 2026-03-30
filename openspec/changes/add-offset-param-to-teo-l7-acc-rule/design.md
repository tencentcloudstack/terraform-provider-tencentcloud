## Context

当前 tencentcloud_teo_l7_acc_rule 数据源通过 DescribeL7AccRules API 查询七层访问规则列表。该 API 已经支持 Offset 参数来实现分页查询，但数据源未暴露此参数，导致用户无法控制查询偏移量，只能查询从第一条开始的数据。

现有数据源架构：
- 文件位置: `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go`
- SDK 调用: 通过 tencentcloud-sdk-go 的 teo 包调用 DescribeL7AccRules API
- 当前支持的参数: ZoneId, Type, Protocol, RuleName 等，但缺少 Offset

技术约束：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- Offset 参数为 Optional，不影响现有使用场景
- 遵循 Terraform Plugin SDK v2 的 schema 定义规范

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 数据源 schema 中添加 Offset 参数（Optional, Int）
- 将 Offset 参数传递给 DescribeL7AccRules API，实现分页查询功能
- 添加 Offset 参数验证，确保其为非负整数
- 更新数据源文档，说明 Offset 参数的使用方法
- 添加测试用例验证 Offset 参数功能

**Non-Goals:**
- 不修改现有数据源的任何其他参数或行为
- 不实现 Limit 参数（如果需要可后续添加）
- 不改变数据源的返回数据结构
- 不涉及其他 TEO 数据源的修改

## Decisions

### Schema 定义方式
**决策**: 使用 `schema.Int` 类型定义 Offset 参数，设置为 Optional，不设置默认值

**理由**:
- Optional 参数不会影响现有配置，保持向后兼容
- 不设置默认值可以让 API 使用其默认行为（通常为 0）
- Type: Int 符合 API 参数要求，且是 Terraform 常用类型

**替代方案考虑**:
- 设置 Default: 0 - 会导致 state 文件中始终包含此字段，增加不必要的数据
- 使用 Float 类型 - API 不需要浮点数，类型不匹配

### 参数传递策略
**决策**: 在 datasource 的 Read 函数中，从 d.Get("offset") 获取值，有值时添加到 API 请求参数中

**理由**:
- 简单直接，遵循现有代码模式
- 只有当用户明确指定 Offset 时才传递给 API，避免不必要的参数传递
- 与现有参数（如 Type, Protocol）的处理方式保持一致

### 参数验证实现
**决策**: 使用 schema 的 ValidateDiagFunc 对 Offset 参数进行验证

**理由**:
- 在 schema 层面进行验证，可以在用户配置阶段就发现错误
- 提供清晰的错误提示，提升用户体验
- 符合 Terraform Provider 的最佳实践

**替代方案考虑**:
- 在 Read 函数中验证 - 会在运行时才报错，不够及时

## Risks / Trade-offs

### Risk: API 不支持 Offset 参数
**风险**: 如果旧版本的 SDK 不支持 Offset 参数，可能导致 API 调用失败

**缓解措施**:
- 在 proposal 中已明确依赖 TencentCloud Go SDK v1.0.831+
- 在测试中验证 Offset 参数功能
- 检查 vendor 目录中的 SDK 版本是否满足要求

### Risk: 大 Offset 值导致性能问题
**风险**: 用户设置过大的 Offset 值可能导致 API 响应时间过长

**缓解措施**:
- API 层面通常有最大 Offset 值限制，由后端控制
- 在文档中建议合理使用 Offset，通常配合 Limit 参数使用（虽然 Limit 不在此次变更范围）

### Trade-off: 不实现 Limit 参数
**权衡**: 本次变更只添加 Offset 参数，不添加 Limit 参数

**原因**:
- 保持变更范围最小化，降低风险
- Offset 和 Limit 通常一起使用，但可以分阶段实现
- 如果后续需要 Limit，可以作为独立的变更

## Migration Plan

此次变更不需要迁移策略：
- 新增 Optional 参数，不影响现有配置和 state
- 用户可以自然地开始使用 Offset 参数
- 无需更新现有 Terraform 代码或 state 文件

部署步骤：
1. 修改数据源代码，添加 Offset 参数
2. 更新文档
3. 添加测试用例
4. 运行测试确保功能正常
5. 合并到主分支并发布新版本

回滚策略：
- 如果发现问题，可以简单地移除 Offset 参数的代码
- 由于参数是 Optional 的，回滚不会影响现有用户
- state 文件中未使用 Offset 参数的配置不受影响

## Open Questions

无。此次变更范围明确，实现方式清晰，不需要进一步的技术决策。
