## Context

当前 `tencentcloud_teo_l7_acc_rule` 资源仅返回 L7 访问控制规则的详细信息列表（rules），但不提供规则总数的统计信息。DescribeL7AccRules API 响应中已经包含 `TotalCount` 字段，表示符合条件的规则总数，但在当前资源实现中未将其映射到 Terraform schema 中。

资源的主要调用路径是：
1. 资源的 Read 方法调用 `service.DescribeTeoL7AccRuleById(ctx, zoneId, "")`
2. `DescribeTeoL7AccRuleById` 函数调用腾讯云 Teo V20220901 服务的 `DescribeL7AccRules` API
3. API 返回 `DescribeL7AccRulesResponseParams`，其中包含 `TotalCount` 字段

这是一个相对简单的字段映射变更，不涉及架构调整或数据模型改变。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule` 资源的 schema 中添加 `total_count` Computed 字段
- 在 Read 方法中将 API 响应的 `TotalCount` 值设置到该字段
- 保持向后兼容，不破坏现有配置和状态
- 确保新增字段能够正确反映当前站点下 L7 访问控制规则的总数

**Non-Goals:**
- 不修改资源的 Create、Update、Delete 方法
- 不添加任何新的 API 调用
- 不修改现有的 rules 列表结构
- 不添加任何新的依赖或 SDK 引用

## Decisions

1. **字段类型选择**: 使用 `TypeInt` 作为 `total_count` 字段的类型
   - **理由**: 腾讯云 SDK 的 `TotalCount` 字段为 `*int64` 类型，对应 Terraform 的 `TypeInt`
   - **替代方案**: 使用 `TypeString` 可能会引入不必要的类型转换，且不符合 API 的原始类型

2. **字段属性设置**: 设置为 `Computed` 和 `Optional: false`
   - **理由**: `total_count` 仅用于输出，由 API 返回值决定，不允许用户手动设置
   - **向后兼容**: Computed 字段不会破坏现有配置，Terraform 会自动忽略配置文件中未定义的字段

3. **字段位置**: 将 `total_count` 添加在 schema 的顶层，与 `zone_id` 和 `rules` 平级
   - **理由**: 该字段描述的是整个规则集合的统计信息，而不是单个 rule 的属性
   - **替代方案**: 将其添加到 `rules` 块内的每个规则中是不合适的，因为它是总数而不是单个规则的属性

4. **错误处理**: 不对 `TotalCount` 为 nil 的情况做特殊处理
   - **理由**: 当 API 返回 `TotalCount = nil` 时，Terraform 会自然地将该字段值设为 0（int 类型的零值），这是合理的行为
   - **替代方案**: 添加显式的 nil 检查会增加代码复杂度，但没有额外收益

## Risks / Trade-offs

- **[API 字段变更风险]** 如果腾讯云 API 在未来版本中移除或修改 `TotalCount` 字段，可能导致此字段失效
  - **缓解措施**: 这是一个只读的 Computed 字段，即使 API 字段变化也不会破坏用户的 Terraform 配置或状态文件。需要通过单元测试和集成测试验证 API 兼容性

- **[Terraform 状态迁移风险]** 虽然这是新增字段，但理论上不应该影响现有状态
  - **缓解措施**: Computed 字段不会破坏现有状态，Terraform 在读取状态时会自动添加该字段。通过运行 acceptance test 验证状态兼容性

- **[语义歧义风险]** `total_count` 可能被误解为当前返回的 `rules` 列表长度，而不是 API 查询到的总数
  - **缓解措施**: 在字段的 Description 中明确说明这是 API 返回的总数，并提供清晰的文档说明

## Migration Plan

这是一个无破坏性的变更，不需要特殊的迁移步骤：

1. 代码变更完成后，直接发布新版本的 Terraform Provider
2. 用户在升级 Provider 后，下次运行 `terraform plan` 或 `terraform apply` 时，Terraform 会自动检测到 `total_count` 是一个新增的 Computed 字段
3. 现有的 Terraform 配置文件无需修改
4. 用户可以在 `terraform output` 或 `terraform state show` 中看到新增的 `total_count` 字段

Rollback 策略：
- 如果发现 `total_count` 字段存在问题，可以通过降级 Provider 版本来回退
- 由于这是只读字段，回退不会影响用户的状态文件

## Open Questions

None. 此变更范围明确，依赖现有的 API 字段，不存在未解决的技术决策。
