## Context

tencentcloud_teo_l7_acc_rule 资源是用于管理腾讯云 EdgeOne L7 访问规则的 Terraform Provider 资源。该资源当前通过 DescribeL7AccRules API 读取规则信息。API 响应中包含 TotalCount 字段（规则总数），但该字段目前未在 Terraform Provider 资源的 Schema 中暴露。

用户需要在查询资源时能够获取规则总数信息，以便更好地了解当前的规则规模。

当前约束：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 新增字段应为 Optional 和 Computed，不要求用户在配置中指定
- 仅从 API 响应中读取，不需要在 Create/Update/Delete 中处理

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 资源中新增 TotalCount 字段
- 确保 TotalCount 字段正确从 DescribeL7AccRules API 响应中读取
- 更新单元测试以验证新字段的读取逻辑
- 保持完全的向后兼容性

**Non-Goals:**
- 不修改 Create/Update/Delete 操作逻辑
- 不要求用户在 TF 配置中指定 TotalCount
- 不改变现有字段的行为

## Decisions

### 1. 字段类型和属性
**决策**: TotalCount 字段定义为 `schema.TypeInt`，并设置为 `Optional: true` 和 `Computed: true`

**理由**:
- `TypeInt` 匹配 API 返回的整数类型
- `Optional` 允许用户不指定该字段
- `Computed` 确保字段仅从 API 响应中填充，不会持久化到 state

**替代方案考虑**:
- 如果设置为 `Required`，会破坏现有配置（不可接受）
- 如果不设置 `Computed`，字段会被持久化到 state（不需要）

### 2. 字段读取位置
**决策**: 在 Read 函数中从 DescribeL7AccRules API 响应中读取 TotalCount

**理由**:
- TotalCount 是 API 响应中的字段
- Read 函数是读取资源状态的标准位置
- Create/Update/Delete 不需要处理该字段

### 3. 测试策略
**决策**: 在单元测试中添加对 TotalCount 字段的验证

**理由**:
- 确保字段正确读取
- 验证字段类型和值
- 保持测试覆盖率

## Risks / Trade-offs

**风险**: API 响应中 TotalCount 字段可能不存在或为 null
→ **缓解**: 在代码中添加空值检查，使用 `d.Set("total_count", 0)` 作为默认值

**风险**: 用户可能误解 TotalCount 字段的含义
→ **缓解**: 在文档中明确说明该字段为只读字段，表示规则总数

**权衡**: 新增字段会增加资源 schema 的复杂度
→ **缓解**: 该字段为只读且简单，不会显著增加复杂度，用户价值明确

## Migration Plan

### 部署步骤
1. 修改 tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go 文件，在 Schema 中添加 TotalCount 字段
2. 修改 Read 函数，从 API 响应中读取并设置 TotalCount 字段
3. 修改单元测试文件 resource_tc_teo_l7_acc_rule_test.go，添加 TotalCount 字段的测试用例
4. 运行单元测试验证修改
5. 更新文档（如需要）

### 回滚策略
由于该改动仅新增一个只读字段，回滚策略为：
1. 从 Schema 中移除 TotalCount 字段
2. 从 Read 函数中移除 TotalCount 字段的设置逻辑
3. 从测试中移除相关测试用例

回滚不会影响用户的 TF 配置和 state，因为该字段为 Computed，不会被持久化。

## Open Questions

无

