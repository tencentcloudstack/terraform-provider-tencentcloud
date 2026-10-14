---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_local_gateway"
sidebar_current: "docs-tencentcloud-resource-vpc_local_gateway"
description: |-
  Provides a resource to create a vpc local_gateway
---

# tencentcloud_vpc_local_gateway

Provides a resource to create a vpc local_gateway

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_local_gateway" "example" {
  local_gateway_name = "tf-example"
  vpc_id             = tencentcloud_vpc.vpc.id
  cdc_id             = "cluster-j9gyu1iy"
}
```

## Argument Reference

The following arguments are supported:

* `cdc_id` - (Required, String) CDC instance ID.
* `local_gateway_name` - (Required, String) Local gateway name.
* `vpc_id` - (Required, String) VPC instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLocalGateway](https://www.tencentcloud.com/document/product/215/40144) - This API is used to create a local gateway for a CDC instance.
* Read: [DescribeLocalGateway](https://www.tencentcloud.com/document/product/215/40142) - This API is used to query local gateways of a CDC instance.
* Update: [ModifyLocalGateway](https://www.tencentcloud.com/document/product/215/40141) - This API is used to modify the local gateway of a CDC instance.
* Delete: [DeleteLocalGateway](https://www.tencentcloud.com/document/product/215/40143) - This API is used to delete the local gateway of a CDC instance.

## Import

vpc local_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_local_gateway.local_gateway local_gateway_id
```

