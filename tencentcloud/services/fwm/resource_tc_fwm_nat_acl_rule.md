Provides a resource to create a FWM NAT ACL rule

Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "fwmrg-in3q0qkq0n"
  direction      = 0
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "accept"
  scope          = "ALL"
  port           = "-1/-1"
  description    = "remark."
}
```

Import

FWM NAT ACL rule can be imported using the groupId#ruleId#direction, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-in3q0qkq0n#7520b898-30ba-415d-ab27-ff4542071ea9#0
```
