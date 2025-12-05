Provides a resource to create a VPC private nat gateway translation acl rule

Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = ""
  translation_direction = ""
  translation_type      = ""
  translation_ip        = ""
  original_ip           = ""
  translation_acl_rule {
    protocol         = ""
    source_port      = ""
    source_cidr      = ""
    destination_port = ""
    destination_cidr = ""
    action           = ""
    description      = ""
  }
}
```

Import

VPC private nat gateway translation acl rule can be imported using the natGatewayId#translationDirection#translationType#translationIp#aclRuleId, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example vpc_private_nat_gateway_translation_acl_rule_id
```
