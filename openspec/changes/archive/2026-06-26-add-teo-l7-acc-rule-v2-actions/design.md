## Context

当前 `tencentcloud_teo_l7_acc_rule_v2` 资源（文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`）通过 `TencentTeoL7RuleBranchBasicInfo(1)` 在 `branches` 子 schema 中已包含 `actions` 字段。该字段对应 SDK `RuleBranch.Actions`（类型：`[]*RuleEngineAction`），用于定义加速规则的操作列表（如 Cache、CacheKey、AccessURLRedirect 等）。

SDK 中四个与 L7 加速规则相关的 API 为：
- `CreateL7AccRulesRequest`: 包含 `Rules []*RuleEngineItem`，`RuleEngineItem.Branches` 下可包含 `Actions`
- `ModifyL7AccRuleRequest`: 包含 `Rule *RuleEngineItem`，同上路径
- `DescribeL7AccRulesResponse`: 返回 `Rules []*RuleEngineItem`，同上路径读出
- `DeleteL7AccRulesRequest`: 仅需 `ZoneId` + `RuleIds`，不涉及 Actions

当前资源通过在 `branches` 嵌套中访问 `actions`，用户配置时需要多一层 `branches {}` 包裹。本设计将 `actions` 提升为资源顶级 Optional 参数，简化配置。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源 Schema 顶层新增 `actions` 参数，类型为 TypeList，元素为 `RuleEngineAction` 子 schema（与 `TencentTeoL7RuleBranchBasicInfo` 中的 `actions` 定义一致）
- Create 时：当 `actions` 有值时，将其映射到 SDK `RuleEngineItem.Branches[0].Actions`
- Read 时：从 API 响应 `Rules[0].Branches[0].Actions` 中提取值并设置到 terraform state 的顶层 `actions`
- Update 时：当 `actions` 有变更时，将其映射到 SDK `RuleEngineItem.Branches[0].Actions` 进行更新
- 保持向后兼容：`actions` 为 Optional，现有配置不受影响

**Non-Goals:**
- 不修改 `branches` 子 schema 中已有的 `actions` 实现
- 不修改 `DeleteL7AccRules` 接口调用（该接口不涉及 Actions）
- 不修改 `resource_tc_teo_l7_acc_rule_extension.go` 中的共享辅助函数

## Decisions

### 1. `actions` Schema 定义复用已有 `RuleEngineAction` 结构
**选择**：直接复用 `TencentTeoL7RuleBranchBasicInfo` 中 `actions` 的 Schema 定义（包括 `name`、`cache_parameters`、`cache_key_parameters` 等所有子字段）。

**理由**：避免重复定义，确保顶层 `actions` 与 `branches.actions` 行为一致。

**替代方案**：重新定义简化的 Schema → 拒绝，因为会导致两层 actions 行为不一致，增加维护负担。

### 2. 当 `actions` 与 `branches` 同时配置时的处理
**选择**：`actions` 和 `branches` 独立配置，Create/Update 时合并处理 —— 若 `actions` 有值且 `branches` 未设置，则创建一个仅包含 Actions 的 `RuleBranch`；若两者都设置，保持不变（branches 优先，actions 作为补充）。

**理由**：保持灵活性，`actions` 作为简化入口，`branches` 作为高级入口。

### 3. Read 时 actions 的提取路径
**选择**：从 `DescribeL7AccRules` 返回的 `Rules[0].Branches[0].Actions` 提取值设置到顶层 `actions`。

**理由**：SDK 中 `Actions` 位于 `RuleBranch` 层级，RO 路径必须与 API 结构对应。

### 4. 不修改 Delete 逻辑
**选择**：Delete 接口仅需 `ZoneId` + `RuleIds`，不涉及 `Actions` 字段，无需变更。

**理由**：Delete 操作的 SDK 请求不包含 Actions 字段，无需处理。

## Risks / Trade-offs

- **[风险] `actions` 与 `branches.actions` 语义重叠**：顶层 `actions` 和 `branches[0].actions` 在 SDK 层面映射到同一路径（`Branches[0].Actions`），可能导致用户困惑。
  - **缓解**：在文档中明确说明：`actions` 是简化用法，`branches` 是高级用法，建议只使用其中一种。

- **[风险] Read 后 state 漂移**：若只配了 `actions`，Read 返回后可能由于 API 响应的 Branches 结构变化导致差异。
  - **缓解**：在 Read 中确保如果响应 `Branches` 存在且用户使用了顶层 `actions`，则正确提取并设置。

## Migration Plan

1. 新增 `actions` 参数到资源 Schema（Optional，无默认值）
2. 在 Create、Read、Update 方法中添加 actions 处理逻辑
3. 更新 `.md` 文档添加 Example Usage 展示 actions 用法
4. 运行现有单元测试确保无回归
