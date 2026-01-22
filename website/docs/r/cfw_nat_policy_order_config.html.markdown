---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_policy_order_config"
sidebar_current: "docs-tencentcloud-resource-cfw_nat_policy_order_config"
description: |-
  Provides a resource to create a CFW nat policy order config
---

# tencentcloud_cfw_nat_policy_order_config

Provides a resource to create a CFW nat policy order config

~> **NOTE:** If resource `tencentcloud_cfw_nat_policy_order_config` is used to sort resource `tencentcloud_cfw_nat_policy`, all instances of resource `tencentcloud_cfw_nat_policy` must be configured simultaneously, and the sorting of this resource cannot be declared elsewhere.

~> **NOTE:** At any given time, resource `tencentcloud_cfw_nat_policy_order_config` can only be sorted against resources `tencentcloud_cfw_nat_policy` of the same `direction`.

## Example Usage

```hcl
resource "tencentcloud_cfw_nat_policy" "in_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example3" {
  source_content = "6.6.6.6/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "UDP"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}


resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_nat_policy.in_example3.uuid,
    tencentcloud_cfw_nat_policy.in_example1.uuid,
    tencentcloud_cfw_nat_policy.in_example2.uuid,
  ]

  outbound_rule_uuid_list = [
    tencentcloud_cfw_nat_policy.out_example2.uuid,
    tencentcloud_cfw_nat_policy.out_example1.uuid,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `inbound_rule_uuid_list` - (Optional, List: [`Int`]) The unique IDs of the inbound rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.
* `outbound_rule_uuid_list` - (Optional, List: [`Int`]) The unique IDs of the outbound rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CFW nat policy order config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_nat_policy_order_config.example GedqV07VpNU0ob8LuOXw==
```

