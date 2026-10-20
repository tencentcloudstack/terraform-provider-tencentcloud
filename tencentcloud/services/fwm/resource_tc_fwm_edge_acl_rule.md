Provides a resource to create a FWM edge acl rule

Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id = "group-xxxx"
  rules {
    order_index     = 1
    direction       = 1
    source_content  = "1.1.1.1/32"
    source_type     = "ip"
    target_content  = "0.0.0.0/0"
    target_type     = "ip"
    port            = "80"
    protocol        = "TCP"
    rule_action     = "accept"
    description     = "example rule"
    scope           = "serial"
  }
}
```

Import

FWM edge acl rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example group-xxxx#rule-xxxx
```