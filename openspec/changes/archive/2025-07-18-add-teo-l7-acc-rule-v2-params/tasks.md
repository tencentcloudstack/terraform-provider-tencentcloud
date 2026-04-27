## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` 的 Schema 中添加 `rule_ids` 字段，定义为 `TypeList`、`Computed: true`、`Elem: &schema.Schema{Type: schema.TypeString}`，描述为 "Rule ID list."

## 2. CRUD 函数修改

- [x] 2.1 在 `ResourceTencentCloudTeoL7AccRuleV2Create` 函数中，CreateL7AccRules API 调用成功后，将 `result.Response.RuleIds` 扁平化并设置到 `d.Set("rule_ids", ...)` 中
- [x] 2.2 在 `ResourceTencentCloudTeoL7AccRuleV2Read` 函数中，从 `respData.Rules` 中收集所有 `RuleId` 并设置到 `d.Set("rule_ids", ...)` 中

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` 中补充 `rule_ids` 字段的单元测试用例，使用 gomonkey mock 云 API

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` 文件，添加 `rule_ids` 属性的文档说明
