Provides a resource to create a VPC private nat gateway translation nat rule

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional translation nat rules resources of this nat gateway elsewhere.

~> **NOTE:** Append-only convention for typed lists: The provider preserves the API's authoritative order in state (rules are currently ordered by creation time at the backend, and may gain user-defined ordering semantics in the future). To keep plans clean and avoid spurious churn, **add new rules ONLY at the end of each typed list, and do NOT reorder or insert rules in the middle of an existing list**. Inserting in the middle, removing from the middle, or reordering will produce a connected `~` plan and may trigger uniqueness errors at apply time (e.g. `MODIFY` on slot N attempts to rename rule A's IP to rule B's still-existing IP). When that happens, the SDK error message points at the offending IP — fix the HCL by appending instead of inserting.

Example Usage

Recommended typed list usage

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

Deprecated translation_nat_rules usage

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

Import

VPC private nat gateway translation nat rule can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example intranat-r46f6pxl
```
