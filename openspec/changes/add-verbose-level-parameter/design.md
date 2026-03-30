## Context

当前 `tencentcloud_nats` 数据源在调用 TencentCloud API `DescribeNatGateways` 时未传递任何控制输出详细程度的参数。用户无法根据实际需求获取不同粒度的数据，导致在某些场景下获取过多不必要的信息，影响查询效率和性能。

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_nats` 数据源添加 `VerboseLevel` 参数，允许用户控制 API 返回数据的详细程度
- 确保参数正确传递到 `DescribeNatGateways` API 调用中
- 保持向后兼容性，不影响现有用户的使用
- 更新相关文档和测试

**Non-Goals:**
- 修改现有数据源的其他参数或行为
- 修改腾讯云 API 的实现
- 实现数据过滤逻辑（API 端处理）

## Decisions

### 参数类型选择
**决策**：`VerboseLevel` 参数定义为整型（Int），可选参数，默认值不设置

**理由**：
- DescribeNatGateways API 接受整型的 VerboseLevel 参数
- 作为可选参数，不设置时使用 API 默认行为，保持向后兼容
- 整型类型与 API 原始参数类型一致，无需类型转换

### 参数传递方式
**决策**：在 `data_source_tc_nats.go` 的 schema 中添加参数，在查询函数中读取并传递给 API 调用

**理由**：
- 符合 Terraform Provider 的标准模式
- 参数直接传递给 API，不在 Provider 端进行额外处理
- 简单直接，易于维护

### 文档更新策略
**决策**：在数据源文档中添加参数说明，包括类型、可选值和用途

**理由**：
- 用户需要了解如何使用新参数
- 遵循项目文档规范
- 提供清晰的使用示例

## Risks / Trade-offs

### 风险 1：参数值验证不足
**风险**：如果用户传入无效的 VerboseLevel 值，API 可能返回错误或不预期的结果
**缓解**：依赖 API 端的参数验证，不在 Provider 端添加额外的验证逻辑，保持简单

### 风险 2：向后兼容性
**风险**：新参数可能影响现有配置（如果参数名称与其他参数冲突）
**缓解**：参数名称 `VerboseLevel` 与 API 参数名一致，不与现有 schema 字段冲突

### 权衡：简单性 vs 灵活性
**选择优先考虑简单性**：直接传递参数给 API，不在 Provider 端进行复杂的参数处理或验证
**影响**：减少了代码复杂度，但用户需要参考 API 文档了解有效参数值

## Migration Plan

无需迁移计划：
- 新增参数为可选参数，不影响现有配置
- 用户可选择性地使用新参数，无需修改现有代码
- 部署后立即生效，无需数据迁移

## Open Questions

（无，所有技术决策已明确）
