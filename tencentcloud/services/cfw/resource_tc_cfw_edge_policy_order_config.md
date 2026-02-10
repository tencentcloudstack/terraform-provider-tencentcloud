Provides a resource to create a CFW edge policy order config

~> **NOTE:** If resource `tencentcloud_cfw_edge_policy_order_config` is used to sort resource `tencentcloud_cfw_edge_policy`, all instances of resource `tencentcloud_cfw_edge_policy` must be configured simultaneously, and the sorting of this resource cannot be declared elsewhere.

~> **NOTE:** At any given time, resource `tencentcloud_cfw_edge_policy_order_config` can only be sorted against resources `tencentcloud_cfw_edge_policy` of the same `direction`.

Example Usage

```hcl
resource "tencentcloud_cfw_edge_policy" "in_example1" {
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

resource "tencentcloud_cfw_edge_policy" "in_example2" {
  source_content = "2.2.2.2/0"
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

resource "tencentcloud_cfw_edge_policy" "in_example3" {
  source_content = "3.3.3.3/0"
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

resource "tencentcloud_cfw_edge_policy" "out_example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "out_example2" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.in_example3.uuid,
    tencentcloud_cfw_edge_policy.in_example1.uuid,
    tencentcloud_cfw_edge_policy.in_example2.uuid,
  ]

  outbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.out_example2.uuid,
    tencentcloud_cfw_edge_policy.out_example1.uuid,
  ]
}
```

Import

CFW edge policy order config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_edge_policy_order_config.example GedqV07VpNU0ob8LuOXw==
```
