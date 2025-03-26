Provides a resource to create a CFW nat policy

Example Usage

```hcl
resource "tencentcloud_cfw_nat_policy" "example" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "ALL"
}
```

Import

CFW nat policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_policy.example 134123
```