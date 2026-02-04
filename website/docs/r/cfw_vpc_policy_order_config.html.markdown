---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_vpc_policy_order_config"
sidebar_current: "docs-tencentcloud-resource-cfw_vpc_policy_order_config"
description: |-
  Provides a resource to create a CFW vpc policy order config
---

# tencentcloud_cfw_vpc_policy_order_config

Provides a resource to create a CFW vpc policy order config

~> **NOTE:** If resource `tencentcloud_cfw_vpc_policy_order_config` is used to sort resource `tencentcloud_cfw_vpc_policy`, all instances of resource `tencentcloud_cfw_vpc_policy` must be configured simultaneously, and the sorting of this resource cannot be declared elsewhere.

## Example Usage

```hcl
resource "tencentcloud_cfw_vpc_policy" "example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  dest_content   = "192.168.0.1"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy" "example2" {
  source_content = "2.2.2.2/0"
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
  source_content = "3.3.3.3/0"
  source_type    = "net"
  dest_content   = "192.168.0.3"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy_order_config" "example" {
  rule_uuid_list = [
    tencentcloud_cfw_vpc_policy.example3.uuid,
    tencentcloud_cfw_vpc_policy.example1.uuid,
    tencentcloud_cfw_vpc_policy.example2.uuid,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `rule_uuid_list` - (Optional, List: [`Int`]) The unique IDs of the rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CFW vpc policy order config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_vpc_policy_order_config.example GedqV07VpNU0ob8LuOXw==
```

