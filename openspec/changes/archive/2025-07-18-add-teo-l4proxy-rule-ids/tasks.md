## 1. Schema & CRUD 代码修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule.go` 的 Schema 中添加 `l4proxy_rule_ids` 计算属性（TypeList of TypeString, Computed: true）
- [x] 1.2 在 `resourceTencentCloudTeoL4ProxyRuleCreate` 函数中，API 调用成功后将 `response.Response.L4ProxyRuleIds` 设置到 `l4proxy_rule_ids` 属性
- [x] 1.3 在 `resourceTencentCloudTeoL4ProxyRuleRead` 函数中，从复合 ID 中提取 ruleId 并设置到 `l4proxy_rule_ids` 属性

## 2. 单元测试

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule_test.go` 中补充 `l4proxy_rule_ids` 字段的单元测试用例

## 3. 文档更新

- [x] 3.1 更新 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule.md` 示例文件，添加 `l4proxy_rule_ids` 属性说明
