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

## Example Usage

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
  description    = "111"
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  uuid        = tencentcloud_cfw_nat_policy.example.id
  order_index = 1
}
```

## Argument Reference

The following arguments are supported:

* `order_index` - (Required, Int) Rule sequence number.
* `uuid` - (Required, Int) The unique ID of the rule, which is not required when you create a rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CFW nat policy order config can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_policy_order_config.example 151517
```

