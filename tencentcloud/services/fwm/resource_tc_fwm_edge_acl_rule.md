Provides a resource to create a Edge ACL rule of Cloud Firewall (CFW)

Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id       = "cfw-group-xxxx"
  order_index    = 1
  direction      = 0
  source_content = "1.1.1.1"
  source_type    = "ip"
  target_content = "2.2.2.2"
  target_type     = "ip"
  protocol        = "TCP"
  port            = "80"
  rule_action     = "accept"
  scope           = "all"
}
```

Import

FWM edge acl rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example cfw-group-xxxx#cfw-rule-xxxx
```