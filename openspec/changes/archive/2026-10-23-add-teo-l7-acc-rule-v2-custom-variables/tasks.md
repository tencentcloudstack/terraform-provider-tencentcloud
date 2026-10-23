## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go` 的 `TencentTeoL7RuleBranchBasicInfo` 函数的 actions Elem.Schema 中新增 `custom_action_parameters`（TypeList, Optional, MaxItems=1），其 Elem.Schema 含 `custom_actions`（TypeList, Optional），每个 custom_actions 元素含 `action_id`（TypeString, Optional）与 `parameters`（TypeList, Optional，含 `name`/`value_type`/`string_value`/`integer_value`/`float_value`/`boolean_value`/`string_array_value`/`integer_array_value`/`float_array_value`）
- [x] 1.2 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go` 中新增独立的 `custom_variable_operations` 分支 schema 定义函数（含 branches（condition/actions（name/set_parameters（name/value））/sub_rules（branches 递归/description））/description），对应 SDK 的 `CustomVariableOperationRuleBranch`/`CustomVariableOperationRuleAction`/`SetParameters`/`CustomVariableOperationSubRule`
- [x] 1.3 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` 的 `ResourceTencentCloudTeoL7AccRuleV2()` schema map 顶层新增 `custom_variables`（TypeList, Optional，含 `name`/`initial_value`/`description`）与 `custom_variable_operations`（TypeList, Optional，Elem 复用 1.2 的 schema）
- [x] 1.4 更新 actions 的 `name` 字段 Description，补充 `CustomAction`（定制配置）作为合法操作名

## 2. CRUD 函数扩展

- [x] 2.1 在 `ResourceTencentCloudTeoL7AccRuleV2Create` 中，新增对 `custom_variables` 与 `custom_variable_operations` 的读取与映射，设置到 `rule.CustomVariables` / `rule.CustomVariableOperations`（request.Rules[0]）
- [x] 2.2 在 `ResourceTencentCloudTeoL7AccRuleV2Update` 的 `d.HasChange` 条件中加入 `custom_variables` / `custom_variable_operations`，并在 rule 构造分支中追加这两个字段的映射到 `ModifyL7AccRuleRequest.Rule`
- [x] 2.3 在 `ResourceTencentCloudTeoL7AccRuleV2Read` 中，从 `respData.Rules[0]` 读取 `CustomVariables` / `CustomVariableOperations`，回填到 `d.Set("custom_variables", ...)` / `d.Set("custom_variable_operations", ...)`，回填前判断 nil

## 3. 辅助函数扩展

- [x] 3.1 新增 `resourceTencentCloudTeoL7AccRuleGetCustomVariables` 函数，将 TF state 的 custom_variables 列表转换为 `[]*teov20220901.CustomVariable`
- [x] 3.2 新增 `resourceTencentCloudTeoL7AccRuleSetCustomVariables` 函数，将 `[]*teov20220901.CustomVariable` 扁平化为 TF state 列表（字段 nil 安全判断）
- [x] 3.3 新增 `resourceTencentCloudTeoL7AccRuleGetCustomVariableOperations` 函数，将 TF state 的 custom_variable_operations（含递归 branches/actions/set_parameters/sub_rules）转换为 `[]*teov20220901.CustomVariableOperation`
- [x] 3.4 新增 `resourceTencentCloudTeoL7AccRuleSetCustomVariableOperations` 函数，将 `[]*teov20220901.CustomVariableOperation` 扁平化为 TF state（字段 nil 安全判断）
- [x] 3.5 在 `resourceTencentCloudTeoL7AccRuleGetBranchs` 的 actions 遍历中追加 `custom_action_parameters` 分支，将 custom_actions（含 action_id/parameters）映射到 `RuleEngineAction.CustomActionParameters`
- [x] 3.6 在 `resourceTencentCloudTeoL7AccRuleSetBranchs` 的 actions 遍历中追加 `custom_action_parameters` 分支，将 `RuleEngineAction.CustomActionParameters` 扁平化回 TF state（字段 nil 安全判断）

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` 中新增 `custom_variables` 的 Create/Read/Update 单元测试用例（gomonkey mock `CreateL7AccRules` / `DescribeL7AccRules` / `ModifyL7AccRule`），断言 `RuleEngineItem.CustomVariables` 映射正确
- [x] 4.2 新增 `custom_variable_operations`（含 branches/actions/set_parameters/sub_rules/description）的 Create/Read/Update 单元测试用例，断言嵌套结构映射正确
- [x] 4.3 新增 `custom_action_parameters`（在 actions 内，含 custom_actions/action_id/parameters 多类型 Value）的 Create/Read/Update 单元测试用例，断言 `RuleEngineAction.CustomActionParameters` 映射正确
- [x] 4.4 保证已有测试用例（不含新参数的场景）继续通过，不修改其行为

## 5. 文档同步

- [x] 5.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` 中新增展示 `custom_variables` 与 `custom_variable_operations` 用法的 HCL 示例
- [x] 5.2 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` 中新增展示 `custom_action_parameters`（在 actions 内）用法的 HCL 示例
- [ ] 5.3 执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website/ 目录）

## 6. 验证

- [ ] 6.1 运行 `go build ./...` 确保编译通过
- [ ] 6.2 运行 `go vet ./tencentcloud/services/teo/...` 确保无静态检查错误
- [ ] 6.3 运行 `go test ./tencentcloud/services/teo/ -run TestAccTeoL7AccRuleV2 -v -count=1 -gcflags="all=-l"` 确认所有单测通过（含新增用例）