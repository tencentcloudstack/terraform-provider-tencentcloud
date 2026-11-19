---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_private_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-vpc_private_nat_gateway"
description: |-
  Provides a resource to create a VPC private nat gateway
---

# tencentcloud_vpc_private_nat_gateway

Provides a resource to create a VPC private nat gateway

## Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway" "example" {
  nat_gateway_name = "tf-example"
  vpc_id           = "vpc-i5yyodl9"
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_name` - (Required, String) Private network gateway name.
* `ccn_id` - (Optional, String) Cloud Connect Network type The Cloud Connect Network instance ID required to be bound to the private network NAT gateway.
* `cross_domain` - (Optional, Bool) Cross-domain parameters. Cross-domain binding of VPCs is supported only when the value is True.
* `tags` - (Optional, Map) Tag description of the instance.
* `vpc_id` - (Optional, String) Private Cloud instance ID. This parameter is required when creating a VPC type private network NAT gateway or a private network NAT gateway of private network gateway.
* `vpc_type` - (Optional, Bool) VPC type private network NAT gateway. Only when the value is True will a VPC type private network NAT gateway be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `nat_gateway_id` - Private network NAT gateway instance ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreatePrivateNatGateway](https://www.tencentcloud.com/search/?wd=CreatePrivateNatGateway&type=document&lang=en)
  - [DescribePrivateNatGateways](https://www.tencentcloud.com/search/?wd=DescribePrivateNatGateways&type=document&lang=en)
* Read: [DescribePrivateNatGateways](https://www.tencentcloud.com/search/?wd=DescribePrivateNatGateways&type=document&lang=en)
* Update:
  - [DescribePrivateNatGateways](https://www.tencentcloud.com/search/?wd=DescribePrivateNatGateways&type=document&lang=en)
  - [ModifyPrivateNatGatewayAttribute](https://www.tencentcloud.com/search/?wd=ModifyPrivateNatGatewayAttribute&type=document&lang=en)
* Delete:
  - [DeletePrivateNatGateway](https://www.tencentcloud.com/search/?wd=DeletePrivateNatGateway&type=document&lang=en)
  - [DescribePrivateNatGateways](https://www.tencentcloud.com/search/?wd=DescribePrivateNatGateways&type=document&lang=en)

## Import

VPC private nat gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway.example intranat-ljdy849x
```

