## Context

tencentcloud_teo_l7_acc_rule 资源目前实现了 L7 加速规则的 CRUD 操作，通过 DescribeL7AccRules API 读取规则列表。API 响应中包含 TotalCount 字段（规则总数），但该字段未被映射到资源的 schema 中，导致用户无法通过 Terraform 查询规则总数。

当前实现中，Read 函数调用 DescribeTeoL7AccRuleById 方法（内部调用 DescribeL7AccRules API），返回的 DescribeL7AccRulesResponseParams 包含 TotalCount *int64 字段，但该字段被忽略，没有设置到 Terraform 状态中。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 资源 schema 中新增 `total_count` Computed 字段
- 在 Read 函数中从 API 响应读取 TotalCount 并设置到状态中
- 保持完全向后兼容，不影响现有配置和 state
- 添加相应的文档和测试

**Non-Goals:**
- 不修改 Create、Update、Delete 操作
- 不修改现有的 rules 列表字段的行为
- 不添加任何验证逻辑（TotalCount 仅作为只读信息）

## Decisions

1. **Schema 设计**：将 `total_count` 定义为 Computed 类型，TypeInt
   - **理由**：TotalCount 是 API 返回的只读信息，用户不能直接设置。设置为 Computed 让 Terraform 在 Read 时自动从 API 获取并更新状态。TypeInt 类型与 SDK 的 *int64 类型匹配。

2. **字段位置**：将 `total_count` 作为资源的顶级字段（与 zone_id、rules 同级）
   - **理由**：TotalCount 表示当前站点的规则总数，是资源级别的元数据，与 rules 列表同级更符合语义。如果放在 rules 列表的每个元素中会重复且不自然。

3. **Read 函数实现**：在 resourceTencentCloudTeoL7AccRuleRead 中，从 respData.TotalCount 读取值并使用 `d.Set("total_count", *respData.TotalCount)` 设置到状态
   - **理由**：这是 Terraform Provider 的标准模式。需要检查 TotalCount 不为 nil 避免空指针异常。

4. **文档更新**：在 resource_tc_teo_l7_acc_rule.md 中添加 `total_count` 字段的文档说明
   - **理由**：根据项目规范，所有资源字段必须有对应的文档。

5. **测试策略**：在现有的资源测试中添加对 TotalCount 字段的验证
   - **理由**：确保 TotalCount 字段正确返回，并与规则列表的长度匹配（可选，但能增强测试覆盖）。

## Risks / Trade-offs

- [API 不返回 TotalCount] → 如果 API 在某些情况下不返回 TotalCount（字段为 nil），代码需要处理这种情况。使用 if respData.TotalCount != nil 检查避免空指针。
- [与 rules 列表长度不一致] → TotalCount 应该等于 len(rules)，但在某些边缘情况下可能不一致（如 API 返回延迟）。这不是代码问题，而是 API 数据一致性保证。文档中说明 TotalCount 来自 API，不进行额外验证。
- [向后兼容性风险] → 极低风险，因为新增的是 Computed 字段，不影响现有配置。Terraform 会在下一次 apply 时自动读取新字段并更新状态。

## Migration Plan

由于这是一个纯新增 Computed 字段的变更，不需要特殊的迁移步骤：

1. 用户升级 Terraform Provider 到新版本
2. 下次运行 `terraform refresh` 或 `terraform apply` 时，新字段会自动从 API 读取并更新到状态中
3. 不需要手动修改 Terraform 配置文件

**回滚策略**：如果需要回滚，只需回退 Provider 版本即可。由于新字段是 Computed，回退不会导致状态不兼容。

## Open Questions

无。这是一个简单明确的字段新增变更，所有技术决策都很清晰。
