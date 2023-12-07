Use this resource to create a dayu CC self-define https policy

~> **NOTE:** creating CC self-define https policy need a valid resource `tencentcloud_dayu_l7_rule`; The resource only support Anti-DDoS of resource type `bgpip`.

Example Usage

```hcl
resource "tencentcloud_dayu_cc_https_policy" "test_policy" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.test_rule.rule_id
  domain        = tencentcloud_dayu_l7_rule.test_rule.domain
  name          = "policy_test"
  action        = "drop"
  switch        = true

  rule_list {
    skey     = "cgi"
    operator = "include"
    value    = "123"
  }
}

```