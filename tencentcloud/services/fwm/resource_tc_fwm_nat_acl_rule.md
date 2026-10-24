Provides a resource to create a FWM NAT ACL rule

Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id        = "fwmrg-abc"
  direction       = 0
  source_content  = "1.1.1.1"
  source_type     = "ip"
  target_content  = "0.0.0.0/0"
  target_type     = "ip"
  protocol        = "TCP"
  rule_action     = "accept"
  order_index     = 1
  scope           = "ALL"
  port            = "80"
  description     = "tf example nat acl rule"
}
```

Import

FWM NAT ACL rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-abc#fwmrule-xyz
```
