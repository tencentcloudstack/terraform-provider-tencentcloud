Provides a resource to create a cfw sg_rule

Example Usage

```hcl
resource "tencentcloud_sg_rule" "sg_rule" {
  enable = 1

  data {
    description         = "1111112"
    dest_content        = "0.0.0.0/0"
    dest_type           = "net"
    port                = "-1/-1"
    protocol            = "ANY"
    rule_action         = "accept"
    service_template_id = "ppm-l9u5pf1y"
    source_content      = "0.0.0.0/0"
    source_type         = "net"
  }
}
```

Import

cfw sg_rule can be imported using the id, e.g.

```
terraform import tencentcloud_sg_rule.sg_rule rule_id
```
