---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-nat_gateway"
description: |-
  Provides a resource to create a NAT gateway.
---

# tencentcloud_nat_gateway

Provides a resource to create a NAT gateway.

~> **NOTE:** If `nat_product_version` is `1`, `max_concurrent` valid values is `1000000`, `3000000`, `10000000`.

~> **NOTE:** If set `stock_public_ip_addresses_bandwidth_out`, do not set the `internet_max_bandwidth_out` parameter of resource `tencentcloud_eip` at the same time, otherwise conflicts may occur.

## Example Usage

### Create a traditional NAT gateway.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name                = "tf_example_nat_gateway"
  vpc_id              = tencentcloud_vpc.vpc.id
  nat_product_version = 1
  bandwidth           = 100
  max_concurrent      = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
  }
}
```

### Create a standard NAT gateway.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name                = "tf_example_nat_gateway"
  vpc_id              = tencentcloud_vpc.vpc.id
  nat_product_version = 2
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
  }
}
```

### Or set stock public ip addresses bandwidth out

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name                                    = "tf_example_nat_gateway"
  vpc_id                                  = tencentcloud_vpc.vpc.id
  nat_product_version                     = 2
  stock_public_ip_addresses_bandwidth_out = 100
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `assigned_eip_set` - (Required, Set: [`String`]) EIP IP address set bound to the gateway. The value of at least 1 and at most 10 if do not apply for a whitelist.
* `name` - (Required, String) Name of the NAT gateway.
* `vpc_id` - (Required, String, ForceNew) ID of the vpc.
* `bandwidth` - (Optional, Int) The maximum public network output bandwidth of NAT gateway (unit: Mbps). Valid values: `20`, `50`, `100`, `200`, `500`, `1000`, `2000`, `5000`. Default is `100`. When the value of parameter `nat_product_version` is 2, which is the standard NAT type, this parameter does not need to be filled in and defaults to `5000`.
* `max_concurrent` - (Optional, Int) The upper limit of concurrent connection of NAT gateway. Valid values: `1000000`, `3000000`, `10000000`. Default is `1000000`. When the value of parameter `nat_product_version` is 2, which is the standard NAT type, this parameter does not need to be filled in and defaults to `2000000`.
* `nat_product_version` - (Optional, Int, ForceNew) 1: traditional NAT, 2: standard NAT, default value is 1.
* `stock_public_ip_addresses_bandwidth_out` - (Optional, Int, ForceNew) The elastic public IP bandwidth value (unit: Mbps) for binding NAT gateway. When this parameter is not filled in, it defaults to the bandwidth value of the elastic public IP, and for some users, it defaults to the bandwidth limit of the elastic public IP of that user type.
* `subnet_id` - (Optional, String, ForceNew) Subnet of NAT.
* `tags` - (Optional, Map) The available tags within this NAT gateway.
* `zone` - (Optional, String) The availability zone, such as `ap-guangzhou-3`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Create time of the NAT gateway.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNatGateway](https://cloud.tencent.com/document/api/215/36721) - This API is used to create a NAT Gateway.
* Create: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Create: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - This API is used to query NAT gateways.
* Create: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Read: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - This API is used to query NAT gateways.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [AssociateNatGatewayAddress](https://cloud.tencent.com/document/api/215/36722) - This API is used to bind an EIP to a NAT gateway.
* Update: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Update: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - This API is used to query NAT gateways.
* Update: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeVpcTaskResult](https://cloud.tencent.com/document/api/215/59037) - This API is used to query the execution result of a VPC task.
* Update: [DisassociateNatGatewayAddress](https://cloud.tencent.com/document/api/215/36716) - This API is used to unbind an EIP from a NAT gateway.
* Update: [ModifyNatGatewayAttribute](https://cloud.tencent.com/document/api/215/36715) - This API is used to modify the attributes of a NAT gateway.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [ResetNatGatewayConnection](https://cloud.tencent.com/document/api/215/36713) - This API is used to adjust concurrent connection cap for the NAT gateway.
* Delete: [DeleteNatGateway](https://cloud.tencent.com/document/api/215/36719) - This API is used to delete a NAT gateway.
* Delete: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - This API is used to query NAT gateways.

## Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.example nat-1asg3t63
```

