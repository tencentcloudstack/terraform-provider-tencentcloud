Provides a resource to create a VPC private nat gateway translation acl rule

Example Usage

If translation_type is NETWORK_LAYER

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

If translation_type is NETWORK_LAYER

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

Import

VPC private nat gateway translation acl rule can be imported using the natGatewayId#translationDirection#translationType#translationIp[#originalIp]#aclruleId, e.g.

```
# If translation_type is NETWORK_LAYER
terraform import tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example intranat-bw389ya1#LOCAL#NETWORK_LAYER#2.2.2.2#1.1.1.1#23

# If translation_type is NETWORK_LAYER
terraform import tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example intranat-bw389ya1#LOCAL#TRANSPORT_LAYER#3.3.3.3#25
```
