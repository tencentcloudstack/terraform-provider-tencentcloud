## Context

当前 `tencentcloud_nats` 数据源通过调用 DescribeNatGateways API 查询 NAT 网关信息。API 支持一个可选的 `VerboseLevel` 参数，用于控制返回数据的详细程度。该参数可以优化查询性能，特别是在用户只需要基本信息而不需要完整的 NAT 规则和自定义路由信息时。

当前实现未暴露此参数，用户无法根据实际需求控制返回数据的详细程度。在某些场景下（如大规模批量查询），返回不必要的数据会增加网络传输开销和解析时间。

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_nats 数据源添加 verbose_level 参数支持
- 允许用户选择返回数据的详细程度（DETAIL/COMPACT/SIMPLE）
- 保持向后兼容性，不破坏现有配置
- 提供清晰的文档说明不同详细级别的区别

**Non-Goals:**
- 不修改 API 的默认行为
- 不影响其他数据源或资源
- 不改变现有参数的行为

## Decisions

**1. 参数类型和位置**
- 在 schema 的顶层新增 `verbose_level` 参数
- 类型为 `TypeString`，属性为 `Optional`
- 使用 `ValidateFunc` 限制有效值为 "DETAIL"、"COMPACT"、"SIMPLE"

**理由**：
- 该参数与 vpc_id、name、state 等筛选参数处于同一层级，都是控制查询行为的参数
- 使用 String 类型与 API 的 VerboseLevel 字段类型一致
- 添加验证函数可以提前捕获无效值，提供更好的用户体验

**2. 参数验证**
- 创建自定义验证函数 `validateVerboseLevel`
- 仅允许三个有效值：DETAIL、COMPACT、SIMPLE
- 错误信息提示用户可选值

**理由**：
- API 只支持这三个值，提前验证可以避免无效请求
- 与 provider 中其他类似的参数验证保持一致（如 ValidateStringLengthInRange）

**3. API 调用方式**
- 在 `dataSourceTencentCloudNatsRead` 函数中，检查 verbose_level 参数
- 如果设置了该参数，将其赋值给 `request.VerboseLevel`
- 如果未设置，保持默认行为（不传该参数）

**理由**：
- API 的 VerboseLevel 是可选参数，不传时使用默认行为
- 保持向后兼容，现有配置不受影响

**4. 文档更新**
- 在参数说明中添加 verbose_level 的文档
- 说明三种详细级别的区别和适用场景
- 提供使用示例

**理由**：
- 用户需要了解不同详细级别的差异才能正确使用
- 示例可以帮助用户快速上手

## Risks / Trade-offs

**风险 1**: 用户可能会选择错误的详细级别，导致缺少所需信息
- **缓解措施**: 在文档中清晰说明每个级别的数据内容，建议默认使用 DETAIL

**风险 2**: API 可能会更改 VerboseLevel 的支持值
- **缓解措施**: 使用 ValidateFunc 时允许 future-proof 设计，但如果 API 变更需要更新 provider

**权衡**: 添加参数增加了 schema 的复杂度，但提供了更好的性能优化能力
- **决策**: 性能优化的收益大于参数增加的复杂度，值得添加

## Open Questions

无
