---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_private_nat_gateway_translation_acl_rule"
sidebar_current: "docs-tencentcloud-resource-vpc_private_nat_gateway_translation_acl_rule"
description: |-
  Provides a resource to create a VPC private nat gateway translation acl rule
---

# tencentcloud_vpc_private_nat_gateway_translation_acl_rule

Provides a resource to create a VPC private nat gateway translation acl rule

## Example Usage

### If translation_type is NETWORK_LAYER

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = "intranat-bw389ya1"
  translation_direction = "LOCAL"
  translation_type      = "NETWORK_LAYER"
  translation_ip        = "2.2.2.2"
  original_ip           = "1.1.1.1"
  translation_acl_rules {
    protocol         = "TCP"
    source_port      = "80"
    destination_port = "8080"
    destination_cidr = "8.8.8.8"
    description      = "remark."
  }
}
```

### If translation_type is NETWORK_LAYER

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = "intranat-bw389ya1"
  translation_direction = "LOCAL"
  translation_type      = "TRANSPORT_LAYER"
  translation_ip        = "3.3.3.3"
  translation_acl_rules {
    protocol         = "ALL"
    source_port      = "0-65535"
    source_cidr      = "2.2.2.2"
    destination_port = "0-65535"
    destination_cidr = "3.3.3.3"
    description      = "remark."
  }
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String, ForceNew) The unique ID of the private NAT gateway, in the format: `intranat-xxxxxxxx`.
* `translation_acl_rules` - (Required, List) Access control list.
* `translation_direction` - (Required, String, ForceNew) The target of the translation rule, optional value: LOCAL.
* `translation_ip` - (Required, String, ForceNew) The mapped IP address. When the translation rule type is layer 4, it represents an IP pool.
* `translation_type` - (Required, String, ForceNew) The type of translation rule, optional values: NETWORK_LAYER, TRANSPORT_LAYER. Corresponding to layer 3 and layer 4 respectively.
* `original_ip` - (Optional, String, ForceNew) The original IP address before mapping. Valid when the translation rule type is layer 3.

The `translation_acl_rules` object supports the following:

* `destination_cidr` - (Required, String) Destination address.
* `destination_port` - (Required, String) Destination port.
* `protocol` - (Required, String) ACL protocol type, optional values: `ALL`, `TCP`, `UDP`.
* `source_port` - (Required, String) Source port.
* `action` - (Optional, Int) Whether to match.
* `description` - (Optional, String) ACL rule description.
* `source_cidr` - (Optional, String) Source address. Supports `ip` or `cidr` format `xxx.xxx.xxx.000/xx`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPC private nat gateway translation acl rule can be imported using the natGatewayId#translationDirection#translationType#translationIp[#originalIp]#aclruleId, e.g.

```
# If translation_type is NETWORK_LAYER
terraform import tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example intranat-bw389ya1#LOCAL#NETWORK_LAYER#2.2.2.2#1.1.1.1#23

# If translation_type is NETWORK_LAYER
terraform import tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example intranat-bw389ya1#LOCAL#TRANSPORT_LAYER#3.3.3.3#25
```

