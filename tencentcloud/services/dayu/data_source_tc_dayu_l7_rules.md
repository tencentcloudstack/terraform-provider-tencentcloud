Use this data source to query dayu layer 7 rules

Example Usage

```hcl
data "tencentcloud_dayu_l7_rules" "domain_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  domain        = tencentcloud_dayu_l7_rule.test_rule.domain
}
data "tencentcloud_dayu_l7_rules" "id_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.test_rule.rule_id
}
```