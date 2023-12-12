Provides a resource to create a cfw vpc_policy

Example Usage

```hcl
resource "tencentcloud_cfw_vpc_policy" "example" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}
```

Import

cfw vpc_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_policy.vpc_policy vpc_policy_id
```