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

~> **NOTE:** Append-only convention for typed lists: The provider preserves the API's authoritative order in state (rules are currently ordered by creation time at the backend, and may gain user-defined ordering semantics in the future). To keep plans clean and avoid spurious churn, **add new rules ONLY at the end of each typed list, and do NOT reorder or insert rules in the middle of an existing list**. Inserting in the middle, removing from the middle, or reordering will produce a connected `~` plan and may trigger uniqueness errors at apply time (e.g. `MODIFY` on slot N attempts to rename rule A's IP to rule B's still-existing IP). When that happens, the SDK error message points at the offending IP — fix the HCL by appending instead of inserting.

## Example Usage

### Recommended typed list usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_nat_rule" "example" {
  nat_gateway_id = "intranat-r46f6pxl"

  local_network_layer_rules {
    translation_ip = "2.2.2.2"
    original_ip    = "1.1.1.1"
    description    = "LOCAL three-layer rule."
  }

  local_transport_layer_rules {
    translation_ip = "3.3.3.3"
    description    = "LOCAL four-layer rule."
  }

  peer_network_layer_rules {
    translation_ip = "5.5.5.5"
    original_ip    = "4.4.4.4"
    description    = "PEER three-layer rule."
  }
}
```

### Deprecated translation_nat_rules usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_nat_rule" "example_deprecated" {
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
* `local_network_layer_rules` - (Optional, List) Translation rules for the LOCAL direction at the NETWORK_LAYER (three-layer). Identity is keyed by `original_ip` (unique within this bucket). Editing `description` or `translation_ip` is applied in place via ModifyPrivateNatGatewayTranslationNatRule; editing `original_ip` is treated as deleting the old rule and creating a new one.
* `local_transport_layer_rules` - (Optional, List) Translation rules for the LOCAL direction at the TRANSPORT_LAYER (four-layer). Identity is keyed by `translation_ip` (unique within this bucket). `original_ip` is not applicable for transport-layer rules and is intentionally not exposed.
* `peer_network_layer_rules` - (Optional, List) Translation rules for the PEER direction at the NETWORK_LAYER (three-layer). Identity is keyed by `original_ip` (unique within this bucket). Editing `description` or `translation_ip` is applied in place via ModifyPrivateNatGatewayTranslationNatRule; editing `original_ip` is treated as deleting the old rule and creating a new one.
* `translation_nat_rules` - (Optional, Set, **Deprecated**) It has been deprecated from version 1.82.98, please use local_network_layer_rules / local_transport_layer_rules / peer_network_layer_rules instead. Cannot be used together with the new fields. (Deprecated) Translation rule object array. Use the typed list fields `local_network_layer_rules`, `local_transport_layer_rules`, and `peer_network_layer_rules` instead. The legacy field continues to work but does not benefit from in-place ModifyPrivateNatGatewayTranslationNatRule support.

The `local_network_layer_rules` object supports the following:

* `original_ip` - (Required, String) Original (mapped-before) IP for this rule. Acts as the rule identity within this bucket.
* `translation_ip` - (Required, String) Translated (mapped-after) IP for this rule. Can be modified in place.
* `description` - (Optional, String) Translation rule description.

The `local_transport_layer_rules` object supports the following:

* `translation_ip` - (Required, String) Translated IP pool for this rule (transport-layer rules use an IP pool). Acts as the rule identity within this bucket.
* `description` - (Optional, String) Translation rule description.

The `peer_network_layer_rules` object supports the following:

* `original_ip` - (Required, String) Original (mapped-before) IP for this rule. Acts as the rule identity within this bucket.
* `translation_ip` - (Required, String) Translated (mapped-after) IP for this rule. Can be modified in place.
* `description` - (Optional, String) Translation rule description.

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

