---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway_snat"
sidebar_current: "docs-tencentcloud-resource-nat_gateway_snat"
description: |-
  Provides a resource to create a NAT Gateway SNat rule.
---

# tencentcloud_nat_gateway_snat

Provides a resource to create a NAT Gateway SNat rule.

## Example Usage

### If resource_type is SUBNET

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id    = "nat-jmf9955g"
  resource_type     = "SUBNET"
  subnet_id         = "subnet-hhi88a58"
  subnet_cidr_block = "10.0.30.0/24"
  description       = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

### If resource_type is NETWORKINTERFACE

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id           = "nat-jmf9955g"
  resource_type            = "NETWORKINTERFACE"
  instance_id              = "ins-p97oyxt6"
  instance_private_ip_addr = "10.0.30.62"
  description              = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

### If resource_type is USERDEFINED

```hcl
resource "tencentcloud_nat_gateway_snat" "example" {
  nat_gateway_id           = "nat-jmf9955g"
  resource_type            = "USERDEFINED"
  instance_private_ip_addr = "10.0.40.0/24"
  description              = "remark."
  public_ip_addr = [
    "106.55.167.123",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Description.
* `nat_gateway_id` - (Required, String, ForceNew) NAT gateway ID.
* `public_ip_addr` - (Required, List: [`String`]) Elastic IP address pool.
* `resource_type` - (Required, String, ForceNew) Resource type. Valid values: SUBNET, NETWORKINTERFACE, USERDEFINED.
* `instance_id` - (Optional, String, ForceNew) Instance ID, required when `resource_type` is NETWORKINTERFACE.
* `instance_private_ip_addr` - (Optional, String, ForceNew) Private IPs of the instance's primary ENI, required when `resource_type` is NETWORKINTERFACE or USERDEFINED.
* `subnet_cidr_block` - (Optional, String, ForceNew) The IPv4 CIDR of the subnet, required when `resource_type` is SUBNET.
* `subnet_id` - (Optional, String, ForceNew) Subnet instance ID, required when `resource_type` is SUBNET.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `snat_id` - SNAT rule ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateNatGatewaySourceIpTranslationNatRule](https://www.tencentcloud.com/document/api/215/39646)
  - [DescribeNatGatewaySourceIpTranslationNatRules](https://www.tencentcloud.com/document/api/215/39644)
* Read: [DescribeNatGatewaySourceIpTranslationNatRules](https://www.tencentcloud.com/document/api/215/39644)
* Update: [ModifyNatGatewaySourceIpTranslationNatRule](https://www.tencentcloud.com/document/api/215/39643)
* Delete: [DeleteNatGatewaySourceIpTranslationNatRule](https://www.tencentcloud.com/document/api/215/39645)

## Import

VPN gateway route can be imported using the id, the id format must be '{nat_gateway_id}#{resource_id}', resource_id range `subnet_id`, `instance_id` and `""`, e.g.

SUBNET SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#subnet-hhi88a58
```

NETWORKINTERFACT SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#ins-p97oyxt6
```

USERDEFINED SNat

```
terraform import tencentcloud_nat_gateway_snat.example nat-jmf9955g#
```

