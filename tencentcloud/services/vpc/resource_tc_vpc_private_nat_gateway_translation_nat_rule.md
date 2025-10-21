Provides a resource to create a VPC private nat gateway translation nat rule

~> **NOTE:** 

Example Usage

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

Import

VPC private nat gateway translation nat rule can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example intranat-r46f6pxl
```
