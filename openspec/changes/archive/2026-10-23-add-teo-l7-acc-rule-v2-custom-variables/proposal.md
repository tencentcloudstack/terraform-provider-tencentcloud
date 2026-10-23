## Why

`tencentcloud_teo_l7_acc_rule_v2` 资源当前仅支持规则的基本字段（status、rule_name、description、branches 及其子分支 actions/sub_rules），但在边缘函数加速规则场景中，用户需要使用规则级自定义变量（CustomVariables）及其运算（CustomVariableOperations），以及通过 `CustomAction` 操作名调用定制配置（CustomActionParameters）。腾讯云 TEO SDK `teov20220901` 的 `RuleEngineItem` 已提供 `CustomVariables`、`CustomVariableOperations` 字段，`RuleEngineAction` 已提供 `CustomActionParameters` 字段，且 `CreateL7AccRules` / `ModifyL7AccRule` / `DescribeL7AccRules` 三个接口均支持这些字段的读写，但 Terraform Provider 尚未暴露这些能力，导致用户无法通过 Terraform 声明式管理自定义变量相关的高级规则配置。

## What Changes

- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源的顶层 schema 新增可选参数 `custom_variables`（列表类型），对应 `RuleEngineItem.CustomVariables`，每个元素包含 `name`、`initial_value`、`description` 字段。
- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源的顶层 schema 新增可选参数 `custom_variable_operations`（列表类型），对应 `RuleEngineItem.CustomVariableOperations`，每个元素包含 `branches`（含 `condition`、`actions`（含 `name`、`set_parameters`（含 `name`、`value`））、`sub_rules`（含 `branches`（递归）、`description`））和 `description` 字段。
- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源的 `branches.actions` schema 中新增可选参数 `custom_action_parameters`（列表类型，MaxItems=1），对应 `RuleEngineAction.CustomActionParameters`，包含 `custom_actions` 列表（每个元素含 `action_id` 和 `parameters`）。
- 扩展 Create 函数：将 `custom_variables`、`custom_variable_operations` 映射到 `CreateL7AccRulesRequest.Rules[0]` 的对应字段。
- 扩展 Update 函数：将 `custom_variables`、`custom_variable_operations` 映射到 `ModifyL7AccRuleRequest.Rule` 的对应字段，并在 `d.HasChange` 检测中加入这两个字段。
- 扩展 Read 函数：从 `DescribeL7AccRules` 返回的 `RuleEngineItem` 中读取 `CustomVariables`、`CustomVariableOperations` 回填到 state。
- 扩展 `resourceTencentCloudTeoL7AccRuleGetBranchs` / `resourceTencentCloudTeoL7AccRuleSetBranchs` 辅助函数以处理 actions 内的 `custom_action_parameters`。
- 更新资源文档 `resource_tc_teo_l7_acc_rule_v2.md`，补充新参数的示例与说明。

向后兼容：所有新增字段均为 Optional，未配置时资源行为与之前完全一致。

## Capabilities

### New Capabilities
<!-- 无新增 capability，本次仅在已有 capability 基础上新增参数 -->

### Modified Capabilities
- `teo-l7-acc-rule-v2-actions`: 新增对 `custom_variables`、`custom_variable_operations` 顶层参数以及 `branches.actions` 内 `custom_action_parameters` 参数的支持，覆盖 Create/Update/Read 三个接口的字段映射。

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`（顶层 schema 新增 custom_variables / custom_variable_operations，Create/Update/Read 逻辑扩展）
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go`（`TencentTeoL7RuleBranchBasicInfo` 的 actions 内新增 custom_action_parameters schema；GetBranchs/SetBranchs 辅助函数扩展 custom_action_parameters 的扁平化与还原逻辑；新增 custom_variables / custom_variable_operations 的辅助映射函数）
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`（新增覆盖 custom_variables / custom_variable_operations / custom_action_parameters 的单元测试用例）
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`（示例与字段说明同步）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.RuleEngineItem.CustomVariables`、`CustomVariableOperations`、`RuleEngineAction.CustomActionParameters`，无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由自动生成流程读取 `.md` 文件）。