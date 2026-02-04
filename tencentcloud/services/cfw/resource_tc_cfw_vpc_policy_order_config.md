Provides a resource to create a CFW vpc policy order config

~> **NOTE:** If resource `tencentcloud_cfw_vpc_policy_order_config` is used to sort resource `tencentcloud_cfw_vpc_policy`, all instances of resource `tencentcloud_cfw_vpc_policy` must be configured simultaneously, and the sorting of this resource cannot be declared elsewhere.

Example Usage

```hcl
resource "tencentcloud_cfw_vpc_policy" "example1" {
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

resource "tencentcloud_cfw_vpc_policy" "example2" {
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

resource "tencentcloud_cfw_vpc_policy" "example3" {
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

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  rule_uuid_list = [
    tencentcloud_cfw_vpc_policy.in_example3.uuid,
    tencentcloud_cfw_vpc_policy.in_example1.uuid,
    tencentcloud_cfw_vpc_policy.in_example2.uuid,
  ]
}
```

Import

CFW vpc policy order config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_vpc_policy_order_config.example GedqV07VpNU0ob8LuOXw==
```
