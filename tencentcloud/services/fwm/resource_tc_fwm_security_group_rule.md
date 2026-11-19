Provides a resource to create a FWM security group rule.

Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id           = "fwmrg-ghlu11dv"
  ip_version         = "ipv4"
  source_type        = "net"
  source_content     = "0.0.0.0/0"
  dest_type          = "net"
  dest_content       = "0.0.0.0/0"
  protocol           = "ANY"
  port               = "-1/-1"
  rule_action        = "drop"
  scope              = "SG"
  protocol_port_type = 0
  description        = "remark."
}
```

Import

FWM security group rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example fwmrg-ghlu11dv#a5f39ae7-2a81-4735-9027-5d6f83d814a5
```