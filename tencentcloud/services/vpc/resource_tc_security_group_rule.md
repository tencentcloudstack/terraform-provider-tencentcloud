Provides a resource to create security group rule.

~> **NOTE:** It is recommended to use `tencentcloud_security_group_rule_set` to manage multiple security group rules, because the index number of a single rule changes dynamically as other rules change.

Example Usage

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group example."
}

resource "tencentcloud_security_group_rule" "example" {
  security_group_id = tencentcloud_security_group.example.id
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80,8080"
  policy            = "ACCEPT"
  description       = "allow access to the web service."
}
```
