Provides a resource to create a Edge ACL rule of Cloud Firewall (CFW)

Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id       = "fwmrg-9z03vakjj8"
  direction      = 1
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  scope          = "all"
  description    = "remark."
}
```

Import

FWM edge acl rule can be imported using the groupId#ruleId#direction, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example fwmrg-9z03vakjj8#2fe0b492-5623-4646-b40e-7b0908307ca8#1
```