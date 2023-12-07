Provides a resource to create a cfw edge_policy

Example Usage

```hcl
resource "tencentcloud_cfw_edge_policy" "example" {
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
  scope          = "all"
}
```

If target_type is tag

```hcl
resource "tencentcloud_cfw_edge_policy" "example" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = jsonencode({"Key":"test","Value":"dddd"})
  target_type    = "tag"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}
```

Import

cfw edge_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_edge_policy.example edge_policy_id
```