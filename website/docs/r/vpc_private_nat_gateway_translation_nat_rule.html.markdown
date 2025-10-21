---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_private_nat_gateway_translation_nat_rule"
sidebar_current: "docs-tencentcloud-resource-vpc_private_nat_gateway_translation_nat_rule"
description: |-
  Provides a resource to create a VPC private nat gateway translation nat rule
---

# tencentcloud_vpc_private_nat_gateway_translation_nat_rule

Provides a resource to create a VPC private nat gateway translation nat rule

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional translation nat rules resources of this nat gateway elsewhere.

## Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_nat_rule" "example" {
  nat_gateway_id = "intranat-r46f6pxl"
  translation_nat_rules {
    translation_direction = "LOCAL"
    translation_type      = "NETWORK_LAYER"
    translation_ip        = "2.2.2.2"
    description           = "remark."
    original_ip           = "1.1.1.1"
  }

  translation_nat_rules {
    translation_direction = "LOCAL"
    translation_type      = "TRANSPORT_LAYER"
    translation_ip        = "3.3.3.3"
    description           = "remark."
  }
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String, ForceNew) Private NAT gateway unique ID, such as: `intranat-xxxxxxxx`.
* `translation_nat_rules` - (Required, Set) Translation rule object array.

The `translation_nat_rules` object supports the following:

* `description` - (Required, String) Translation rule description.
* `translation_direction` - (Required, String) Translation rule target, optional values "LOCAL","PEER".
* `translation_ip` - (Required, String) Translation IP, when translation rule type is transport layer, it is an IP pool.
* `translation_type` - (Required, String) Translation rule type, optional values "NETWORK_LAYER","TRANSPORT_LAYER".
* `original_ip` - (Optional, String) Source IP, valid when translation rule type is network layer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPC private nat gateway translation nat rule can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example intranat-r46f6pxl
```

