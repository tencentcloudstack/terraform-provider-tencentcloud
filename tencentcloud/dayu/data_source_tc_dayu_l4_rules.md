Use this data source to query dayu layer 4 rules

Example Usage

```hcl
data "tencentcloud_dayu_l4_rules" "name_test" {
  resource_type = tencentcloud_dayu_l4_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.test_rule.resource_id
  name          = tencentcloud_dayu_l4_rule.test_rule.name
}
data "tencentcloud_dayu_l4_rules" "id_test" {
  resource_type = tencentcloud_dayu_l4_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l4_rule.test_rule.rule_id
}
```