## 1. Schema 定义

- [x] 1.1 在 `resource_tc_teo_l7_acc_rule_v2.go` 的 `Schema` map 中新增顶层参数 `custom_variables`（TypeList, Optional），元素为 Resource schema，包含 `name`、`initial_value`、`description`（均 TypeString, Optional），映射 `CustomVariable` 结构体
- [x] 1.2 在 `resource_tc_teo_l7_acc_rule_v2.go` 的 `Schema` map 中新增顶层参数 `custom_variable_operations`（TypeList, Optional, MaxItems: 1），元素包含 `branches`（TypeList，单个 `CustomVariableOperationRuleBranch`：`condition`/`actions`/`sub_rules`）和 `description`（TypeList of TypeString）
- [x] 1.3 在 `custom_variable_operations.branches.actions` 元素 schema 中定义 `name`（TypeString, Optional）和 `set_parameters`（TypeList, MaxItems: 1，含 `name`/`value`，映射 `SetParameters`）
- [x] 1.4 在 `custom_variable_operations.branches.sub_rules` 元素 schema 中定义 `branches`（TypeList of `CustomVariableOperationRuleBranch`，含 `condition`/`actions`/`sub_rules`，单层不再递归）和 `description`（TypeList of TypeString）
- [x] 1.5 在 `resource_tc_teo_l7_acc_rule_extension.go` 的 `TencentTeoL7RuleBranchBasicInfo` actions 块中新增 `custom_action_parameters`（TypeList, Optional, MaxItems: 1），包含 `custom_actions`（TypeList of `CustomAction`：`action_id` TypeString + `parameters` TypeList of `CustomActionParameter`，每个 parameter 含 `name` 和 `value` TypeString）
- [x] 1.6 更新 actions 的 `name` 字段 description，追加 `CustomAction: Custom configuration` 取值说明

## 2. Create 函数（resource_tc_teo_l7_acc_rule_v2.go）

- [x] 2.1 在 `ResourceTencentCloudTeoL7AccRuleV2Create` 中，当 `custom_variables` 存在时，将其展开为 `[]*CustomVariable` 并赋值给 `rule.CustomVariables`
- [x] 2.2 在 `ResourceTencentCloudTeoL7AccRuleV2Create` 中，当 `custom_variable_operations` 存在时，调用专用展开辅助函数生成 `[]*CustomVariableOperation` 并赋值给 `rule.CustomVariableOperations`

## 3. Update 函数（resource_tc_teo_l7_acc_rule_v2.go）

- [x] 3.1 在 `ResourceTencentCloudTeoL7AccRuleV2Update` 的 `HasChange` 守卫条件中追加 `custom_variables` 和 `custom_variable_operations`
- [x] 3.2 在 Update 的 rule 构建块中，将 `custom_variables` 展开为 `rule.CustomVariables`，将 `custom_variable_operations` 展开为 `rule.CustomVariableOperations`

## 4. Expand/Flatten 辅助函数（resource_tc_teo_l7_acc_rule_extension.go）

- [x] 4.1 新增 `custom_variable_operations` 的展开辅助函数：将 Terraform map 展开为 `[]*CustomVariableOperation`（含 `Branches`→`CustomVariableOperationRuleBranch`、`Actions`→`CustomVariableOperationRuleAction`+`SetParameters`、`SubRules`→`CustomVariableOperationSubRule`、`Description`），所有层级做非空判断
- [x] 4.2 新增 `custom_variable_operations` 的 flatten 辅助函数：将 SDK `[]*CustomVariableOperation` 扁平化为 Terraform map，每一层 SDK 字段非 nil 才 set，与展开函数对称
- [x] 4.3 在 `resourceTencentCloudTeoL7AccRuleGetBranchs` 的 actions 循环中新增 `custom_action_parameters` 展开块：读取 `actionsMap["custom_action_parameters"]`，展开为 `RuleEngineAction.CustomActionParameters`（含 `CustomActions` 列表，每个含 `ActionId` 与 `Parameters`，`Parameters` 中按需设置 typed value 字段），与现有 `shield_parameters`/`site_failover_parameters` 展开块并列
- [x] 4.4 在 `resourceTencentCloudTeoL7AccRuleSetBranchs` 的 actions flatten 中新增 `custom_action_parameters` flatten 块：当 `actions.CustomActionParameters` 非 nil 时，将其扁平化为 `actionsMap["custom_action_parameters"]`（含 `custom_actions` 列表，每个含 `action_id` 与 `parameters`，`parameters` 中将 populated 的 typed value 序列化为单一 `value` 字符串），与现有 `shield_parameters`/`site_failover_parameters` flatten 块并列

## 5. Read 函数（resource_tc_teo_l7_acc_rule_v2.go）

- [x] 5.1 在 `ResourceTencentCloudTeoL7AccRuleV2Read` 中，当 `respData.Rules[0]` 存在时，调用 flatten 辅助函数将 `rule.CustomVariables` 设置到 `custom_variables`、将 `rule.CustomVariableOperations` 设置到 `custom_variable_operations`，所有 SDK 字段非 nil 才 set

## 6. 单元测试（resource_tc_teo_l7_acc_rule_v2_test.go）

- [x] 6.1 使用 gomonkey mock 云 API 的方式，补充 `custom_variables` 参数的单元测试用例（mock `CreateL7AccRules`/`DescribeL7AccRules`/`ModifyL7AccRule`，验证 schema 与 CRUD 逻辑）
- [x] 6.2 补充 `custom_variable_operations` 参数（含 branches/actions/set_parameters/sub_rules 嵌套）的单元测试用例
- [x] 6.3 补充 `custom_action_parameters` 参数（含 custom_actions/action_id/parameters 嵌套）的单元测试用例

## 7. 验证

- [x] 7.1 使用 `git diff` 核对所有变更文件，确认 schema 定义、CRUD 展开/flatten 逻辑、测试用例与 design.md/spec.md 一致，且未引入未使用变量或编译错误
