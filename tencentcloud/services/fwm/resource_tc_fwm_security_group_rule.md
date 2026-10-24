Provides a resource to create a FWM security group rule.

~> **NOTE:** The `tencentcloud_fwm_security_group_rule` resource belongs to a rule group (Product=enterprise_sg) which must be created in advance by `tencentcloud_fwm_rule_group`. This resource does not create the rule group.

Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id        = "group-id-example"
  ip_version      = "ipv4"
  order_index     = 1
  source_type      = "net"
  source_content   = "1.1.1.1"
  dest_type        = "net"
  dest_content     = "2.2.2.2"
  protocol         = "TCP"
  port             = "80"
  rule_action      = "accept"
  scope            = "SG"
  description      = "tf example"
}
```

Import

FWM security group rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example group-id-example#rule-id-example
```