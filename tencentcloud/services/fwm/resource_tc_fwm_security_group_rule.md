Provides a resource to manage an enterprise security group rule of Tencent Cloud Firewall Manager (FWM).

Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_name     = "tf-example-fwm-group"
  ip_version     = "4"
  source_content = "1.1.1.1/0"
  source_type    = "net"
  dest_content   = "0.0.0.0/0"
  dest_type      = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  description    = "rule description."
}
```

To add the rule into an existing enterprise security group rule group instead of letting this resource create one, use `group_id`:

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id       = "fwmrg-xxxxxxxx"
  ip_version     = "4"
  source_content = "1.1.1.1/0"
  source_type    = "net"
  dest_content   = "0.0.0.0/0"
  dest_type      = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  description    = "rule description."
}
```

Import

Tencent Cloud Firewall Manager enterprise security group rule can be imported using the composite id, which is `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example fwmrg-xxxxxxxx#fwmrule-yyyyyyyy
```

-> **Note:** Both segments of the composite id are required, because every read/update/delete call of this resource needs the rule group id together with the rule id.
