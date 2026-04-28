## 1. Schema 与 CRUD 函数修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` 的 Schema 中添加 `rule_ids` 参数定义（TypeList, Computed, Elem: TypeString, Description）
- [x] 1.2 在 `ResourceTencentCloudTeoL7AccRuleV2Create` 函数中，CreateL7AccRules 调用成功后，将 `result.Response.RuleIds` 转换为 `[]string` 并通过 `d.Set("rule_ids", ...)` 设置到 ResourceData
- [x] 1.3 在 `ResourceTencentCloudTeoL7AccRuleV2Read` 函数中，从 DescribeL7AccRules 返回的数据中提取 rule IDs 列表并设置到 `rule_ids`

## 2. 单元测试

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` 中补充 `rule_ids` 参数相关的单元测试用例，使用 gomonkey mock 云 API

## 3. 文档更新

- [x] 3.1 更新 `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` 文件，补充 `rule_ids` 参数的说明

## 4. 代码验证

- [x] 4.1 使用 `go test -gcflags=all=-l` 运行 `resource_tc_teo_l7_acc_rule_v2_test.go` 中的单元测试，确保测试通过
