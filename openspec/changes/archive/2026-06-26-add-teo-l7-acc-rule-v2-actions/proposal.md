## Why

为 `tencentcloud_teo_l7_acc_rule_v2` 资源新增顶层 `actions` 参数。当前资源在 `branches` 子 schema 中已包含 `actions` 字段（对应 SDK `RuleBranch.Actions`），用户需要通过嵌套层级访问。本变更将 `actions` 提升为资源顶层参数，简化 Terraform 配置，使得用户可以直接在 `tencentcloud_teo_l7_acc_rule_v2` 资源根级别配置加速规则的操作列表（对应 SDK `RuleEngineItem.Branches[0].Actions`）。

## What Changes

- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源 Schema 顶层新增 **Optional** 参数 `actions`（类型：TypeList，元素类型为 `RuleEngineAction` 子 schema），对应 SDK `CreateL7AccRulesRequest.Rules[0].Branches[0].Actions` / `ModifyL7AccRuleRequest.Rule.Branches[0].Actions` / `DescribeL7AccRulesResponse.Rules[0].Branches[0].Actions` 的 `[]*RuleEngineAction` 字段
- 在 Create / Read / Update 方法中添加 `actions` 参数的读写逻辑
- 新增参数为 Optional，不影响已有配置的向后兼容性

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-actions`: 在 `tencentcloud_teo_l7_acc_rule_v2` 资源顶层新增 `actions` 参数，允许用户直接在资源根级别定义加速规则操作列表（L7 加速规则 actions）

### Modified Capabilities
（无）

## Impact

- 受影响文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`（新增 Schema 定义及 CRUD 逻辑）
- 受影响文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`（更新文档示例）
- 不受影响：vendor SDK 无需更新（`RuleBranch.Actions` 字段已存在于 SDK 中）
- 向后兼容：`actions` 为 Optional 参数，不填写时行为不变
